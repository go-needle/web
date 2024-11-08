package web

import (
	"net/http"
)

type store interface {
	insert(parts []string, handlerFunc HandlerFunc) int
	search(parts []string) (*node, map[string]string)
}

type router struct {
	store map[string]store
	total int
}

func newRouter() *router {
	return &router{
		store: make(map[string]store),
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
	if _, has := r.store[method]; !has {
		r.store[method] = newTrieTree()
	}
	r.total += r.store[method].insert(parts, handler)
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	tree, ok := r.store[method]
	if !ok {
		return nil, nil
	}
	return tree.search(searchParts)
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		c.handlers = append(c.handlers, n.handle)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
