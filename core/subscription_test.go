package mist

import (
	"strings"
	"testing"
)

// // BenchmarkAddRemoveKeys benchmarks to see how fast mist can add/remove keys to
// // a subscription
// func BenchmarkAddRemoveKeys(b *testing.B) {
// 	node := newNode()
//
// 	// create a giant slice of random keys
// 	keys := [][]string{}
// 	for i := 0; i < b.N; i++ {
// 		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
// 	}
//
// 	b.ResetTimer()
//
// 	// add/remove keys
// 	for _, k := range keys {
// 		node.Add(k)
// 		node.Remove(k)
// 	}
// }
//
// // BenchmarkMatch benchmarks to see how fast mist can match a set of keys on a
// // subscription
// func BenchmarkMatch(b *testing.B) {
// 	node := newNode()
//
// 	// create a giant slice of random keys
// 	keys := [][]string{}
// 	for i := 0; i < b.N; i++ {
// 		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
// 	}
//
// 	b.ResetTimer()
//
// 	// add/match keys
// 	for _, k := range keys {
// 		node.Add(k)
// 		node.Match(k)
// 	}
// }

// TestEmptySubscription
func TestEmptySubscription(t *testing.T) {
	node := newNode()
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Unexpected tags in new subscription!")
	}
}

// TestAddRemoveSimple
func TestAddRemoveSimple(t *testing.T) {
	node := newNode()

	//
	node.Add([]string{"a"})
	if len(node.ToSlice()) != 1 {
		t.Fatalf("Failed to add node")
	}

	node.Remove([]string{"a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove node")
	}
}

// TestAddRemoveComplex
func TestAddRemoveComplex(t *testing.T) {
	node := newNode()

	// add/remove unordered keys; should remove
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"c", "b", "a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove node")
	}

	// add/remove incomplete keys; should not remove
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a"})
	node.Remove([]string{"b"})
	node.Remove([]string{"c"})
	node.Remove([]string{"a", "b"})
	node.Remove([]string{"b", "c"})
	node.Remove([]string{"a", "c"})
	node.Remove([]string{"b", "c", "d"})
	node.Remove([]string{"a", "b", "c", "d"})
	if len(node.ToSlice()) != 1 {
		t.Fatalf("Node unexpectedly removed")
	}

	// add duplicate keys; should only add once
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 1 {
		t.Fatalf("Duplicate nodes added")
	}
	node.Remove([]string{"a", "b", "c"})

	// remove duplicate keys; should only remove once
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a", "b", "c"})
	node.Remove([]string{"c", "b", "a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove nodes")
	}

	// add duplicate remote one; should leave no nodes
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove nodes")
	}
}

// TestList
func TestList(t *testing.T) {
	node := newNode()

	// test simple list; length should be 1 and value should be "a"
	node.Add([]string{"a"})
	list := node.ToSlice()
	if len(list) != 1 {
		t.Fatalf("Wrong number of keys - Expecting 1 got %v", len(list))
	}
	if len(list[0]) != 1 {
		t.Fatalf("Wrong number of keys - Expecing 2 got %v", len(list[0]))
	}
	if strings.Join(list[0], ",") != "a" {
		t.Fatalf("Wrong tags - Expecing 'a' got %v", list[0])
	}

	node.Add([]string{"a", "b"})
	list = node.ToSlice()
	if len(list) != 2 {
		t.Fatalf("Wrong number of keys - Expecting 2 got %v", len(list))
	}
	if len(list[1]) != 2 {
		t.Fatalf("Wrong number of keys - Expecing 2 got %v", len(list[1]))
	}
	if strings.Join(list[1], ",") != "a,b" {
		t.Fatalf("Wrong tags - Expecing 'a,b' got %v", list[1])
	}

	node.Add([]string{"a", "b", "c"})
	list = node.ToSlice()
	if len(list) != 3 {
		t.Fatalf("wrong length of list. Expecting 3 got %v", len(list))
	}
	if len(list[2]) != 3 {
		t.Fatalf("Wrong number of keys - Expecing 3 got %v", len(list[2]))
	}
	if strings.Join(list[2], ",") != "a,b,c" {
		t.Fatalf("Wrong tags - Expecing 'a,b,c' got %v", list[1])
	}
}

