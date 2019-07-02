package router

import (
	"net/http"
	"net/url"
	"strings"
)

// route is the handler and middleware for a given http verb on a path
type route struct {
	handler    http.HandlerFunc
	middleware []http.HandlerFunc
}

// node refers to a path in our route tree, contains handlers for each method
// and middleware. unique refers to if the node takes in parameters e.g :id
type node struct {
	prefix   string
	unique   bool
	children []*node
	methods  map[string]*route
}

func (n *node) add(method, path string, handle http.HandlerFunc) {

	// split the path at every "/", ignoring the first "/"
	paths := strings.Split(path, "/")[1:]
	count := len(paths)

	for {
		current, path := n.walk(paths, nil)

		// if existing node is being updated
		if current.prefix == path && count == 1 {
			r := route{handler: handle}
			current.methods[method] = &r
			return
		}

		newNode := node{
			prefix:  path,
			unique:  false,
			methods: make(map[string]*route),
		}

		// if this is a unique parameter
		if len(path) > 0 && path[0] == ':' {
			newNode.unique = true
		}

		// reached end of path, add handler to given method
		if count == 1 {
			r := route{handler: handle}
			newNode.methods[method] = &r
		}

		current.children = append(current.children, &newNode)
		count--
		if count == 0 {
			break
		}

	}

}

// walk will move along the tree parsing unique parameters
// returns the node and path found
func (n *node) walk(paths []string, params url.Values) (*node, string) {
	path := paths[0]

	if len(n.children) > 0 {
		for _, child := range n.children {
			if path == child.prefix || child.unique {
				if child.unique && params != nil {
					params.Add(child.prefix[1:], path)
				}

				next := paths[1:]
				if len(next) > 0 {
					return child.walk(next, params)
				}
				return child, path
			}
		}
	}

	return n, path
}
