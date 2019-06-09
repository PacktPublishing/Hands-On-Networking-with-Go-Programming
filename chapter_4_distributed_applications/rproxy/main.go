package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/didip/tollbooth"

	"github.com/didip/tollbooth/limiter"
)

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

func withAPIKey(apiKeyToUser map[string]string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := apiKeyToUser[r.Header.Get("X-API-Key")]; !ok {
			http.Error(w, "missing or invalid X-API-Key header", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type limitingMiddleware struct {
	limiter *limiter.Limiter
	handler http.Handler
}

func (lm limitingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := "X-API-Key"
	lm.limiter.SetHeader(header, []string{r.Header.Get(header)})
	lm.handler.ServeHTTP(w, r)
}

func withRateLimiting(next http.Handler) http.Handler {
	perSecond := 1.0
	l := tollbooth.NewLimiter(perSecond, nil)
	return limitingMiddleware{
		limiter: l,
		handler: tollbooth.LimitHandler(l, next),
	}
}

func urlMust(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return u
}

func main() {
	backend1 := urlMust("https://5cf6b7a146583900149cbb82.mockapi.io")
	healthCheck1 := urlMust("https://5cf6b7a146583900149cbb82.mockapi.io/health")
	backend2 := urlMust("http://localhost:1111")
	healthCheck2 := urlMust("http://localhost:1111/health")
	// rp := newReverseProxy(backend1)
	rp, cancel := NewBalancer(context.Background(), time.Second*5,
		NewOriginServer(backend1, healthCheck1),
		NewOriginServer(backend2, healthCheck2))

	cached := NewCache(time.Minute, rp)

	limited := withRateLimiting(cached)

	apiKeyToUsername := map[string]string{
		"9a5f509f-b81b-4e91-a728-1f5d55adbd55": "user1",
		"829d5c8c-a782-4bbf-af00-c50044111d0e": "user2",
	}
	authenticated := withAPIKey(apiKeyToUsername, limited)

	allowedMethods := []string{
		"GET", "POST", // OPTIONS is often required, but I'm leaving it out.
	}
	allowedQueryKeys := []string{"q"} // Just allow one querystring param.
	maxQueryValueLength := 10         // Only small querystring values are allowed.
	filtered := withFiltering(allowedMethods, allowedQueryKeys, maxQueryValueLength, authenticated)

	err := http.ListenAndServe(":9898", filtered)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cancel()
}

func withFiltering(allowMethods []string, allowQueryKeys []string, maxQueryValueLength int, next http.Handler) http.Handler {
	// Put the items into a map for faster lookup.
	ms := make(map[string]struct{})
	for _, am := range allowMethods {
		ms[am] = struct{}{}
	}
	aq := make(map[string]struct{})
	for _, k := range allowQueryKeys {
		aq[k] = struct{}{}
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		if _, methodAllowed := ms[r.Method]; !methodAllowed {
			dumped, _ := httputil.DumpRequest(r, true)
			log.Printf("filtered request (invalid HTTP method): %s", string(dumped))
			http.Error(w, "invalid HTTP method", http.StatusMethodNotAllowed)
			return
		}
		for k, v := range r.URL.Query() {
			if _, queryParamAllowed := aq[k]; !queryParamAllowed {
				dumped, _ := httputil.DumpRequest(r, true)
				log.Printf("filtered request (invalid query param): %s", string(dumped))
				http.Error(w, "invalid query param", http.StatusUnprocessableEntity)
				return
			}
			if len(v) > maxQueryValueLength {
				dumped, _ := httputil.DumpRequest(r, true)
				log.Printf("filtered request (query param length): %s", string(dumped))
				http.Error(w, "invalid query param length", http.StatusUnprocessableEntity)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func newReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
