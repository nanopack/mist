// -*- mode: go; tab-width: 2; indent-tabs-mode: 1; st-rulers: [70] -*-
// vim: ts=4 sw=4 ft=lua noet
//--------------------------------------------------------------------
// @author Daniel Barney <daniel@nanobox.io>
// Copyright (C) Pagoda Box, Inc - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly
// prohibited. Proprietary and confidential
//
// @doc
//
// @end
// Created :   10 November 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
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
