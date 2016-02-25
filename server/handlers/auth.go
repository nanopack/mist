package handlers

import (
	"fmt"
	"strings"

	"github.com/nanopack/mist/core"
  "github.com/nanopack/mist/auth"
)

//
func GenerateAuthCommands(auth auth.Authenticator) map[string]mist.TCPHandler {
	return map[string]mist.TCPHandler{
		"register":   {2, handleRegister(auth)},
		"unregister": {1, handleUnregister(auth)},
		"set":        {2, handleSet(auth)},
		"unset":      {2, handleUnset(auth)},
		"tags":       {1, handleTags(auth)},
	}
}

// handleRegister
func handleRegister(a auth.Authenticator) func(mist.Client, []string) string {
	return func(client mist.Client, args []string) string {
		var err error

		token := args[0]
		tags := strings.Split(args[1], ",")

		//
		if err = a.AddToken(token); err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

		//
		if err = a.AddTags(token, tags); err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

    //
		return fmt.Sprintf("registered [%v] to '%v'", tags, token)
	}
}

// handleUnregister
func handleUnregister(a auth.Authenticator) func(mist.Client, []string) string {
	return func(client mist.Client, args []string) string {

		//
		if err := a.RemoveToken(args[0]); err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

    //
		return fmt.Sprintf("unregistered '%v'", args[0])
	}
}

// handleSet
func handleSet(a auth.Authenticator) func(mist.Client, []string) string {
	return func(client mist.Client, args []string) string {

		token := args[0]
		tags := strings.Split(args[1], ",")

		//
		if err := a.AddTags(token, tags); err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

		return fmt.Sprintf("added [%v] to '%v'", tags, token)
	}
}

// handleUnset
func handleUnset(a auth.Authenticator) func(mist.Client, []string) string {
	return func(client mist.Client, args []string) string {

		token := args[0]
		tags := strings.Split(args[1], ",")

		//
		if err := a.RemoveTags(token, tags); err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

		return fmt.Sprintf("removed [%v] from '%v'", tags, token)
	}
}

// handleTags
func handleTags(a auth.Authenticator) func(mist.Client, []string) string {
	return func(client mist.Client, args []string) string {

		tags, err := a.GetTagsForToken(args[0])
		if err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

		return fmt.Sprintf("tags for '%s' are [%s]", args[0], strings.Join(tags, ","))
	}
}
