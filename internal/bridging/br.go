package bridging

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Claim struct {
	BridgeId  string    `json:"i"`
	ExpiresAt time.Time `json:"e"`
}

func ttl(
	ttl time.Duration,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ttl)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

type Bridger struct {
	sync.Mutex
	bridges map[string]*Bridge
}

func NewBridger() Bridger {
	return Bridger{
		bridges: make(map[string]*Bridge),
	}
}

func (b *Bridger) Get(id string) *Bridge {
	b.Lock()
	defer b.Unlock()
	if bridge, ok := b.bridges[id]; ok {
		return bridge
	}
	bridge := newBridge(func() {
		b.Del(id)
	})
	b.bridges[id] = bridge
	return bridge
}

func (b *Bridger) Del(id string) {
	b.Lock()
	defer b.Unlock()
	delete(b.bridges, id)
}

type Bridge struct {
	clean func()
	l     *Conn
	r     *Conn
}

func newBridge(clean func()) *Bridge {
	l := make(chan http.ResponseWriter, 1)
	r := make(chan http.ResponseWriter, 1)
	b := &Bridge{
		clean: clean,
		l:     newConn(r, l),
		r:     newConn(l, r),
	}
	return b
}

func (b *Bridge) Join(w http.ResponseWriter) (*Conn, error) {
	switch {
	case b.l.TryLock():

		b.r.receive <- w

		return b.l, nil
	case b.r.TryLock():

		b.l.receive <- w

		return b.r, nil
	default:
		return nil, fmt.Errorf("no available connections")
	}
}

func (b *Bridge) Wait() {
	defer b.clean()
	<-b.l.ctx.Done()
	<-b.r.ctx.Done()
}

type Conn struct {
	sync.Mutex
	ctx     context.Context
	cancel  func()
	send    chan http.ResponseWriter
	receive chan http.ResponseWriter
}

func newConn(send, rec chan http.ResponseWriter) *Conn {
	ctx, cancel := context.WithCancel(context.Background())
	return &Conn{
		ctx:     ctx,
		cancel:  cancel,
		send:    send,
		receive: rec,
	}
}

func (c *Conn) Receive() chan http.ResponseWriter {
	return c.receive
}

func (c *Conn) Close() {
	defer c.Unlock()
	close(c.send)
	c.cancel()
}
