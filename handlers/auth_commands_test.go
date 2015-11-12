// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package handlers_test

import (
	"github.com/nanopack/mist/authenticate"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/handlers"
	"testing"
)

func get(cmds map[string]mist.Handler, cmd string) func(mist.Client, []string) string {
	return cmds[cmd].Handle
}

func TestAdditionalCommands(test *testing.T) {
	auth := authenticate.NewMemoryAuthenticator()
	cmds := handlers.GenerateAdditionalCommands(auth)

	reg := get(cmds, "register")
	unreg := get(cmds, "unregister")
	set := get(cmds, "set")
	unset := get(cmds, "unset")
	tags := get(cmds, "tags")

	// the client parameter is not used.
	if res := reg(nil, []string{"1,2,3,4", "token"}); res != "" {
		test.Log(res)
		test.FailNow()
	}

	if res := set(nil, []string{"a,b,c,d", "token"}); res != "" {
		test.Log(res)
		test.FailNow()
	}

	if tags := tags(nil, []string{"token"}); tags == "" {
		test.Log("wrong tags were returned")
		test.FailNow()
	}

	if res := unset(nil, []string{"a,b,c,d", "token"}); res != "" {
		test.Log(res)
		test.FailNow()
	}

	if res := unreg(nil, []string{"token"}); res != "" {
		test.Log(res)
		test.FailNow()
	}

	if tags := tags(nil, []string{"token"}); tags != "error Token not found" {
		test.Log("wrong tags were returned", tags)
		test.FailNow()
	}

}
