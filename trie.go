package web

import (
	"fmt"
	"log"
)

type node struct {
	pattern   string
	isUse     bool
	children  map[string]*node
	wildChild *node
	isStop    bool // stop at '*'
}

// The first successfully matched node for insertion
func (n *node) matchChild(part string) *node {
	if _, has := n.children[part]; has {
		return n.children[part]
	}
	if n.wildChild != nil {
		return n.wildChild
	}
	return nil
}

// All successfully matched nodes are used to find
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	if _, has := n.children[part]; has {
		nodes = append(nodes, n.children[part])
	}
	if n.wildChild != nil {
		nodes = append(nodes, n.wildChild)
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height || n.isStop {
		isSame := n.pattern == pattern
		if n.isUse && !isSame {
			panic(fmt.Errorf("Find the route \"%s\" is in conflict with \"%s\". It may not be safe. ", n.pattern, pattern))
		}
		if n.isUse {
			log.Printf("[Warning] the route \"%s\" is covered", n.pattern)
		}
		n.pattern = pattern
		n.isUse = true
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{}
		if part[0] == ':' {
			n.wildChild = child
		} else if part[0] == '*' {
			child.isStop = true
			n.wildChild = child
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
	if len(parts) == height || n.isStop {
		if !n.isUse {
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
