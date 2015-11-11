// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package authenticate_test

import (
	"github.com/nanopack/mist/authenticate"
	"os/user"
	"testing"
)

type (
	postgresql    string
	Authenticator interface {
		TagsForToken(token string) ([]string, error)
		AddTags(token string, tags []string) error
		RemoveTags(token string, tags []string) error
		AddToken(token string) error
		RemoveToken(token string) error
	}
)

func TestPostgresql(test *testing.T) {
	usr, err := user.Current()
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	pg, err := authenticate.NewPostgresqlAuthenticator(usr.Username, "postgres", "127.0.0.1:5432")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	if err = pg.Clear(); err != nil {
		test.Log(err)
		test.FailNow()
	}
	testDb(test, pg)
}

func TestMemory(test *testing.T) {
	memory := authenticate.NewMemoryAuthenticator()
	testDb(test, memory)
}

func testDb(test *testing.T, auth Authenticator) {

	tags, err := auth.TagsForToken("token")
	if err == nil {
		test.Log("there should have been an error")
		test.FailNow()
	}
	if len(tags) != 0 {
		test.Log("wrong number of tags were returned")
		test.FailNow()
	}

	err = auth.AddToken("token")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	err = auth.AddTags("token", []string{"a", "b"})
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	tags, err = auth.TagsForToken("token")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
	if len(tags) != 2 {
		test.Log("wrong number of tags were returned", tags)
		test.FailNow()
	}

	err = auth.RemoveTags("token", []string{"a", "b"})
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	err = auth.RemoveToken("token")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
}
