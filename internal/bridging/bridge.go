package bridging

import (
	"context"
	"fmt"
	"io"
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
		
	}
}

func (b *Bridger) Join(claim Claim, w http.ResponseWriter) Bridge {
	return Bridge{
		
	}
}

type Bridge struct {
	clean func()
	sync.Mutex
	sync.WaitGroup
	send http.ResponseWriter
	receive chan http.ResponseWriter
}

func (b *Bridge) Receive() chan  {

}

func (b *Bridge) Wait() {
	defer clean()
}
