package rest

import (
	"fmt"
	"net/http"
)

func Serve(handler http.HandlerFunc) {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.String())
		handler(w, r)
	}))
}
