// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package handlers

import (
	"github.com/nanobox-io/golang-mist/core"
	"strings"
)

type (
	Authenticator interface {
		TagsForToken(token string) ([]string, error)
		AddTags(token string, tags []string) error
		RemoveTags(token string, tags []string) error
		AddToken(token string) error
		RemoveToken(token string) error
	}

	handleFun func(mist.Client, []string) string
)

func GenerateAdditionalCommands(auth Authenticator) map[string]mist.Handler {
	return map[string]mist.Handler{
		"register":   {1, handleRegister(auth)},
		"unregister": {1, handleUnregister(auth)},
		"set":        {2, handleSet(auth)},
		"unset":      {2, handleUnset(auth)},
		"tags":       {1, handleGetTags(auth)},
	}
}

func handleRegister(auth Authenticator) handleFun {
	return func(client mist.Client, args []string) string {
		err := auth.AddToken(args[0])
		if err != nil {
			return "error " + err.Error()
		}
		return ""
	}
}

func handleUnregister(auth Authenticator) handleFun {
	return func(client mist.Client, args []string) string {
		err := auth.RemoveToken(args[0])
		if err != nil {
			return "error " + err.Error()
		}
		return ""
	}
}

func handleSet(auth Authenticator) handleFun {
	return func(client mist.Client, args []string) string {
		tags := strings.Split(args[1], ",")
		err := auth.AddTags(args[0], tags)
		if err != nil {
			return "error " + err.Error()
		}
		return ""
	}
}

func handleUnset(auth Authenticator) handleFun {
	return func(client mist.Client, args []string) string {
		tags := strings.Split(args[1], ",")
		err := auth.RemoveTags(args[0], tags)
		if err != nil {
			return "error " + err.Error()
		}
		return ""
	}
}

func handleGetTags(auth Authenticator) handleFun {
	return func(client mist.Client, args []string) string {
		tags, err := auth.TagsForToken(args[0])
		if err != nil {
			return "error " + err.Error()
		}
		return "tags " + args[0] + " " + strings.Join(tags, ",")
	}
}
