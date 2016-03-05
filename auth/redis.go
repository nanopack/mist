package auth

import (
	"net/url"
)

//
func init() {
	authenticators["redis"] = newRedis
}

//
func newRedis(url *url.URL) error {
	return nil
}
