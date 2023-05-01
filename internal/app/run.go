package app

import "net/http"

func Run() {
	http.ListenAndServe(Service())
}
