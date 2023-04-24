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
	r.HandleFunc("/", health).Methods(http.MethodGet)
	r.HandleFunc("/bridges", createBridge).Methods(http.MethodGet)
	r.HandleFunc("/bridges/{bridgeId}", joinBridge).Methods(http.MethodGet, http.MethodPut, http.MethodPost)
}

func getBridgeId(r *http.Request) string {
	return mux.Vars(r)["bridgeId"]
}

func health(w http.ResponseWriter, r *http.Request) {

}
