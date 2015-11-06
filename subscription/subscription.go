// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package subscription

import (
	"sort"
)

const (
	create = iota
	remove
	nothing
)

type (
	node struct {
		id       uint64
		key      string
		count    int
		children map[string]*node
		parent   *node
		leafs    map[uint64]*node
	}
)

func NewNode() *node {
	child := newNode()
	child.leafs = map[uint64]*node{}
	return child
}

func newNode() *node {
	return &node{
		id:       0,
		count:    0,
		children: map[string]*node{},
	}
}

func (root *node) Len() int {
	return len(root.leafs)
}

func (root *node) Add(keys []string) {
	if len(keys) == 0 {
		return
	}
	sort.Sort(sort.StringSlice(keys))
	last, _ := root.traverse(keys, create)
	last.count++
	if last.count == 1 {
		last.id = root.id
		root.leafs[last.id] = last
		root.id++
	}
}

func (root *node) Remove(keys []string) {
	if len(keys) == 0 {
		return
	}
	sort.Sort(sort.StringSlice(keys))
	found, _ := root.traverse(keys, remove)
	if found != nil {
		found.count--
		if found.count == 0 {
			delete(root.leafs, found.id)
		}
	}
}

func (root *node) Match(keys []string) bool {
	sort.Sort(sort.StringSlice(keys))
	last, count := root.traverse(keys, nothing)
	return last != nil && count != -1
}

func (root *node) ToSlice(keys []string) [][]string {
	paths := make([][]string, len(root.leafs))
	for _, leaf := range root.leafs {
		path := make([]string, 1)
		for ; leaf != nil; leaf = leaf.parent {
			path = append(path, leaf.key)
		}
		paths = append(paths, path)
	}
	return paths
}

func (root *node) Find(keys []string) *node {
	sort.Sort(sort.StringSlice(keys))
	child, _ := root.traverse(keys, nothing)
	return child
}

func (root *node) traverse(keys []string, action int) (*node, int) {
	if len(keys) == 0 {
		if root.count == 0 {
			// this node is not a leaf, so return -1 so it doesn't get deleted
			return root, -1
		}
		return root, root.count
	}

	key := keys[0]

	var ok bool
	var child *node
	child, ok = root.children[key]

	switch action {
	case remove:
		if ok {
			found, count := child.traverse(keys[1:], action)
			if found != nil && count == 1 {
				if child.count == 0 && len(child.children) == 0 {
					delete(root.children, key)
				}
			}
			return found, count
		}
		return nil, 0
	case create:
		if !ok {
			child = newNode()
			child.parent = root
			child.key = keys[0] // preserve the original key
			root.children[key] = child
		}
		return child.traverse(keys[1:], action)
	default:
		if ok {
			found, count := child.traverse(keys[1:], action)
			// if 0 or -1 indicate that the traversal didn't really work
			if count != 0 {
				return found, count
			}
		}
		// we didn't find a match with the key, try the rest of the keys
		return root.traverse(keys[1:], action)
	}
	return nil, 0

}
