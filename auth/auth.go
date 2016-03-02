//
package auth

import (
	"errors"
)

//
var (
	ErrTokenNotFound = errors.New("Token not found")
	ErrTokenExist    = errors.New("Token already exists")

	//
	authenticators = map[string]func(uri string, errChan chan<- error) {
		"memory": newMemory,
		"postgres": newPostgres,
		// "redis": newRedis,
		// "scribble": newScribble,
	}
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

// auth, ok := authenticators[viper.GetString("authenticator")]
