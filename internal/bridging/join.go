package bridging

import (
	"net/http"
	"time"
)

func NewJoinHandler(
	getBridgeId func(r *http.Request) string,
) http.HandlerFunc {
	bridger := NewBridger()
	return ttl(time.Minute, func(w http.ResponseWriter, r *http.Request) {
		receiverId := getBridgeId(r)
		senderId := r.URL.Query().Get("from")
		to := receiverId + senderId
		from := senderId + receiverId
		bridger.Join(to, from, w, r)
	})
}
