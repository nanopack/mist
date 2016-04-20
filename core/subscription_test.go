package mist

import "testing"

// BenchmarkAddRemoveKeys benchmarks to see how fast mist can add/remove keys to
// a subscription
func BenchmarkAddRemoveKeys(b *testing.B) {
	sub := newSub()

	// create a giant slice of random keys
	keys := [][]string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
	}

	b.ResetTimer()

	// add/remove keys
	for _, k := range keys {
		sub.Add(k)
		sub.Remove(k)
	}
}

// BenchmarkMatch benchmarks to see how fast mist can match a set of keys on a
// subscription
func BenchmarkMatch(b *testing.B) {
	sub := newSub()

	// create a giant slice of random keys
	keys := [][]string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
	}

	b.ResetTimer()

	// add/match keys
	for _, k := range keys {
		sub.Add(k)
		sub.Match(k)
	}
}

// TestEmptySubscription tests to ensure a new subscription is empty
func TestEmptySubscription(t *testing.T) {
	sub := newSub()
	if len(sub) != 0 {
		t.Fatalf("Unexpected tags in new subscription!")
	}
}

// TestAddRemove tests to ensure that adding/removing tags works as expected
func TestAddRemove(t *testing.T) {
	sub := newSub()

	// add a single tag and remove it
	sub.Add([]string{"a"})
	if len(sub) != 1 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
	}
	sub.Remove([]string{"a"})
	if len(sub) != 0 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
	}

	// add multiple tags and remove them
	sub.Add([]string{"a", "b"})
	if len(sub) != 1 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
	}
	sub.Remove([]string{"a", "b"})
	if len(sub) != 0 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
	}

	// add multiple tags and attempt to remove them unordered
	sub.Add([]string{"a", "b"})
	if len(sub) != 1 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
	}
	sub.Remove([]string{"b", "a"})
	if len(sub) != 0 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
	}
}

// TestList ensure listing tags on a subscription works as expected
func TestList(t *testing.T) {

	sub := newSub()

	var list string

	// list single tag
	sub.Add([]string{"a"})
	list = flattenSliceToString(sub.ToSlice())
	if list != "a" {
		t.Fatalf("Unexpected tag - Expecting %v got %v", "a", list)
	}
	sub.Remove([]string{"a"})

	// list compound tags
	sub.Add([]string{"a", "b"})
	list = flattenSliceToString(sub.ToSlice())
	if list != "a,b" {
		t.Fatalf("Unexpected tags - Expecting %v got %v", "a,b", list)
	}
	sub.Remove([]string{"a", "b"})

	// list multiple tags; we test both configurations here because maps are unordered
	sub.Add([]string{"a", "b"})
	sub.Add([]string{"c"})
	list = flattenSliceToString(sub.ToSlice())
	switch list {
	case "a,bc", "ca,b":
	default:
		t.Fatalf("Unexpected tags - Expecting %v got %v", "'a,bc' OR 'ca,b'", list)
	}
	sub.Remove([]string{"a", "b"})
	sub.Remove([]string{"c"})

	// all tags should be removed
	tags := sub.ToSlice()
	if len(tags) != 0 {
		t.Fatalf("Unexpected length of tags - Expecting %v got %v", 0, len(tags))
	}
}

// TestAddRemoveDuplicate tests to ensure that adding duplicate tags don't actually
// get added to the subscription
func TestAddRemoveDuplicate(t *testing.T) {
	sub := newSub()

	// add dupicate tags and unordered duplicates
	sub.Add([]string{"a", "b"})
	sub.Add([]string{"a", "b"})
	sub.Add([]string{"b", "a"})
	if len(sub) != 1 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 1, len(sub))
	}

	// remove tags; multiple removes should have no effect
	sub.Remove([]string{"a", "b"})
	sub.Remove([]string{"a", "b"})
	sub.Remove([]string{"b", "a"})
	if len(sub) != 0 {
		t.Fatalf("Unexpected number of tags - Expecting %v got %v", 0, len(sub))
	}
}

// TestMatch
func TestMatch(t *testing.T) {
	sub := newSub()

	// add a single tag and test match
	sub.Add([]string{"a"})
	if !sub.Match([]string{"a"}) {
		t.Fatalf("Expected match!")
	}

	// add multiple tags and test for a match against unordered tags
	sub.Add([]string{"a", "b"})
	if !sub.Match([]string{"b", "a"}) {
		t.Fatalf("Expected match!")
	}
}
