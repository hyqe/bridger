package bridging

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/hyqe/bridger/internal/mint"
)

/*
curl -X POST "http://localhost:8080/bridges" \
	-d '{"users":["a","b"]}' \
	-i

*/

type CreateRequest struct {
	With string
}

type CreateResponse struct {
	BridgeId  string    `json:"bridgeId"`
	ExpiresAt time.Time `json:"expiresAt"`
	Token     string    `json:"token"`
}

func NewCreateHandler(
	getSecret func() []byte,
	getUserId func(r *http.Request) string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := getUserId(r)

		var form CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claim := Claim{
			BridgeId:  hashedBridgeId(userId, form.With),
			UserId:    userId,
			ExpiresAt: time.Now().UTC().Add(time.Hour),
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateResponse{
			BridgeId:  claim.BridgeId,
			ExpiresAt: claim.ExpiresAt,
			Token: mint.NewToken(
				getSecret(),
				claim,
			).String(),
		})
	}
}

func hashedBridgeId(users ...string) string {
	sort.Strings(users)
	return hex.Dump(md5.New().Sum([]byte(strings.Join(users, ""))))
}
