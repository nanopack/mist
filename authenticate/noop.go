package authenticate

import (
	"errors"
)

type (
	noop struct{}
)

var (
	Nothing = errors.New("I do nothing")
)

func NewNoopAuthenticator() noop {
	return noop{}
}

func (noop) TagsForToken(token string) ([]string, error) {
	return []string{}, Nothing
}

func (noop) AddTags(token string, tags []string) error {
	return Nothing
}

func (noop) RemoveTags(token string, tags []string) error {
	return Nothing
}

func (noop) AddToken(token string) error {
	return Nothing
}

func (noop) RemoveToken(token string) error {
	return Nothing
}
