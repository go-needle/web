package web

import (
	"fmt"
	"log"
	"strings"
)

type node struct {
	handle    HandlerFunc
	children  map[string]*node
	jumpChild *node //  ':'
	stopChild *node // '*'
	keys      map[int]string
}

func newNode() *node {
	return &node{children: make(map[string]*node)}
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
	nodes := make([]*node, 0, 3)
	if _, has := n.children[part]; has {
		nodes = append(nodes, n.children[part])
	}
	if n.jumpChild != nil {
		nodes = append(nodes, n.jumpChild)
	}
	return nodes
}

type trieTree struct {
	root            *node
	heightNodeCount map[int]int
}

func newTrieTree() *trieTree {
	return &trieTree{newNode(), make(map[int]int)}
}

func (t *trieTree) insert(parts []string, handlerFunc HandlerFunc) {
	cur := t.root
	keys := make(map[int]string)
	height := 0
	for i, part := range parts {
		next := cur.matchChild(part)
		height++
		if (part[0] == ':' || part[0] == '*') && len(part) == 1 {
			panic(fmt.Errorf("the routing path \"%s\" cannot contain nodes with only \"*\" or \":\"", strings.Join(parts, "/")))
		}
		if part[0] == '*' {
			keys[i] = part[1:]
			if next == nil {
				next = newNode()
				cur.stopChild = next
				cur = next
				break
			}
		} else if part[0] == ':' {
			keys[i] = part[1:]
			if next == nil || cur.stopChild == next {
				next = newNode()
				cur.jumpChild = next
			}
		} else {
			if next == nil || cur.jumpChild == next || cur.stopChild == next {
				next = newNode()
				cur.children[part] = next
			}
		}
		cur = next
	}
	if cur.handle != nil {
		t.heightNodeCount[height]--
		log.Printf("[Warning] A route coverage occurred in \"/%s\"", strings.Join(parts, "/"))
	}
	cur.handle = handlerFunc
	cur.keys = keys
	t.heightNodeCount[height]++
}

func (t *trieTree) search(parts []string) (*node, map[string]string) {
	queue := []*node{t.root}
	cur := []*node{t.root}
	height := 0
	var stopNodes []*node
	for len(queue) > 0 && height < len(parts) {
		tmp := make([]*node, 0, t.heightNodeCount[height])
		part := parts[height]
		for len(queue) > 0 {
			head := queue[0]
			if head.stopChild != nil {
				stopNodes = append(stopNodes, head.stopChild)
			}
			nxt := head.matchChildren(part)
			queue = queue[1:]
			tmp = append(tmp, nxt...)
		}
		queue = append(queue, tmp...)
		cur = tmp
		height++
	}

	params := make(map[string]string)
	var nd *node
	isStop := false

	for _, n := range cur {
		if n.handle != nil {
			nd = n
			break
		}
	}

	if nd == nil {
		for i := len(stopNodes) - 1; i >= 0; i-- {
			if stopNodes[i].handle != nil {
				nd = stopNodes[i]
				isStop = true
				break
			}
		}
	}

	if nd == nil {
		return nil, nil
	}

	if isStop {
		maxIdx := -1
		for k, v := range nd.keys {
			params[v] = parts[k]
			maxIdx = max(maxIdx, k)
		}
		params[nd.keys[maxIdx]] = strings.Join(parts[maxIdx:], "/")
	} else {
		for k, v := range nd.keys {
			params[v] = parts[k]
		}
	}
	return nd, params
}
