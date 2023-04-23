package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hyqe/bridger/internal/bridging"
)

const (
	envPORT = "PORT"
)

func Run() {
	http.ListenAndServe(Addr(), Service())
}

func Addr() string {
	port, ok := os.LookupEnv(envPORT)
	if !ok {
		return ":8080"
	}
	return fmt.Sprintf(":%v", port)
}

func Service() http.Handler {
	r := mux.NewRouter()

	Routes(
		r,
		bridging.NewJoinHandler(getBridgeId),
	)

	return r
}
