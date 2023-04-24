package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(
	r *mux.Router,
	createBridge http.HandlerFunc,
	joinBridge http.HandlerFunc,
) {
	r.HandleFunc("/bridges", createBridge).Methods(http.MethodPost)
	r.HandleFunc("/bridges/{bridgeId}", joinBridge).Methods(http.MethodGet, http.MethodPut)
}

func getBridgeId(r *http.Request) string {
	return mux.Vars(r)["bridgeId"]
}
