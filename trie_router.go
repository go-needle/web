package web

import (
	"fmt"
	"log"
	"strings"
)

type nodeR struct {
	handle    HandlerFunc
	children  map[string]*nodeR
	jumpChild *nodeR //  ':'
	stopChild *nodeR // '*'
	keys      map[int]string
}

func newNodeR() *nodeR {
	return &nodeR{children: make(map[string]*nodeR)}
}

func (n *nodeR) matchChild(part string) *nodeR {
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

type trieTreeR struct {
	root              *nodeR
	heightNodeCount   map[int]int
	maxDenseNodeCount int
}

func newTrieTreeR() *trieTreeR {
	return &trieTreeR{newNodeR(), make(map[int]int), 1}
}

func (t *trieTreeR) insert(parts []string, handlerFunc HandlerFunc) int {
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
				next = newNodeR()
				cur.stopChild = next
				cur = next
				break
			}
		} else if part[0] == ':' {
			keys[i] = part[1:]
			if next == nil || cur.stopChild == next {
				next = newNodeR()
				cur.jumpChild = next
			}
		} else {
			if next == nil || cur.jumpChild == next || cur.stopChild == next {
				next = newNodeR()
				cur.children[part] = next
			}
		}
		cur = next
	}
	isAdd := true
	if cur.handle != nil {
		t.heightNodeCount[height]--
		isAdd = false
		log.Printf("[Warning] A route coverage occurred in \"/%s\"", strings.Join(parts, "/"))
	}
	cur.handle = handlerFunc
	cur.keys = keys
	t.heightNodeCount[height]++
	t.maxDenseNodeCount = max(t.maxDenseNodeCount, t.heightNodeCount[height])
	if isAdd {
		return 1
	} else {
		return 0
	}
}

func (t *trieTreeR) search(parts []string) (*nodeR, map[string]string) {
	queue := make([]*nodeR, 1, t.maxDenseNodeCount<<1)
	queue[0] = t.root
	cur := make([]*nodeR, 1, t.maxDenseNodeCount)
	cur[0] = t.root
	height := 0
	var stopNodes []*nodeR
	for len(queue) > 0 && height < len(parts) {
		cur = cur[:0]
		part := parts[height]
		for len(queue) > 0 {
			head := queue[0]
			if head.stopChild != nil {
				stopNodes = append(stopNodes, head.stopChild)
			}
			if _, has := head.children[part]; has {
				cur = append(cur, head.children[part])
			}
			if head.jumpChild != nil {
				cur = append(cur, head.jumpChild)
			}
			queue = queue[1:]
		}
		height++
		if height >= len(parts) {
			break
		}
		queue = append(queue, cur...)
	}

	var nd *nodeR
	isStop := false
	if len(parts) == height {
		for _, n := range cur {
			if n.handle != nil {
				nd = n
				break
			}
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

	params := make(map[string]string, len(nd.keys))

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
