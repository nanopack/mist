package auth

import (
	"fmt"
)

type (
	websocket struct{}
)

//
func NewWS() websocket {
	return websocket{}
}

//
func (ws websocket) AddToken(token string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) RemoveToken(token string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) AddTags(token string, tags []string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) RemoveTags(token string, tags []string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) GetTagsForToken(token string) ([]string, error) {
	return []string{}, fmt.Errorf("I dont do anything...")
}
