package auth

import (
	"fmt"
	"net/url"
	"strings"

	scribbleDB "github.com/nanobox-io/golang-scribble"
)

type (
	scribble struct {
		driver *scribbleDB.Driver
	}
)

//
func init() {
	authenticators["scribble"] = NewScribble
}

//
func NewScribble(url *url.URL) (Authenticator, error) {

	dir := url.Query().Get("db")

	// if no database location is provided, fail
	if dir == "" {
		return nil, fmt.Errorf("Missing database location in scheme (?db=)")
	}

	// create a new scribble at the specified location
	db, err := scribbleDB.New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	return &scribble{driver: db}, nil
}

//
func (a *scribble) AddToken(token string) error {
	return a.driver.Write("tokens", token, token)
}

//
func (a *scribble) RemoveToken(token string) error {
	return a.driver.Delete("tokens", token)
}

//
func (a *scribble) AddTags(token string, tags []string) error {
	return a.driver.Write(fmt.Sprintf("tokens/%s/tags", token), strings.Join(tags, "-"), tags)
}

//
func (a *scribble) RemoveTags(token string, tags []string) error {
	return a.driver.Delete(fmt.Sprintf("tokens/%s/tags", token), strings.Join(tags, "-"))
}

//
func (a *scribble) GetTagsForToken(token string) ([]string, error) {
	return a.driver.ReadAll(token)
}
