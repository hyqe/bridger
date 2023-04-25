package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hyqe/bridger/internal/bridging"
	"github.com/hyqe/bridger/internal/mint"
	"github.com/hyqe/timber"
)

const (
	envPORT = "PORT"
)

// Addr gets address to bind the server to.
func Addr() string {
	port, ok := os.LookupEnv(envPORT)
	if !ok {
		return ":8080"
	}
	return fmt.Sprintf(":%v", port)
}

// Service builds the entire service.
func Service() http.Handler {
	r := mux.NewRouter()

	jack := timber.NewJack()

	Routes(
		r,
		bridging.NewCreateHandler(getSecret, getUserId),
		bridging.NewJoinHandler(getSecret, getBridgeId, getClaim[bridging.Claim]),
	)

	logger := timber.NewMiddleware(jack)

	return logger(r)
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
