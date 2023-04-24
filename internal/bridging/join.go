package bridging

import (
	"io"
	"net/http"
	"time"
)

func NewJoinHandler(
	getBridgeId func(r *http.Request) string,
	getBridgeClaim func(r *http.Request) Claim,
) http.HandlerFunc {
	bridger := NewBridger()
	return ttl(time.Minute, func(w http.ResponseWriter, r *http.Request) {
		id := getBridgeId(r)
		claim := getBridgeClaim(r)

		if id != claim.BridgeId {
			http.Error(w, "invalid bridge id", http.StatusBadRequest)
			return
		}

		if time.Now().After(claim.ExpiresAt) {
			http.Error(w, "bridge has expired", http.StatusBadRequest)
			return
		}

		bridge := bridger.Join(claim, w)

		select {
		case remote := <-bridge.Receive():
			io.Copy(remote.w, r.Body)
			remote.Close()
		}

		bridge.Wait()
	})
}
