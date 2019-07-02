package router

import (
	"errors"
	"net/http"
	"strings"
)

// Router is URL router to route desired requests to correct handler
type Router struct {
	trie *node
}

func NewRouter() *Router {

	n := &node{
		prefix:  "/",
		unique:  false,
		methods: make(map[string]*route),
	}

	return &Router{
		trie: n,
	}
}

// Handle will map a handler to a given method on a path
func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if path[0] != '/' {
		panic("Not a correct path. Must start with /.")
	}

	r.trie.add(method, path, handler)
}

// ServeHTTP to satisfy mux interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	handler, err := r.Match(req)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	handler(w, req)
}

// Match walks the path tree and finds handler for approriate request
func (r *Router) Match(req *http.Request) (http.HandlerFunc, error) {
	req.ParseForm()
	params := req.Form

	path := strings.Split(req.URL.Path, "/")[1:]

	node, _ := r.trie.walk(path, params)

	route := node.methods[req.Method]
	if route != nil {
		return route.handler, nil
	}

	return nil, errors.New("No such handler")
}

// GET Handle
func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle("GET", path, handler)
}

// POST Handle
func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle("POST", path, handler)
}

// PUT Handle
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle("PUT", path, handler)
}

// DELETE Handle
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Handle("DELETE", path, handler)
}

// HEAD Handle
func (r *Router) HEAD(path string, handler http.HandlerFunc) {
	r.Handle("HEAD", path, handler)
}

// PATCH Handle
func (r *Router) PATCH(path string, handler http.HandlerFunc) {
	r.Handle("PATCH", path, handler)
}
