package auth

import (
	"net/url"
)

type (
	redis struct{}
)

//
func init() {
	authenticators["redis"] = NewRedis
}

//
func NewRedis(url *url.URL) (Authenticator, error) {
	return &redis{}, nil
}

//
func (a *redis) AddToken(token string) error {
	return nil
}

//
func (a *redis) RemoveToken(token string) error {
	return nil
}

//
func (a *redis) AddTags(token string, tags []string) error {
	return nil
}

//
func (a *redis) RemoveTags(token string, tags []string) error {
	return nil
}

//
func (a *redis) GetTagsForToken(token string) ([]string, error) {
	return nil, nil
}