// TestMatchSimple
func TestMatchSimple(t *testing.T) {
	node := newNode()

	// simple match
	node.Add([]string{"a"})
	if !node.Match([]string{"a"}) {
		t.Fatalf("Expected match!")
	}

	//
	node.Add([]string{"a", "b"})
	if !node.Match([]string{"a", "b"}) {
		t.Fatalf("Expected match!")
	}

	//
	node.Add([]string{"a", "b", "c"})
	if !node.Match([]string{"a", "b", "c"}) {
		t.Fatalf("Expected match!")
	}
}

// TestMatchComplex
func TestMatchComplex(t *testing.T) {
	node := newNode()

	// match unordered keys; should match
	node.Add([]string{"a", "b", "c"})
	if !node.Match([]string{"c", "b", "a"}) {
		t.Fatalf("Expected match!")
	}

	// match incomplete keys; should match
	// node.Add([]string{"a", "b", "c"})
	// if !node.Match([]string{"a"}) {
	// 	t.Fatalf("Expected match!")
	// }
	// if !node.Match([]string{"a", "b"}) {
	// 	t.Fatalf("Expected match!")
	// }
	// if !node.Match([]string{"b", "c"}) {
	// 	t.Fatalf("Expected match!")
	// }
	// if !node.Match([]string{"a", "c"}) {
	// 	t.Fatalf("Expected match!")
	// }
}

// BenchmarkAddRemove
func BenchmarkAddRemove(b *testing.B) {
	node := newNode()
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Add(keys)
		node.Remove(keys)
	}
}

// BenchmarkMatch
func BenchmarkMatch(b *testing.B) {
	node := newNode()
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	node.Add(keys)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Match(keys)
	}
}

