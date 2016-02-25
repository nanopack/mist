//
package auth

import (
	"errors"
)

//
var (
	ErrTokenNotFound = errors.New("Token not found")
	ErrTokenExist    = errors.New("Token already exists")
)

//
type (

	//
	Authenticator interface {
		AddToken(token string) error
		RemoveToken(token string) error
		AddTags(token string, tags []string) error
		RemoveTags(token string, tags []string) error
		GetTagsForToken(token string) ([]string, error)
	}
)
