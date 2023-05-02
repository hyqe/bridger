package app

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

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
	var mu sync.Mutex
	IPs := make(map[string]time.Time)

	clean := func(ip string) {
		mu.Lock()
		defer mu.Unlock()
		delete(IPs, ip)
	}

	add := func(ip string) {
		mu.Lock()
		defer mu.Unlock()
		IPs[ip] = time.Now()
	}

	hasAny := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(IPs) > 0
	}

	has := func(ip string) bool {
		mu.Lock()
		defer mu.Unlock()
		_, ok := IPs[ip]
		return ok
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			if !hasAny() {
				continue
			}
			for k, v := range IPs {
				if time.Now().After(v.Add(time.Minute)) {
					clean(k)
				}
			}
		}
	}()

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
			if ip, ok := getConnectingIP(r); ok {
				add(ip)
				w.Header().Set("X-Warning", "Your actions are being reported 👀")
			}
			w.Header().Set("Content-Type", "💩")
			w.WriteHeader(http.StatusTeapot)
			fmt.Fprint(w, "💩")
		default:
			if ip, ok := getConnectingIP(r); ok && has(ip) {
				http.Error(w, "", http.StatusTeapot)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

var reHeaderConnectingIP = regexp.MustCompile(`^.*(?i:connect.*[-]ip)$`)

func getConnectingIP(r *http.Request) (ip string, ok bool) {
	for k, v := range r.Header {
		if reHeaderConnectingIP.MatchString(k) {
			if len(v) > 0 {
				return v[0], true
			}
		}
	}
	if strings.TrimSpace(r.RemoteAddr) != "" {
		return r.RemoteAddr, true
	}
	return "", false
}
