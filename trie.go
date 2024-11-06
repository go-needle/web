package web

import (
	"fmt"
	"strings"
)

type node struct {
	handle    HandlerFunc
	children  map[string]*node
	jumpChild *node //  ':'
	stopChild *node // '*'
	keys      map[string]int
}

// The successfully matched node for insertion
func (n *node) matchChild(part string) *node {
	if _, has := n.children[part]; has {
		return n.children[part]
	}
	if n.jumpChild != nil {
		return n.jumpChild
	}
	if n.stopChild != nil {
		return n.stopChild
	}
	return nil
}

// All successfully matched nodes are used to find
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0, 2)
	if _, has := n.children[part]; has {
		nodes = append(nodes, n.children[part])
	}
	if n.jumpChild != nil {
		nodes = append(nodes, n.jumpChild)
	}
	if n.stopChild != nil {
		nodes = append(nodes, n.stopChild)
	}
	return nodes
}

type trieTree struct {
	root *node
}

func newTrieTree() *trieTree {
	return &trieTree{&node{}}
}

func (t *trieTree) insert(parts []string, handlerFunc HandlerFunc) {
	cur := t.root
	keys := make(map[string]int)
	for i, part := range parts {
		next := cur.matchChild(part)
		if next == nil {
			next = &node{}
		}
		if (part[0] == ':' || part[0] == '*') && len(part) == 1 {
			panic(fmt.Errorf("the routing path \"%s\" cannot contain nodes with only \"*\" or \":\"", strings.Join(parts, "/")))
		}
		if part[0] == ':' {
			keys[part[1:]] = i
			cur.jumpChild = next
		} else if part[0] == '*' {
			cur.stopChild = next
		} else {
			cur.children[part] = next
		}
		if cur.stopChild != nil && (len(cur.children) != 0 || cur.jumpChild != nil) {
			panic(fmt.Errorf("a conflict has occurred at route \"/%s\"", strings.Join(parts, "/")))
		}
		if cur.stopChild != nil {
			cur = next
			break
		}
		cur = next
	}
	if cur.handle != nil {
		panic(fmt.Errorf("a conflict has occurred at route \"/%s\"", strings.Join(parts, "/")))
	}
	cur.handle = handlerFunc
	cur.keys = keys
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
