package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots   map[string]*node
	handles map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		roots:   make(map[string]*node),
		handles: make(map[string]HandleFunc)
	}
}

// 解析路径pattern
func parsePattern(pattern string) []string  {
	patternSplit := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range patternSplit{
		if item != ""{
			parts = append(parts, item)
			if item[0] == '*'{
				break
			}
		}
	}
	return parts
}


func (r *router) addRoute(method string, pattern string, handle HandleFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok{
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handles[key] = handle
}

func (r *router) getRouter(method string, path string) (*node, map[string]string)  {
	searchParts := parsePattern()
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handles[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
