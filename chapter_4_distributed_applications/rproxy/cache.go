package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

func NewCache(duration time.Duration, next http.Handler) *Cache {
	return &Cache{
		urlToResponse: make(map[string]CacheResponse),
		Next:          next,
		Duration:      duration,
	}
}

type Cache struct {
	urlToResponse map[string]CacheResponse
	Next          http.Handler
	Duration      time.Duration
}

type CacheResponse struct {
	Header http.Header
	Body   []byte
	Expiry time.Time
}

func (c *Cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodOptions || r.Method == http.MethodHead {
		cr, inCache := c.urlToResponse[r.URL.String()]
		if inCache && time.Now().Before(cr.Expiry) {
			log.Printf("%s: cache hit\n", r.URL.String())
		} else {
			log.Printf("%s: setting up cache\n", r.URL.String())
			// Make the response and cache it.
			ww := httptest.NewRecorder()
			c.Next.ServeHTTP(ww, r)

			body, _ := ioutil.ReadAll(ww.Body)
			cr = CacheResponse{
				Body:   body,
				Header: ww.Header(),
				Expiry: time.Now().Add(c.Duration),
			}
			c.urlToResponse[r.URL.String()] = cr
		}
		// Copy headers.
		for k, v := range cr.Header {
			for _, vv := range v {
				w.Header().Set(k, vv)
			}
		}
		// Write body.
		if cr.Body != nil {
			w.Write(cr.Body)
		}
		return
	}
	c.Next.ServeHTTP(w, r)
}
