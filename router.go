package web

import (
	"fmt"
	"net/http"
)

type router struct {
	tree  map[string]*trieTreeR
	total int
}

func newRouter() *router {
	return &router{
		tree: make(map[string]*trieTreeR),
	}
}

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

func (r *router) addRoute(method string, pattern string, handler Handler) {
	parts := parsePattern(pattern)
	if _, has := r.tree[method]; !has {
		r.tree[method] = newTrieTreeR()
	}
	r.total += r.tree[method].insert(parts, handler)
}

func (r *router) getRoute(method string, path string) (*nodeR, map[string]string) {
	searchParts := parsePattern(path)
	tree, ok := r.tree[method]
	if !ok {
		return nil, nil
	}
	return tree.search(searchParts)
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.params = params
		c.handlers = append(c.handlers, n.handler)
	} else {
		c.handlers = append(c.handlers, HandlerFunc(func(c *Context) {
			c.Fail(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s", c.Path))
		}))
	}
	c.Next()
}
