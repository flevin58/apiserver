// Implementing a basic router using only native net/http golang functionality
// The purpose is to show how things work under the hood without using external
// packages like gorilla, gin gonic, etc. for a simple usage.
package router

import (
	"net/http"
	"strings"
)

type Route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

type Router struct {
	routes []Route
}

func New() *Router {
	return &Router{
		routes: []Route{},
	}
}

func (rr *Router) AddRoute(path, method string, handler http.HandlerFunc) {
	route := Route{
		path:    path,
		method:  method,
		handler: handler,
	}
	rr.routes = append(rr.routes, route)
}

func (rr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rr.routes {
		// Check if the method is valid for this route
		if r.Method != route.method {
			continue
		}

		// Check if the path is valid for this route
		if r.URL.Path == route.path {
			route.handler(w, r)
			return
		}

		// Check if this route ends with a '/' and if so, pass the remaining path to var "id"
		// Waiting for the new version of go when patterns will be implemented !!!
		if strings.HasSuffix(route.path, "/") && strings.HasPrefix(r.URL.Path, route.path) {
			value := r.URL.Path[len(route.path):]
			r.ParseForm()
			r.Form.Add("id", value)
			route.handler(w, r)
			return
		}
	}
	http.Error(w, "bad route: "+r.URL.Path, http.StatusBadRequest)
}
