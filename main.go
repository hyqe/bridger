package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hyqe/bridger/internal/rest"
)

func main() {
	rest.Serve(ttl(time.Minute, handler()))
}

func ttl(
	ttl time.Duration,
	next http.Handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ttl)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func handler() http.HandlerFunc {
	bridger := NewBridger()
	return func(w http.ResponseWriter, r *http.Request) {
		receiverId := strings.TrimPrefix(r.URL.Path, "/")
		senderId := r.URL.Query().Get("from")

		to := receiverId + senderId
		from := senderId + receiverId

		switch r.Method {

		// PUT /<bridgeId> {DATA}
		case http.MethodPut:
			bridger.Join(to, from, w, r)

		// GET /<bridgeId>
		// only retrieves data from the other client. it does not send data.
		case http.MethodGet:

		default:
			http.Error(w, "", http.StatusNotFound)
		}
	}
}

type Conn struct {
	mu  sync.Mutex
	ctx context.Context
	w   http.ResponseWriter
}

func newConn(ctx context.Context, w http.ResponseWriter) *Conn {
	c := &Conn{
		ctx: ctx,
		w:   w,
	}
	c.mu.Lock()
	return c
}

// Close a connection. this is not safe to be called multiple times.
func (c *Conn) Close() {
	c.mu.Unlock()
}

func (c *Conn) Wait() {
	c.mu.Lock()
}

type Bridger struct {
	sync.Mutex
	index map[string]*Bridge
}

func NewBridger() *Bridger {
	return &Bridger{
		index: make(map[string]*Bridge),
	}
}

type Bridge struct {
	sync.Mutex
	to      string
	from    string
	send    *Conn
	receive chan *Conn
}

func (b *Bridger) Join(to, from string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conn := newConn(ctx, w)

	select {

	case <-ctx.Done():
		http.Error(w, "request polling limit reached", http.StatusRequestTimeout)
		return

	case remote, ok := <-b.bridge(to, from, conn):
		if !ok {
			http.Error(w, "unable to connect", http.StatusConflict)
			return
		}
		remote.w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		_, err := io.Copy(remote.w, r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable forward body: %v", err), http.StatusInternalServerError)
			return
		}
		remote.Close()
		conn.Wait()
	}
}

// bridge shares the connections of two clients, so they can write
// to eachother.
func (b *Bridger) bridge(to, from string, send *Conn) chan *Conn {
	receive := make(chan *Conn)
	go func() {
		b.Lock()
		defer b.Unlock()

		// check if a remote "to" bridge is already waiting
		bridge, ok := b.index[to]

		// the remote is waiting
		if ok {
			bridge.Lock()
			defer bridge.Unlock()
			bridge.receive <- send
			close(bridge.receive)
			receive <- bridge.send
			close(receive)
			delete(b.index, to)
			return
		}

		// check if a "from" bridge already exists
		if _, ok := b.index[from]; ok {
			// do nothing the index already exists
			close(receive)
			return
		}

		// add the "from" bridge.
		// the client was the first to join.
		b.index[from] = &Bridge{
			to:      to,
			from:    from,
			send:    send,
			receive: receive,
		}

		// when the host connection to the bridge ends, the bridge will be cleaned up.
		clean := func() {
			<-send.ctx.Done()
			b.Lock()
			defer b.Unlock()

			if bridge, ok := b.index[from]; ok {
				bridge.Lock()
				defer bridge.Unlock()
				delete(b.index, from)
			}
		}
		go clean()
	}()
	return receive
}