// // TestEmptySubscription
// func TestEmptySubscription(t *testing.T) {
// 	node := newNode()
// 	if node.Len() != 0 {
// 		t.Fatalf("empty subscripition should be empty")
// 	}
// }
//
// // TestAddRemove
// func TestAddRemove(t *testing.T) {
// 	node := newNode()
//
// 	//
// 	node.Add([]string{"a", "b"})
// 	if node.Len() != 1 {
// 		t.Fatalf("add should have added a node")
// 	}
//
// 	//
// 	node.Remove([]string{"a", "b"})
// 	if node.Len() != 0 {
// 		t.Fatalf("remove should have removed a node")
// 	}
//
// 	//
// 	node.Add([]string{"a", "b"})
// 	if node.Len() != 1 {
// 		t.Fatalf("add should have added a node")
// 	}
//
// 	//
// 	node.Remove([]string{"b", "a"})
// 	if node.Len() != 0 {
// 		t.Fatalf("remove should have removed a node")
// 	}
// }
//
// // TestList
// func TestList(t *testing.T) {
// 	node := newNode()
//
// 	//
// 	node.Add([]string{"a", "b"})
// 	list := node.ToSlice()
// 	if len(list) != 1 {
// 		t.Fatalf("wrong length of list. Expecting 1 got %v", list)
// 	}
// 	if len(list[0]) != 2 {
// 		t.Fatalf("wrong numer of tags. Expecing 2 got %v", len(list[0]))
// 	}
// }
//
// // TestAddRemoveDuplicate
// func TestAddRemoveDuplicate(t *testing.T) {
// 	node := newNode()
//
// 	//
// 	node.Add([]string{"a", "b"})
// 	node.Add([]string{"a", "b"})
// 	if node.Len() != 1 {
// 		t.Fatalf("duplicate add should not have added a node")
// 	}
//
// 	//
// 	node.Remove([]string{"a", "b"})
// 	if node.Len() != 1 {
// 		t.Fatalf("duplicate remove should not have removed a node with a count > 1")
// 	}
//
// 	//
// 	node.Remove([]string{"a", "b"})
// 	if node.Len() != 0 {
// 		t.Fatalf("duplicate remove should have removed a node")
// 	}
// }
//
// // TestAddRemoveSimilarDuplicate
// func TestAddRemoveSimilarDuplicate(t *testing.T) {
// 	node := newNode()
//
// 	//
// 	node.Add([]string{"a", "b"})
// 	if node.Len() != 1 {
// 		t.Fatalf("similar duplicate add should have added a node")
// 	}
//
// 	//
// 	node.Add([]string{"a", "b", "c"})
// 	if node.Len() != 2 {
// 		t.Fatalf("similar duplicate add should not have added a node")
// 	}
//
// 	//
// 	node.Remove([]string{"a", "b"})
// 	if node.Len() != 1 {
// 		t.Fatalf("similar duplicate remove should not have removed a node with a count > 1")
// 	}
//
// 	//
// 	node.Remove([]string{"a", "b", "c"})
// 	if node.Len() != 0 {
// 		t.Fatalf("similar duplicate remove should have removed a node")
// 	}
// }
//
// // TestMatches
// func TestMatches(t *testing.T) {
// 	matches := []matchTest{
// 		{
// 			[]string{},
// 			[]subTest{
// 				{[]string{"a"}, false},
// 				{[]string{}, false},
// 			},
// 		},
// 		{
// 			[]string{"a", "b"},
// 			[]subTest{
// 				{[]string{}, false},
// 				{[]string{"a", "b"}, true},
// 				{[]string{"b", "a"}, true},
// 				{[]string{"b", "a", "c"}, true},
// 				{[]string{"c", "a", "b"}, true},
// 				{[]string{"c", "a"}, false},
// 				{[]string{"c", "b"}, false},
// 				{[]string{"c"}, false},
// 			},
// 		},
// 	}
//
// 	for _, match := range matches {
// 		node := newNode()
// 		sort.Sort(sort.StringSlice(match.sub))
//
// 		//
// 		node.Add(match.sub)
//
// 		//
// 		for _, mt := range match.tests {
// 			sort.Sort(sort.StringSlice(mt.test))
// 			if node.Match(mt.test) != mt.result {
// 				t.Fatalf("match failed: %v:%v expected %v", match.sub, mt.test, mt.result)
// 			}
// 		}
//
// 		//
// 		node.Remove(match.sub)
//
// 		//
// 		for _, mt := range match.tests {
// 			if node.Match(mt.test) {
// 				t.Fatalf("match failed: []:%v expected false", mt.test)
// 			}
// 		}
//
// 	}
// }
//
// // BenchmarkAddRemove
// func BenchmarkAddRemove(b *testing.B) {
// 	node := newNode()
// 	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		node.Add(keys)
// 		node.Remove(keys)
// 	}
// }
//
// // BenchmarkMatch
// func BenchmarkMatch(b *testing.B) {
// 	node := newNode()
// 	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
// 	node.Add(keys)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		node.Match(keys)
// 	}
// }

