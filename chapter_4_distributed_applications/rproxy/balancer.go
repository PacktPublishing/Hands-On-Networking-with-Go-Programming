package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

func NewOriginServer(target, healthCheck *url.URL) *OriginServer {
	return &OriginServer{
		Target:              target,
		HealthCheckEndpoint: healthCheck,
	}
}

type OriginServer struct {
	Target              *url.URL
	IsHealthy           bool
	HealthCheckEndpoint *url.URL
}

type Pool []*OriginServer

func (p Pool) GetTarget() (u *url.URL, ok bool) {
	var healthy []int
	for i, t := range p {
		if t.IsHealthy {
			healthy = append(healthy, i)
		}
	}
	if len(healthy) == 0 {
		return
	}
	idx := rand.Int() % len(healthy)
	u = p[healthy[idx]].Target
	ok = true
	return
}

func (p Pool) HealthChecks() {
	var wg sync.WaitGroup
	for _, t := range p {
		wg.Add(1)
		go func(tt *OriginServer) {
			defer wg.Done()
			resp, err := http.Get(tt.HealthCheckEndpoint.String())
			tt.IsHealthy = err == nil && resp.StatusCode == http.StatusOK
			fmt.Printf("%s is healthy = %v\n", tt.HealthCheckEndpoint.String(), tt.IsHealthy)
		}(t)
	}
	wg.Wait()
}

type Balancer struct {
	Pool             Pool
	HealthCheckEvery time.Duration
	Proxy            http.Handler
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.Proxy.ServeHTTP(w, r)
}

func (b *Balancer) healthChecks(ctx context.Context) {
	timer := time.After(b.HealthCheckEvery)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer:
			b.Pool.HealthChecks()
			timer = time.After(b.HealthCheckEvery)
		}
	}
}

func NewBalancer(ctx context.Context, healthCheckEvery time.Duration, origins ...*OriginServer) (b *Balancer, cancel context.CancelFunc) {
	b = &Balancer{
		HealthCheckEvery: healthCheckEvery,
		Pool:             origins,
	}
	director := func(req *http.Request) {
		target, ok := b.Pool.GetTarget()
		if !ok {
			log.Println("no available servers, issuing 502 error!")
			return
		}
		targetQuery := target.RawQuery
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
	b.Proxy = &httputil.ReverseProxy{
		Director: director,
	}
	ctx, cancel = context.WithCancel(ctx)
	go b.healthChecks(ctx)
	return
}
