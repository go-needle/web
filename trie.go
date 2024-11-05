package web

import (
	"strings"
)

type node struct {
	pattern      string
	part         string
	children     map[string]*node
	wildChildren []*node
}

// The first successfully matched node for insertion
func (n *node) matchChild(part string) *node {
	if _, has := n.children[part]; has {
		return n.children[part]
	}
	if len(n.wildChildren) > 0 {
		return n.wildChildren[0]
	}
	return nil
}

// All successfully matched nodes are used to find
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	if _, has := n.children[part]; has {
		nodes = append(nodes, n.children[part])
	}
	nodes = append(nodes, n.wildChildren...)
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part}
		if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*") {
			n.wildChildren = append(n.wildChildren, child)
		} else {
			if n.children == nil {
				n.children = make(map[string]*node)
			}
			n.children[part] = child
		}
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
