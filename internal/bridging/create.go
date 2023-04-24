package bridging

import (
	"encoding/json"
	"net/http"
	"time"
)

type CreateRequest struct {
	Users []string
}

type CreateResponse struct {
	BridgeId  string    `json:"bridgeId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var form CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: implement tokening into BridgeId here.
		// The token should be very minimal, contains only
		// whats need to auth a bridge join.

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateResponse{
			BridgeId:  "TODO",
			ExpiresAt: time.Now().UTC().Add(time.Hour),
		})
	}
}
