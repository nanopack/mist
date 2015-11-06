// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package subscription_test

import (
	"github.com/nanobox-io/golang-mist/subscription"
	"testing"
)

type (
	subTest struct {
		test   []string
		result bool
	}
	matchTest struct {
		sub   []string
		tests []subTest
	}
)

func TestEmptySubscription(test *testing.T) {
	node := subscription.NewNode()
	if node.Len() != 0 {
		test.Log("empty subscripition should be empty")
		test.FailNow()
	}
}

func TestAddRemove(test *testing.T) {
	node := subscription.NewNode()
	node.Add([]string{"a", "b"})
	if node.Len() != 1 {
		test.Log("add should have added a node")
		test.FailNow()
	}

	node.Remove([]string{"a", "b"})
	if node.Len() != 0 {
		test.Log("remove should have removed a node")
		test.FailNow()
	}
}

func TestAddRemoveDuplicate(test *testing.T) {
	node := subscription.NewNode()
	node.Add([]string{"a", "b"})
	if node.Len() != 1 {
		test.Log("duplicate add should have added a node")
		test.FailNow()
	}
	node.Add([]string{"a", "b"})
	if node.Len() != 1 {
		test.Log("duplicate add should not have added a node")
		test.FailNow()
	}

	node.Remove([]string{"a", "b"})
	if node.Len() != 1 {
		test.Log("duplicate remove should not have removed a node with a count > 1")
		test.FailNow()
	}
	node.Remove([]string{"a", "b"})
	if node.Len() != 0 {
		test.Log("duplicate remove should have removed a node")
		test.FailNow()
	}
}

func TestAddRemoveSimilarDuplicate(test *testing.T) {
	node := subscription.NewNode()
	node.Add([]string{"a", "b"})
	if node.Len() != 1 {
		test.Log("similar duplicate add should have added a node")
		test.FailNow()
	}
	node.Add([]string{"a", "b", "c"})
	if node.Len() != 2 {
		test.Log("similar duplicate add should not have added a node")
		test.FailNow()
	}

	node.Remove([]string{"a", "b"})
	if node.Len() != 1 {
		test.Log("similar duplicate remove should not have removed a node with a count > 1")
		test.FailNow()
	}
	node.Remove([]string{"a", "b", "c"})
	if node.Len() != 0 {
		test.Log("similar duplicate remove should have removed a node", node.Len())
		test.FailNow()
	}
}

func TestMatches(test *testing.T) {
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
		node := subscription.NewNode()
		node.Add(match.sub)

		for _, t := range match.tests {
			if node.Match(t.test) != t.result {
				test.Logf("match failed: %v:%v expected %v", match.sub, t.test, t.result)
				test.Fail()
			}
		}
		node.Remove(match.sub)
		for _, t := range match.tests {
			if node.Match(t.test) {
				test.Logf("match failed: []:%v expected false", t.test)
				test.Fail()
			}
		}

	}
}
