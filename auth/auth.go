// package auth provides a pluggable set of "Authenticators". These Authenticators
// represent a different database used to store authorized token/tag combinations
// that, when enabled, will allow/deny access to mist methods for these authorized
// token/tags combinations.
package auth

import (
	"fmt"
	"net/url"
)

//
var (
	DefaultAuth Authenticator // this is the current authenticator for the package; this is set during an authenticator start
	Token       string        // used by the server package when determining if auth command handlers should be added

	//
	ErrTokenNotFound = fmt.Errorf("Token not found\n")
	ErrTokenExist    = fmt.Errorf("Token already exists\n")

	// the list of available authenticators
	authenticators = map[string]func(url *url.URL) (Authenticator, error){}
)

//
type (

	// an authenticator represnets a database of authorized token/tag combinations.
	// These combinations are used as a way to allow access to mist methods for a
	// particular token/tag combination (when authentication is desired)
	Authenticator interface {
		AddToken(token string) error                    // add a token to list of authorized tokens
		RemoveToken(token string) error                 // remove a token from the list of authorized tokens
		AddTags(token string, tags []string) error      // add authorized tags to a token
		RemoveTags(token string, tags []string) error   // remove authorized tags from a token
		GetTagsForToken(token string) ([]string, error) // get the authorized tags for a token
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

	// set the package token
	Token = token

	// parse the uri string into a url object
	url, err := url.Parse(uri)
	if err != nil {
		return err
	}

	// check to see if the scheme is supported; if not, indicate as such and continue
	auth, ok := authenticators[url.Scheme]
	if !ok {
		return fmt.Errorf("Unsupported scheme '%v'", url.Scheme)
	}

	// set DefaultAuth by attempting to start the desired authenticator
	DefaultAuth, err = auth(url)
	if err != nil {
		return err
	}

	//
	return nil
}
