package saber

import (
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (t *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		t.pattern = pattern
		return
	}
	part := parts[height]
	child := func() *node {
		for _, next := range t.children {
			if next.isWild || next.part == part {
				return next
			}
		}
		return nil
	}()
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		t.children = append(t.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (t *node) search(parts []string, height int) *node {
	if len(parts) == 0 {
		return t
	}
	if len(parts) == height || strings.HasPrefix(t.part, "*") {
		if t.part == "" {
			return nil
		}
		return t
	}
	part := parts[height]
	children := func() []*node {
		nodes := make([]*node, 0)
		for _, next := range t.children {
			if next.part == part || next.isWild {
				nodes = append(nodes, next)
			}
		}
		return nodes
	}()
	for _, child := range children {
		ret := child.search(parts, height+1)
		if ret != nil {
			return ret
		}
	}
	return nil
}
