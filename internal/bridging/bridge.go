package bridging

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Claim struct {
	BridgeId  string    `json:"i"`
	UserId    string    `json:"u"`
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

	b.bridges[id] = bridge
	return bridge
}

type Bridge struct {
	sync.WaitGroup
	ctx context.Context
	l   *Conn
	r   *Conn
}

func NewBridge() *Bridge {
	return &Bridge{}

}

func (b *Bridge) Join(w http.ResponseWriter) *Conn {

	return conn
}
func (b *Bridge) Wait() {
	defer b.clean()
	b.WaitGroup.Wait()
}

type Conn struct {
	sync.Mutex
	send    chan http.ResponseWriter
	receive chan http.ResponseWriter
}

func (b *Conn) Receive() chan http.ResponseWriter {
	return nil
}

func (b *Conn) Close() {

}
