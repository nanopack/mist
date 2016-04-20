package mist

import (
	"sort"
	"strings"
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

	// Sub ...
	Sub map[string]struct{}
)

// newSub creates a new subscription
func newSub() (sub Sub) {
	return make(Sub)
}

// Add adds tags to the subscription
func (sub Sub) Add(tags []string) {

	// sort tags
	sort.Strings(tags)

	// join them into a string
	chk := strings.Join(tags, ",")

	// see if the key is in the map; if not insert it
	if _, found := sub[chk]; !found {
		sub[chk] = struct{}{}
	}
}

// Remove removes tags from the subscription
func (sub Sub) Remove(tags []string) {

	// sort tags
	sort.Strings(tags)

	// delete tags
	delete(sub, strings.Join(tags, ","))
}

// Match checks to see if the subscription has tags
func (sub Sub) Match(tags []string) (found bool) {

	// sort tags
	sort.Strings(tags)

	// see if the key is in the map
	_, found = sub[strings.Join(tags, ",")]

	//
	return
}

// ToSlice converts the subscription (map) to a slice of tags (slices)
func (sub Sub) ToSlice() (slice [][]string) {

	// for each key in sub...
	for k := range sub {

		// split the key into a slice...
		split := strings.Split(k, ",")

		// ...and append that slice to the final
		slice = append(slice, split)
	}

	//
	return
}
