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

// The successfully matched node
func (n *node) matchChild(part string) *node {
	if _, has := n.children[part]; has {
		return n.children[part]
	}
	if n.wildChild != nil {
		return n.wildChild
	}
	return nil
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
			break
		} else {
			cur.children[part] = next
		}
		nodePath = append(nodePath, next.part)
		cur = next
	}
	if cur.isUse {
		pattern, oldPattern := strings.Join(parts, "/"), strings.Join(nodePath, "/")
		panic(fmt.Errorf("Find the route \"%s\" is in conflict with \"%s\". It may not be safe. ", oldPattern, pattern))
	}
	cur.handle = handlerFunc
	cur.isUse = true
}

func (t *trieTree) search(parts []string) (*node, map[string]string) {
	cur := t.root
	params := make(map[string]string)
	for i, part := range parts {
		if cur.isStop {
			params[cur.part[1:]] = strings.Join(parts[i-1:], "/")
			break
		}
		next := cur.matchChild(part)
		if next == nil {
			return nil, nil
		}
		if len(next.part) > 0 && next.part[0] == ':' {
			params[cur.part[1:]] = part
		}
		cur = next
	}
	return cur, params
}
