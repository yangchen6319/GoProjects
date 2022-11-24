package gee

import (
	"log"
	"net/http"
)

type router struct {
	handles map[string]HandleFunc
}

func newRouter() *router {
	return &router{handles: make(map[string]HandleFunc)}
}

func (r *router) addRoute(method string, pattern string, handle HandleFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handles[key] = handle
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handles[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
