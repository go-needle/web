package web

import (
	"fmt"
	"strings"
)

type node struct {
	handle    HandlerFunc
	part      string
	isUse     bool
	children  map[string]*node
	wildChild *node //  check ':'
	isStop    bool  // stop at '*'
}

func newNode(part string) *node {
	return &node{nil, part, false, map[string]*node{}, nil, false}
}

// The successfully matched node for insertion
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
	nodes := make([]*node, 0, 2)
	if _, has := n.children[part]; has {
		nodes = append(nodes, n.children[part])
	}
	if n.wildChild != nil {
		nodes = append(nodes, n.wildChild)
	}
	return nodes
}

type trieTree struct {
	root *node
}

func newTrieTree() *trieTree {
	return &trieTree{newNode("")}
}

func (t *trieTree) insert(parts []string, handlerFunc HandlerFunc) {
	cur := t.root
	var nodePath []string
	for _, part := range parts {
		next := cur.matchChild(part)
		if next == nil {
			next = newNode(part)
		}
		if part[0] == ':' {
			cur.wildChild = next
		} else if part[0] == '*' {
			next.isStop = true
			cur.wildChild = next
		} else {
			cur.children[part] = next
		}
		nodePath = append(nodePath, next.part)
		if cur.wildChild != nil && cur.wildChild.part[0] == '*' && len(cur.children) != 0 {
			panic(fmt.Errorf("A conflict has occurred at route \"/%s\"", strings.Join(parts, "/")))
		}
		cur = next
		if cur.isStop {
			break
		}
	}
	if cur.isUse {
		panic(fmt.Errorf("A conflict has occurred at route \"/%s\"", strings.Join(parts, "/")))
	}
	cur.handle = handlerFunc
	cur.isUse = true
}

func (t *trieTree) search(parts []string) (*node, map[string]string) {
	cur := t.root
	params := make(map[string]string)
	for i, part := range parts {
		next := cur.matchChild(part)
		if next == nil {
			return nil, nil
		}
		if len(next.part) > 0 && next.part[0] == ':' {
			params[next.part[1:]] = part
		}
		cur = next
		if cur.isStop {
			params[cur.part[1:]] = strings.Join(parts[i-1:], "/")
			break
		}
	}
	if !cur.isUse {
		return nil, nil
	}
	return cur, params
}
