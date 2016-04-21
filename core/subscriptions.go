package mist

const (
	create = iota
	remove
	nothing
)

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
type (

	// Node ...
	Node struct {
		id       uint64
		key      string
		count    int
		children map[string]*Node
		parent   *Node
		leafs    map[uint64]*Node
	}

	// Sub ...
	// Sub map[string]struct{}
)

// // newSub creates a new subscription
// func newSub() (sub Sub) {
// 	return make(Sub)
// }
//
// // Add adds tags to the subscription
// func (sub Sub) Add(tags []string) {
//
// 	// sort tags
// 	sort.Strings(tags)
//
// 	// join them into a string
// 	chk := strings.Join(tags, ",")
//
// 	// see if the key is in the map; if not insert it
// 	if _, found := sub[chk]; !found {
// 		sub[chk] = struct{}{}
// 	}
// }
//
// // Remove removes tags from the subscription
// func (sub Sub) Remove(tags []string) {
//
// 	// sort tags
// 	sort.Strings(tags)
//
// 	// delete tags
// 	delete(sub, strings.Join(tags, ","))
// }
//
// // Match checks to see if the subscription has tags
// func (sub Sub) Match(tags []string) (found bool) {
//
// 	// sort tags
// 	sort.Strings(tags)
//
// 	// see if the key is in the map
// 	_, found = sub[strings.Join(tags, ",")]
//
// 	//
// 	return
// }
//
// // ToSlice converts the subscription (map) to a slice of tags (slices)
// func (sub Sub) ToSlice() (slice [][]string) {
//
// 	// for each key in sub...
// 	for k := range sub {
//
// 		// split the key into a slice...
// 		split := strings.Split(k, ",")
//
// 		// ...and append that slice to the final
// 		slice = append(slice, split)
// 	}
//
// 	//
// 	return
// }

//
func newNode() (child *Node) {
	child = &Node{
		id:       0,
		count:    0,
		children: map[string]*Node{},
		leafs:    map[uint64]*Node{},
	}

	//
	return
}

// Len ...
func (root *Node) Len() int {
	return len(root.leafs)
}

// Add ...
func (root *Node) Add(keys []string) {
	if len(keys) == 0 {
		return
	}
	last, _ := root.traverse(keys, create)
	last.count++
	if last.count == 1 {
		last.id = root.id
		root.leafs[last.id] = last
		root.id++
	}
}

// Remove ...
func (root *Node) Remove(keys []string) {
	if len(keys) == 0 {
		return
	}
	found, _ := root.traverse(keys, remove)
	if found != nil {
		found.count--
		if found.count == 0 {
			delete(root.leafs, found.id)
		}
	}
}

// Match ...
func (root *Node) Match(keys []string) bool {
	last, count := root.traverse(keys, nothing)
	return last != nil && count != -1
}

// ToSlice ...
func (root *Node) ToSlice() [][]string {
	paths := make([][]string, len(root.leafs))
	for idx, leaf := range root.leafs {
		var path []string
		for ; leaf != nil && leaf != root; leaf = leaf.parent {
			path = append(path, leaf.key)
		}
		paths[idx] = path
	}
	return paths
}

// Find ...
func (root *Node) Find(keys []string) *Node {
	child, _ := root.traverse(keys, nothing)
	return child
}

// traverse ...
func (root *Node) traverse(keys []string, action int) (*Node, int) {
	if len(keys) == 0 {
		if root.count == 0 {
			// this node is not a leaf, so return -1 so it doesn't get deleted
			return root, -1
		}
		return root, root.count
	}

	key := keys[0]
	child, ok := root.children[key]

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
}
