package main

import (
	"net/http"
	"time"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/router"
)

func main() {
	time := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().UTC().String()))
	})

	r := router.New().
		AddRoute("/time", http.MethodGet, time)

	users := map[string]string{
		"jonty": "shsgfhjdsf",
		"terry": "123bmsda!",
	}

	h := withBasicAuth(users, r)

	http.ListenAndServe(":8773", h)
}

func withBasicAuth(userNameToPassword map[string]string, next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "basic auth not used", http.StatusUnauthorized)
			return
		}
		if pp, ok := userNameToPassword[u]; !ok || p != pp {
			http.Error(w, "unkown user or invalid password", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
