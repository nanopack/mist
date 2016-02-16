package handlers

import (
	"testing"
)

type (
	singleTest struct {
		a      []string
		b      []string
		result bool
	}
)

func TestCorrectAuth(test *testing.T) {
	tests := []singleTest{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"c", "a", "b"}, true},
		{[]string{"a", "b", "c"}, []string{"b", "c", "a"}, true},
		{[]string{"a", "b", "c"}, []string{"a", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"c"}, true},
		{[]string{"a", "b", "c"}, []string{}, false},
		{[]string{}, []string{"a", "b", "c"}, false},
	}
	for _, t := range tests {
		if haveSameTags(t.a, t.b) != t.result {
			test.Log("got the wrong result (%v) for %v:%v ", !t.result, t.a, t.b)
		}
	}
}
