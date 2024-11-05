package web

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	start := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '/' && i == start {
			start = i + 1
			continue
		}
		if pattern[i] == '/' && i > start {
			item := pattern[start:i]
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
			start = i + 1
		}
	}
	if start < len(pattern) {
		parts = append(parts, pattern[start:])
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	pattern = "/" + strings.Join(parts, "/")
	log.Printf("[Route] %4s - %s", method, pattern)
	if _, has := r.roots[method]; !has {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)

	if _, has := r.handlers[method]; !has {
		r.handlers[method] = make(map[string]HandlerFunc)
	}
	r.handlers[method][pattern] = handler

}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		r.handlers[c.Method][n.pattern](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
