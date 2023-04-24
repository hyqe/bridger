package bridging

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/hyqe/bridger/internal/mint"
)

type CreateRequest struct {
	Users []string
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
		// userId := getUserId(r)

		var form CreateRequest

		//if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//	return
		//}

		form.Users = r.URL.Query()["u"]

		if IsEmpty(form.Users...) {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		claim := Claim{
			BridgeId:  hashedBridgeId(form.Users...),
			ExpiresAt: time.Now().UTC().Add(time.Hour),
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, mint.NewToken(
			getSecret(),
			claim,
		).String())

		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(CreateResponse{
		// 	BridgeId:  claim.BridgeId,
		// 	ExpiresAt: claim.ExpiresAt,
		// 	Token: mint.NewToken(
		// 		getSecret(),
		// 		claim,
		// 	).String(),
		// })
	}
}

func hashedBridgeId(users ...string) string {
	sort.Strings(users)
	return hex.EncodeToString(md5.New().Sum([]byte(strings.Join(users, ""))))
}

func IsEmpty(vs ...string) bool {
	for _, v := range vs {
		if strings.TrimSpace(v) == "" {
			return true
		}
	}
	return false
}
