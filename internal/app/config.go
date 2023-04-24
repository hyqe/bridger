package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hyqe/bridger/internal/bridging"
	"github.com/hyqe/bridger/internal/mint"
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

	Routes(
		r,
		bridging.NewCreateHandler(secret, getUserId),
		bridging.NewJoinHandler(getBridgeId, getBridgeClaim),
	)

	return r
}

func secret() []byte {
	return mint.NewSecret(32)
}

func getBridgeClaim(r *http.Request) bridging.Claim {
	return bridging.Claim{}
}

func getUserId(r *http.Request) string {
	return ""
}
