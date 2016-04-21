package mist

import (
	"fmt"
	"sort"
)

// const (
// 	create = iota
// 	remove
// 	nothing
// )

// interfaces
type (

	//
	subscriptions interface {
		Add([]string)
		Remove([]string)
		Match([]string) bool
		ToSlice() [][]string
	}
)

// structs
// type (
//
// 	// // Node ...
// 	// Node struct {
// 	// 	id       uint64
// 	// 	key      string
// 	// 	count    int
// 	// 	children map[string]*Node
// 	// 	parent   *Node
// 	// 	leafs    map[uint64]*Node
// 	// }
//
// )

// //
// func newNode() (child *Node) {
// 	child = &Node{
// 		id:       0,
// 		count:    0,
// 		children: map[string]*Node{},
// 		leafs:    map[uint64]*Node{},
// 	}
//
// 	//
// 	return
// }
//
// // Len ...
// func (root *Node) Len() int {
// 	return len(root.leafs)
// }
//
// // Add ...
// func (root *Node) Add(keys []string) {
// 	if len(keys) == 0 {
// 		return
// 	}
//
// 	// sort tags
// 	sort.Strings(keys)
//
// 	last, _ := root.traverse(keys, create)
// 	last.count++
// 	if last.count == 1 {
// 		last.id = root.id
// 		root.leafs[last.id] = last
// 		root.id++
// 	}
// }
//
// // Remove ...
// func (root *Node) Remove(keys []string) {
// 	if len(keys) == 0 {
// 		return
// 	}
//
// 	// sort tags
// 	sort.Strings(keys)
//
// 	found, _ := root.traverse(keys, remove)
// 	if found != nil {
// 		found.count--
// 		if found.count == 0 {
// 			delete(root.leafs, found.id)
// 		}
// 	}
// }
//
// // Match ...
// func (root *Node) Match(keys []string) bool {
//
// 	if len(keys) == 0 {
// 		return false
// 	}
//
// 	// sort tags
// 	sort.Strings(keys)
//
// 	last, count := root.traverse(keys, nothing)
// 	return last != nil && count != -1
// }
//
// // ToSlice ...
// func (root *Node) ToSlice() [][]string {
// 	paths := make([][]string, len(root.leafs))
// 	for idx, leaf := range root.leafs {
// 		var path []string
// 		for ; leaf != nil && leaf != root; leaf = leaf.parent {
// 			path = append(path, leaf.key)
// 		}
// 		paths[idx] = path
// 	}
// 	return paths
// }
//
// // Find ...
// func (root *Node) Find(keys []string) *Node {
//
// 	if len(keys) == 0 {
// 		return nil
// 	}
//
// 	// sort tags
// 	sort.Strings(keys)
//
// 	child, _ := root.traverse(keys, nothing)
// 	return child
// }
//
// // traverse ...
// func (root *Node) traverse(keys []string, action int) (*Node, int) {
// 	if len(keys) == 0 {
// 		if root.count == 0 {
// 			// this node is not a leaf, so return -1 so it doesn't get deleted
// 			return root, -1
// 		}
// 		return root, root.count
// 	}
//
// 	key := keys[0]
// 	child, ok := root.children[key]
//
// 	switch action {
// 	case remove:
// 		if ok {
// 			found, count := child.traverse(keys[1:], action)
// 			if found != nil && count == 1 {
// 				if child.count == 0 && len(child.children) == 0 {
// 					delete(root.children, key)
// 				}
// 			}
// 			return found, count
// 		}
//
// 		return nil, 0
// 	case create:
// 		if !ok {
// 			child = newNode()
// 			child.parent = root
// 			child.key = keys[0] // preserve the original key
// 			root.children[key] = child
// 		}
//
// 		return child.traverse(keys[1:], action)
// 	default:
// 		if ok {
// 			found, count := child.traverse(keys[1:], action)
// 			// if 0 or -1 indicate that the traversal didn't really work
// 			if count != 0 {
// 				return found, count
// 			}
// 		}
//
// 		// we didn't find a match with the key, try the rest of the keys
// 		return root.traverse(keys[1:], action)
// 	}
// }

