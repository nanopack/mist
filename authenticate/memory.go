// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package authenticate

import (
	"errors"
	"github.com/deckarep/golang-set"
)

var (
	NotFound = errors.New("Token not found")
	Exist    = errors.New("already exists")
)

type (
	memory map[string]mapset.Set
)

func NewMemoryAuthenticator() memory {
	return memory{}
}

func (m memory) TagsForToken(token string) ([]string, error) {

	stored := m[token].ToSlice()
	tags := make([]string, len(stored))
	for idx, tag := range stored {
		tags[idx] = tag.(string)
	}
	return tags, nil
}

func (m memory) AddTags(token string, tags []string) error {
	current, ok := m[token]
	if !ok {
		return NotFound
	}
	for _, tag := range tags {
		current.Add(tag)
	}
	return nil
}

func (m memory) RemoveTags(token string, tags []string) error {
	current, ok := m[token]
	if !ok {
		return NotFound
	}
	for _, tag := range tags {
		current.Remove(tag)
	}
	return nil
}

func (m memory) AddToken(token string) error {
	_, ok := m[token]
	if ok {
		return Exist
	}
	m[token] = mapset.NewSet()
	return nil
}

func (m memory) RemoveToken(token string) error {
	delete(m, token)
	return nil
}
