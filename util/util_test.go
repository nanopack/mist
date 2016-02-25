package util

import (
  "testing"
)

//
func TestNext(t *testing.T) {
  
}

//
func TestHaveSameTags(t *testing.T) {
	tags := []struct{
		a      []string
		b      []string
		result bool
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"c", "b", "a"}, true},
		{[]string{"a", "b", "c"}, []string{"a", "b"}, true},
		{[]string{"a", "b", "c"}, []string{"c"}, true},
		{[]string{"a", "b", "c"}, []string{}, false},
		{[]string{}, []string{"a", "b", "c"}, false},
	}

	//
	for _, tag := range tags {
		if HaveSameTags(tag.a, tag.b) != tag.result {
			t.Log("got the wrong result (%v) for %v:%v ", !tag.result, tag.a, tag.b)
		}
	}
}