type (

	// Node ...
	Node struct {
		// toSlice  [][]string
		branches map[string]*Node
		leaves   map[string]struct{}
	}
)

//
func newNode() (node *Node) {

	node = &Node{
		// id:       0,
		// count:    0,
		branches: map[string]*Node{},
		leaves:   map[string]struct{}{},
	}

	//
	return
}

// Show ...
func (node *Node) Show(indent string) {
	fmt.Printf("%s%#v\n", indent, node)
	for _, branch := range node.branches {
		branch.Show(fmt.Sprintf("  %s", indent))
	}
}

// Add sorts the keys and then attempts to add them
func (node *Node) Add(keys []string) {

	//
	if len(keys) == 0 {
		return
	}

	sort.Strings(keys)
	node.add(keys)
}

// add ...
func (node *Node) add(keys []string) {

	// if there is only one key remaining we are at the end of the chain, so a leaf
	// is created
	if len(keys) == 1 {
		node.leaves[keys[0]] = struct{}{}
		return
	}

	// see if there is already a branch for the first key. if not create a new one
	// and add it; if a branch already exists we simply use it and continue
	branch, ok := node.branches[keys[0]]
	if !ok {
		branch = newNode()
		node.branches[keys[0]] = branch
	}

	// for the current branch (new or existing) continue to the next set of keys
	// adding them as branches until there is only one left, at which point a leaf
	// is created (above)
	branch.add(keys[1:])
}

// Remove sorts the keys and then attempts to remove them
func (node *Node) Remove(keys []string) {

	//
	if len(keys) == 0 {
		return
	}

	sort.Strings(keys)
	node.remove(keys)
}

// remove ...
func (node *Node) remove(keys []string) {

	// if there is only one key remaining we are at the end of the chain and need
	// to remove just the leaf
	if len(keys) == 1 {
		delete(node.leaves, keys[0])
		return
	}

	// see if a branch for the first key exists; if a branch exists we need to
	// recurse down the branch until we reach the end...
	branch, ok := node.branches[keys[0]]
	if ok {

		// continue key by key until we reach a leaf at which point its removed (above)
		branch.remove(keys[1:])

		// once we reach the end of the line, if there are no more leaves or branch
		// on this branch, we can remove the branch
		if len(branch.leaves) == 0 && len(branch.branches) == 0 {
			delete(node.branches, keys[0])
		}
	}
}

// Match sorts the keys and then attempts to find a match
func (node *Node) Match(keys []string) bool {

	//
	if len(keys) == 0 {
		return false
	}

	sort.Strings(keys)
	return node.match(keys)
}

// ​match ...
func (node *Node) match(keys []string) bool {

	// iterate through each key looking for a leaf, if found it's a match
	for _, key := range keys {
		if _, ok := node.leaves[key]; ok {
			return true
		}
	}

	// see if a branch for the first key exists; if no branch exists we need to
	// try and find a match for the next key, continuing down the chain until we
	// find a leaf (above)
	branch, ok := node.branches[keys[0]]
	if !ok {
		return node.match(keys[1:])
	}

	// if a branch does exist we down the branch until we find a leaf (above)
	return branch.match(keys[1:])
}

// ToSlice recurses down an entire node returning a list of all branches and leaves
// as a slice of slices
func (node *Node) ToSlice() (list [][]string) {

	fmt.Printf("TO SLICE! %#v\n", node)

	// iterate through each leaf appending it as a slice to the list of keys
	for leaf := range node.leaves {
		fmt.Println("LEAF?", leaf)
		list = append(list, []string{leaf})
	}

	// iterate through each branch getting its list of branches and appending those
	// to the list
	for branch, node := range node.branches {

		fmt.Println("BRANCH?", branch)

		// get the current nodes slice of branches and leaves
		nodeSlice := node.ToSlice()

		fmt.Println("SLICE?", nodeSlice)

		// for each branch in the nodes list apppend the key to that key
		for _, nodeKey := range nodeSlice {
			fmt.Println("WHERE AM I?", nodeKey)
			list = append(list, append(nodeKey, branch))
		}
	}

	// sort each list
	for _, l := range list {
		sort.Strings(l)
	}

	return
}
