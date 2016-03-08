package auth

import (
	"fmt"
	"net/url"

	scribbleDB "github.com/nanobox-io/golang-scribble"
)

type (

	//
	scribble struct {
		driver *scribbleDB.Driver
	}

	//
	scribbleToken struct {
		Value string
		Tags  map[string]struct{}
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
		return nil, err
	}

	return &scribble{driver: db}, nil
}

//
func (a *scribble) AddToken(token string) error {

	// look for existing token; we want to fail if a token is found
	entry, err := a.findToken(token)
	if err == nil {
		return ErrTokenExist
	}

	//
	entry.Value = token

	// add new token
	return a.driver.Write("tokens", token, &entry)
}

//
func (a *scribble) RemoveToken(token string) error {
	return a.driver.Delete("tokens", token)
}

//
func (a *scribble) AddTags(token string, tags []string) error {

	// look for existing token
	entry, err := a.findToken(token)
	if err != nil {
		return err
	}

	// if this is the first time tags are being added to the token we need to
	// initialize them
	if entry.Tags == nil {
		entry.Tags = map[string]struct{}{}
	}

	// add new tags
	for _, tag := range tags {
		entry.Tags[tag] = struct{}{}
	}

	//
	return a.driver.Write("tokens", token, entry)
}

//
func (a *scribble) RemoveTags(token string, tags []string) error {

	// look for existing token
	entry, err := a.findToken(token)
	if err != nil {
		return err
	}

	// attempt to find tags and remove them
	for _, tag := range tags {
		delete(entry.Tags, tag)
	}

	// re-write entry w/o tags
	return a.driver.Write("tokens", token, entry)
}

//
func (a *scribble) GetTagsForToken(token string) ([]string, error) {

	// look for existing token
	entry, err := a.findToken(token)
	if err != nil {
		return nil, err
	}

	// convert tags from map to slice
	var tags []string
	for k, _ := range entry.Tags {
		tags = append(tags, k)
	}

	return tags, nil
}

//
func (a *scribble) findToken(token string) (entry scribbleToken, err error) {

	// look for existing token
	if err = a.driver.Read("tokens", token, &entry); err != nil {
		return entry, err
	}

	return entry, nil
}
