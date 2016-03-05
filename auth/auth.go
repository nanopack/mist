//
package auth

import (
	"fmt"
	"net/url"
)

//
var (
	DefaultAuth Authenticator
	Token       string // used by the server package when determining if auth command handlers should be added

	//
	ErrTokenNotFound = fmt.Errorf("Token not found\n")
	ErrTokenExist    = fmt.Errorf("Token already exists\n")

	//
	authenticators = map[string]func(url *url.URL) error{}
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

// Start attempts to start a mist authenticator from the list of available
// authenticators; the authenticator provided is in the uri string format
// (scheme:[//[user:pass@]host[:port]][/]path[?query][#fragment])
func Start(uri, token string) error {

	// no authenticator is wanted
	if uri == "" {
		return nil
	}

	// check to see if a token is provided; an authenticator cannot work without
	// a token and so it should error here informing that.
	if token == "" {
		return fmt.Errorf("An authenticator has been specified but no token provided!\n")
	}

	// parse the uri string into a url object
	url, err := url.Parse(uri)
	if err != nil {
		return err
	}

	// check to see if the scheme is supported; if not, indicate as such and
	// continue
	auth, ok := authenticators[url.Scheme]
	if !ok {
		return fmt.Errorf("Unsupported scheme '%v'", url.Scheme)
	}

	//
	Token = token

	// attempt to start the authenticator
	return auth(url)
}
