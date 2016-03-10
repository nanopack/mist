package mist

import (
	"sort"
	"testing"
)

type (

	//
	subTest struct {
		test   []string
		result bool
	}

	//
	matchTest struct {
		sub   []string
		tests []subTest
	}
)

// TestEmptySubscription
func TestEmptySubscription(t *testing.T) {
	node := newNode()
	if node.Len() != 0 {
		t.Fatalf("empty subscripition should be empty")
	}
}

// TestAddRemove
func TestAddRemove(t *testing.T) {
	node := newNode()

	//
	node.Add([]string{"a", "b"})
	if node.Len() != 1 {
		t.Fatalf("add should have added a node")
	}

	//
	node.Remove([]string{"a", "b"})
	if node.Len() != 0 {
		t.Fatalf("remove should have removed a node")
	}
}

// TestList
func TestList(t *testing.T) {
	node := newNode()

	//
	node.Add([]string{"a", "b"})
	list := node.ToSlice()
	if len(list) != 1 {
		t.Fatalf("wrong length of list. Expecting 1 got %v", list)
	}
	if len(list[0]) != 2 {
		t.Fatalf("wrong numer of tags. Expecing 2 got %v", len(list[0]))
	}
}

// TestAddRemoveDuplicate
func TestAddRemoveDuplicate(t *testing.T) {
	node := newNode()

	//
	node.Add([]string{"a", "b"})
	node.Add([]string{"a", "b"})
	if node.Len() != 1 {
		t.Fatalf("duplicate add should not have added a node")
	}

	//
	node.Remove([]string{"a", "b"})
	if node.Len() != 1 {
		t.Fatalf("duplicate remove should not have removed a node with a count > 1")
	}

	//
	node.Remove([]string{"a", "b"})
	if node.Len() != 0 {
		t.Fatalf("duplicate remove should have removed a node")
	}
}

// TestAddRemoveSimilarDuplicate
func TestAddRemoveSimilarDuplicate(t *testing.T) {
	node := newNode()

	//
	node.Add([]string{"a", "b"})
	if node.Len() != 1 {
		t.Fatalf("similar duplicate add should have added a node")
	}

	//
	node.Add([]string{"a", "b", "c"})
	if node.Len() != 2 {
		t.Fatalf("similar duplicate add should not have added a node")
	}

	//
	node.Remove([]string{"a", "b"})
	if node.Len() != 1 {
		t.Fatalf("similar duplicate remove should not have removed a node with a count > 1")
	}

	//
	node.Remove([]string{"a", "b", "c"})
	if node.Len() != 0 {
		t.Fatalf("similar duplicate remove should have removed a node")
	}
}

// TestMatches
func TestMatches(t *testing.T) {
	matches := []matchTest{
		{
			[]string{},
			[]subTest{
				{[]string{"a"}, false},
				{[]string{}, false},
			},
		},
		{
			[]string{"a", "b"},
			[]subTest{
				{[]string{}, false},
				{[]string{"a", "b"}, true},
				{[]string{"b", "a"}, true},
				{[]string{"b", "a", "c"}, true},
				{[]string{"c", "a", "b"}, true},
				{[]string{"c", "a"}, false},
				{[]string{"c", "b"}, false},
				{[]string{"c"}, false},
			},
		},
	}

	for _, match := range matches {
		node := newNode()
		sort.Sort(sort.StringSlice(match.sub))

		//
		node.Add(match.sub)

		//
		for _, mt := range match.tests {
			sort.Sort(sort.StringSlice(mt.test))
			if node.Match(mt.test) != mt.result {
				t.Fatalf("match failed: %v:%v expected %v", match.sub, mt.test, mt.result)
			}
		}

		//
		node.Remove(match.sub)

		//
		for _, mt := range match.tests {
			if node.Match(mt.test) {
				t.Fatalf("match failed: []:%v expected false", mt.test)
			}
		}

	}
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