// // BenchmarkAddRemoveKeys benchmarks to see how fast mist can add/remove keys to
// // a subscription
// func BenchmarkAddRemoveKeys(b *testing.B) {
// 	sub := newSub()
//
// 	// create a giant slice of random keys
// 	keys := [][]string{}
// 	for i := 0; i < b.N; i++ {
// 		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
// 	}
//
// 	b.ResetTimer()
//
// 	// add/remove keys
// 	for _, k := range keys {
// 		sub.Add(k)
// 		sub.Remove(k)
// 	}
// }
//
// // BenchmarkMatch benchmarks to see how fast mist can match a set of keys on a
// // subscription
// func BenchmarkMatch(b *testing.B) {
// 	sub := newSub()
//
// 	// create a giant slice of random keys
// 	keys := [][]string{}
// 	for i := 0; i < b.N; i++ {
// 		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
// 	}
//
// 	b.ResetTimer()
//
// 	// add/match keys
// 	for _, k := range keys {
// 		sub.Add(k)
// 		sub.Match(k)
// 	}
// }
//
// // TestEmptySubscription tests to ensure a new subscription is empty
// func TestEmptySubscription(t *testing.T) {
// 	sub := newSub()
// 	if len(sub) != 0 {
// 		t.Fatalf("Unexpected tags in new subscription!")
// 	}
// }
//
// // TestAddRemove tests to ensure that adding/removing tags works as expected
// func TestAddRemove(t *testing.T) {
// 	sub := newSub()
//
// 	// add a single tag and remove it
// 	sub.Add([]string{"a"})
// 	if len(sub) != 1 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
// 	}
// 	sub.Remove([]string{"a"})
// 	if len(sub) != 0 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
// 	}
//
// 	// add multiple tags and remove them
// 	sub.Add([]string{"a", "b"})
// 	if len(sub) != 1 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
// 	}
// 	sub.Remove([]string{"a", "b"})
// 	if len(sub) != 0 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
// 	}
//
// 	// add multiple tags and attempt to remove them unordered
// 	sub.Add([]string{"a", "b"})
// 	if len(sub) != 1 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
// 	}
// 	sub.Remove([]string{"b", "a"})
// 	if len(sub) != 0 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
// 	}
// }
//
// // TestList ensure listing tags on a subscription works as expected
// func TestList(t *testing.T) {
//
// 	sub := newSub()
//
// 	var list string
//
// 	// list single tag
// 	sub.Add([]string{"a"})
// 	list = flattenSliceToString(sub.ToSlice())
// 	if list != "a" {
// 		t.Fatalf("Unexpected tag - Expecting %v got %v", "a", list)
// 	}
// 	sub.Remove([]string{"a"})
//
// 	// list compound tags
// 	sub.Add([]string{"a", "b"})
// 	list = flattenSliceToString(sub.ToSlice())
// 	if list != "a,b" {
// 		t.Fatalf("Unexpected tags - Expecting %v got %v", "a,b", list)
// 	}
// 	sub.Remove([]string{"a", "b"})
//
// 	// list multiple tags; we test both configurations here because maps are unordered
// 	sub.Add([]string{"a", "b"})
// 	sub.Add([]string{"c"})
// 	list = flattenSliceToString(sub.ToSlice())
// 	switch list {
// 	case "a,bc", "ca,b":
// 	default:
// 		t.Fatalf("Unexpected tags - Expecting %v got %v", "'a,bc' OR 'ca,b'", list)
// 	}
// 	sub.Remove([]string{"a", "b"})
// 	sub.Remove([]string{"c"})
//
// 	// all tags should be removed
// 	tags := sub.ToSlice()
// 	if len(tags) != 0 {
// 		t.Fatalf("Unexpected length of tags - Expecting %v got %v", 0, len(tags))
// 	}
// }
//
// // TestAddRemoveDuplicate tests to ensure that adding duplicate tags don't actually
// // get added to the subscription
// func TestAddRemoveDuplicate(t *testing.T) {
// 	sub := newSub()
//
// 	// add dupicate tags and unordered duplicates
// 	sub.Add([]string{"a", "b"})
// 	sub.Add([]string{"a", "b"})
// 	sub.Add([]string{"b", "a"})
// 	if len(sub) != 1 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
// 	}
//
// 	// remove tags; multiple removes should have no effect
// 	sub.Remove([]string{"a", "b"})
// 	sub.Remove([]string{"a", "b"})
// 	sub.Remove([]string{"b", "a"})
// 	if len(sub) != 0 {
// 		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
// 	}
// }
//
// // TestMatch
// func TestMatch(t *testing.T) {
// 	sub := newSub()
//
// 	// add a single tag and test match
// 	sub.Add([]string{"a"})
// 	if !sub.Match([]string{"a"}) {
// 		t.Fatalf("Expected match!")
// 	}
//
// 	// add multiple tags and test for a match against unordered tags
// 	sub.Add([]string{"a", "b"})
// 	if !sub.Match([]string{"b", "a"}) {
// 		t.Fatalf("Expected match!")
// 	}
// }
