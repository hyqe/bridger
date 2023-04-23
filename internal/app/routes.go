package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes(
	r *mux.Router,
	JoinBridge http.HandlerFunc,
) {
	r.HandleFunc("/bridges/{bridgeId}", JoinBridge).Methods(http.MethodGet, http.MethodPut)
}

func getBridgeId(r *http.Request) string {
	return mux.Vars(r)["bridgeId"]
}
