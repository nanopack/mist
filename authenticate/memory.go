package authenticate

import (
	"errors"

	"github.com/deckarep/golang-set"
)

var (
	NotFound = errors.New("Token not found")
	Exist    = errors.New("already exists")
)

type (
	memory map[string]mapset.Set
)

func NewMemoryAuthenticator() memory {
	return memory{}
}

func (m memory) TagsForToken(token string) ([]string, error) {

	value, ok := m[token]
	if !ok {
		return []string{}, NotFound
	}
	stored := value.ToSlice()
	tags := make([]string, len(stored))
	for idx, tag := range stored {
		tags[idx] = tag.(string)
	}
	return tags, nil
}

func (m memory) AddTags(token string, tags []string) error {
	current, ok := m[token]
	if !ok {
		return NotFound
	}
	for _, tag := range tags {
		current.Add(tag)
	}
	return nil
}

func (m memory) RemoveTags(token string, tags []string) error {
	current, ok := m[token]
	if !ok {
		return NotFound
	}
	for _, tag := range tags {
		current.Remove(tag)
	}
	return nil
}

func (m memory) AddToken(token string) error {
	_, ok := m[token]
	if ok {
		return Exist
	}
	m[token] = mapset.NewSet()
	return nil
}

func (m memory) RemoveToken(token string) error {
	delete(m, token)
	return nil
}
