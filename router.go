package web

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	trees map[string]*trieTree
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*trieTree),
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
	if _, has := r.trees[method]; !has {
		r.trees[method] = newTrieTree()
	}
	r.trees[method].insert(parts, handler)
	log.Printf("[Route] %4s - /%s", method, strings.Join(parts, "/"))
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	tree, ok := r.trees[method]
	if !ok {
		return nil, nil
	}
	return tree.search(searchParts)
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		n.handle(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
