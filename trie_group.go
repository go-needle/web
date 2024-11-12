package web

import (
	"log"
)

type nodeG struct {
	handle   *RouterGroup
	children map[byte]*nodeG
}

func newNodeG(handle *RouterGroup) *nodeG {
	return &nodeG{handle: handle, children: make(map[byte]*nodeG)}
}

func (n *nodeG) matchChild(b byte) *nodeG {
	if _, has := n.children[b]; has {
		return n.children[b]
	}
	return nil
}

type trieTreeG struct {
	root                 *nodeG
	maxMiddleWaresLength int
}

func newTrieTreeG(rootHandle *RouterGroup) *trieTreeG {
	return &trieTreeG{newNodeG(rootHandle), 0}
}

func (t *trieTreeG) insert(prefix string, routerGroup *RouterGroup) int {
	cur := t.root
	middleWaresLength := 0
	for i := 0; i < len(prefix); i++ {
		next := cur.matchChild(prefix[i])
		if next == nil {
			next = newNodeG(nil)
			cur.children[prefix[i]] = next
		}
		if next.handle != nil {
			middleWaresLength += len(next.handle.middlewares)
		}
		cur = next
	}
	isAdd := true
	if cur.handle != nil {
		isAdd = false
		log.Printf("[Warning] A group coverage occurred in \"/%s\"", prefix)
	}
	cur.handle = routerGroup
	t.maxMiddleWaresLength = max(t.maxMiddleWaresLength, middleWaresLength+len(cur.handle.middlewares)>>1)
	if isAdd {
		return 1
	} else {
		return 0
	}
}

func (t *trieTreeG) search(prefix string) []HandlerFunc {
	cur := t.root
	middleWares := make([]HandlerFunc, 0, t.maxMiddleWaresLength)
	if cur.handle != nil {
		middleWares = append(middleWares, cur.handle.middlewares...)
	}
	for i := 0; i < len(prefix); i++ {
		next := cur.matchChild(prefix[i])
		if next == nil {
			break
		}
		if next.handle != nil {
			middleWares = append(middleWares, next.handle.middlewares...)
		}
		cur = next
	}
	return middleWares
}
