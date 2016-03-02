package auth

import (
	"fmt"
	"strings"
)

type (

  // 
	Handler struct {
		NumArgs int
		Handle  func([]string) string
	}
)

//
func GenerateAuthCommands(auth Authenticator) map[string]Handler {
	return map[string]Handler{
		"register":   {2, handleRegister(auth)},
		"unregister": {1, handleUnregister(auth)},
		"set":        {2, handleSet(auth)},
		"unset":      {2, handleUnset(auth)},
		"tags":       {1, handleTags(auth)},
	}
}

// handleRegister
func handleRegister(a Authenticator) func([]string) string {
	return func(args []string) string {
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
func handleUnregister(a Authenticator) func([]string) string {
	return func(args []string) string {

		//
		if err := a.RemoveToken(args[0]); err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

		//
		return fmt.Sprintf("unregistered '%v'", args[0])
	}
}

// handleSet
func handleSet(a Authenticator) func([]string) string {
	return func(args []string) string {

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
func handleUnset(a Authenticator) func([]string) string {
	return func(args []string) string {

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
func handleTags(a Authenticator) func([]string) string {
	return func(args []string) string {

		tags, err := a.GetTagsForToken(args[0])
		if err != nil {
			return fmt.Sprintf("Error: %s", err.Error())
		}

		return fmt.Sprintf("tags for '%s' are [%s]", args[0], strings.Join(tags, ","))
	}
}
