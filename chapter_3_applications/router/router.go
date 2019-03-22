package router

import "net/http"

// New creates a new Router.
func New() *Router {
	return &Router{
		routes: make(map[route]http.Handler),
	}
}

// Router for HTTP requests.
type Router struct {
	routes map[route]http.Handler
}

// AddRoute to the Router.
func (rtr *Router) AddRoute(path string, method string, handler http.Handler) *Router {
	rtr.routes[route{Path: path, Method: method}] = handler
	return rtr
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	incoming := route{Path: r.URL.Path, Method: r.Method}
	child, exists := rtr.routes[incoming]
	if !exists {
		http.NotFound(w, r)
		return
	}
	child.ServeHTTP(w, r)
}

type route struct {
	Path   string
	Method string
}
