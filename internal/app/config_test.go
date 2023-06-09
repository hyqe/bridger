package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyqe/assert"
)

func Test_spam_StatusTeapot(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest(http.MethodGet, "/foo/.env", nil)
	w := httptest.NewRecorder()
	spam(next).ServeHTTP(w, r)
	assert.Want(t, http.StatusTeapot, w.Code)
}

func Test_spam_next(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r := httptest.NewRequest(http.MethodGet, "/foo/", nil)
	w := httptest.NewRecorder()
	spam(next).ServeHTTP(w, r)
	assert.Want(t, http.StatusOK, w.Code)
}

func Test_getConnectingIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/foo/", nil)
	r.Header.Set("DO-Connecting-IP", "11.0.0.1")
	ip, ok := getConnectingIP(r)
	assert.Want(t, true, ok)
	assert.Want(t, "11.0.0.1", ip)
}
