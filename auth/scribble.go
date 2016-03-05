package auth

import (
	"net/url"
)

//
func init() {
	authenticators["scribble"] = newScribble
}

//
func newScribble(url *url.URL) error {
	return nil
}
