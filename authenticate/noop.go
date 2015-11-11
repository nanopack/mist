// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package authenticate

import (
	"errors"
)

type (
	noop struct{}
)

var (
	Nothing = errors.New("I do nothing")
)

func NewNoopAuthenticator() noop {
	return noop{}
}

func (noop) TagsForToken(token string) ([]string, error) {
	return []string{}, Nothing
}
