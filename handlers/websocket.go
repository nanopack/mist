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
// Created :   12 August 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
package handlers

import (
	"bitbucket.org/nanobox/na-api"
	"encoding/json"
	"github.com/nanobox-io/golang-mist/core"
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
