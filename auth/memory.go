package auth

import (
	"fmt"
	"net/url"

	"github.com/deckarep/golang-set"
)

type (
	memory map[string]mapset.Set
)

//
func init() {
	authenticators["memory"] = newMemory
}

//
func newMemory(url *url.URL) error {

	DefaultAuth = memory{}

	return nil
}

//
func (m memory) AddToken(token string) error {
	if _, ok := m[token]; ok {
		return ErrTokenExist
	}

	m[token] = mapset.NewSet()

	return nil
}

//
func (m memory) RemoveToken(token string) error {
	delete(m, token)
	return nil
}

//
func (m memory) AddTags(token string, tags []string) error {
	current, ok := m[token]
	if !ok {
		return ErrTokenNotFound
	}

	//
	for _, tag := range tags {
		current.Add(tag)
	}

	return nil
}

//
func (m memory) RemoveTags(token string, tags []string) error {
	current, ok := m[token]
	if !ok {
		return ErrTokenNotFound
	}

	//
	for _, tag := range tags {
		current.Remove(tag)
	}

	return nil
}

//
func (m memory) GetTagsForToken(token string) ([]string, error) {

	value, ok := m[token]
	if !ok {
		return nil, ErrTokenNotFound
	}

	fmt.Println("VALUE???", value)

	stored := value.ToSlice()

	fmt.Println("STORED???", stored)

	tags := make([]string, len(stored))

	fmt.Println("TAGS???", tags)

	//
	for idx, tag := range stored {
		tags[idx] = tag.(string)
	}

	//
	return tags, nil
}
