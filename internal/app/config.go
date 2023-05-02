package app

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hyqe/bridger/internal/bridging"
	"github.com/hyqe/bridger/internal/mint"
	"github.com/hyqe/timber"
)

const (
	envPORT = "PORT"
)

// Service builds the entire service.
func Service() (address string, handler http.Handler) {
	r := mux.NewRouter()

	jack := timber.NewJack()

	Routes(
		r,
		bridging.NewCreateHandler(getSecret, getUserId),
		bridging.NewJoinHandler(getSecret, getBridgeId, getClaim[bridging.Claim]),
	)

	logger := timber.NewMiddleware(jack)

	return addr(), logger(spam(r))
}

// addr gets address to bind the server to.
func addr() string {
	port, ok := os.LookupEnv(envPORT)
	if !ok {
		return ":8080"
	}
	return fmt.Sprintf(":%v", port)
}

var secret = mint.NewSecret(32)

func getSecret() []byte {
	SECRET, ok := os.LookupEnv("SECRET")
	if ok {
		return []byte(SECRET)
	}
	return secret
}

func getClaim[T any](r *http.Request) (T, error) {
	var claim T
	rawtoken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	token, err := mint.ParseToken(rawtoken)
	if err != nil {
		return claim, err
	}
	if !token.IsValid(getSecret()) {
		return claim, fmt.Errorf("invalid token")
	}
	token.Into(&claim)
	return claim, err
}

func getUserId(r *http.Request) string {
	return r.URL.Query().Get("from")
}

func spam(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch path.Base(r.URL.Path) {
		case
			".env",
			".remote",
			".local",
			".production",
			"wp-login.php",
			"index.php",
			"administrator":
			http.Error(w, "💩", http.StatusTeapot)
		default:
			next.ServeHTTP(w, r)
		}
	}
}
