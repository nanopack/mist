// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package handlers

import (
	"encoding/json"
	"github.com/nanobox-io/nanobox-api"
	"github.com/nanopack/mist/core"
	"net/http"
)

type (
	Auth interface {
		TagsForToken(string) ([]string, error)
	}

	tagList struct {
		Tags []string `json:"tags"`
	}
)

func LoadWebsocketRoute(authenticator Auth) {
	upgrade := authenticate(authenticator)
	api.Router.Get("/subscribe/websocket", api.TraceRequest(upgrade))
}

func authenticate(authenticator Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if they have no tags registered for the token, then they are
		// not authorized to connect to mist
		token := r.Header.Get("x-auth-token")
		tags, err := authenticator.TagsForToken(token)
		if err != nil || len(tags) == 0 {
			w.WriteHeader(401)
			return
		}

		// we overwrite the subscribe command so that we can add
		// authentication to it.
		additionalCommands := map[string]mist.WebsocketHandler{
			"subscribe": buildWebsocketSubscribe(token, authenticator),
		}
		websocketUpgrade := mist.GenerateWebsocketUpgrade(api.User.(*mist.Mist), additionalCommands)
		websocketUpgrade(w, r)
	}
}

func buildWebsocketSubscribe(token string, authenticator Auth) mist.WebsocketHandler {
	return func(frame []byte, write chan<- string, client mist.Client) error {
		authTags, err := authenticator.TagsForToken(token)
		if err != nil || len(authTags) == 0 {
			write <- "{\"success\":false,\"command\":\"subscribe\"}"
			return nil
		}

		tags := tagList{}
		// error would already be caught by unmarshalling the command
		json.Unmarshal(frame, &tags)

		if !haveSameTags(authTags, tags.Tags) {
			write <- "{\"success\":false,\"command\":\"subscribe\"}"
			return nil
		}
		client.Subscribe(tags.Tags)
		write <- "{\"success\":true,\"command\":\"subscribe\"}"
		return nil
	}

}

func haveSameTags(a, b []string) bool {
	for _, vala := range a {
		for _, valb := range b {
			if vala == valb {
				return true
			}
		}
	}
	return false
}
