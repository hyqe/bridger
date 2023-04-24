package bridging

import (
	"io"
	"net/http"
	"time"

	"github.com/hyqe/bridger/internal/mint"
)

func NewJoinHandler(
	getSecret func() []byte,
	getBridgeId func(r *http.Request) string,
	getBridgeClaim func(r *http.Request) (Claim, error),
) http.HandlerFunc {
	bridger := NewBridger()
	return ttl(time.Minute, func(w http.ResponseWriter, r *http.Request) {
		rawtoken := getBridgeId(r)
		token, err := mint.ParseToken(rawtoken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if !token.IsValid(getSecret()) {
			http.Error(w, "invalid token", http.StatusForbidden)
			return
		}

		var claim Claim
		err = token.Into(&claim)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if time.Now().After(claim.ExpiresAt) {
			http.Error(w, "bridge has expired", http.StatusBadRequest)
			return
		}

		bridge := bridger.Get(claim.BridgeId)
		// defer bridge.Wait()

		conn, err := bridge.Join(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		//defer conn.Close()

		select {
		case <-r.Context().Done():
			return
		case remote := <-conn.Receive():
			io.Copy(remote, r.Body)
			conn.Close()
		}
		bridge.Wait()
	})
}
