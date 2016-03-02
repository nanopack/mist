package auth

import (
	"github.com/deckarep/golang-set"
)

type (
	memory map[string]mapset.Set
)

//
func newMemory(uri string, errChan chan<- error) {
}

//
func NewMemory() memory {
	return memory{}
}

//
func (m memory) AddToken(token string) error {
	_, ok := m[token]
	if ok {
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
	for _, tag := range tags {
		current.Remove(tag)
	}
	return nil
}

//
func (m memory) GetTagsForToken(token string) ([]string, error) {

	value, ok := m[token]
	if !ok {
		return []string{}, ErrTokenNotFound
	}
	stored := value.ToSlice()
	tags := make([]string, len(stored))
	for idx, tag := range stored {
		tags[idx] = tag.(string)
	}
	return tags, nil
}
