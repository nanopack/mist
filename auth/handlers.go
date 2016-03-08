package auth

import (
	"fmt"
	"strings"

	"github.com/nanopack/mist/core"
)

//
func GenerateHandlers() map[string]mist.Handler {
	return map[string]mist.Handler{
		"register":   {2, handleRegister},
		"unregister": {1, handleUnregister},
		"authorize":  {1, handleAuth},
		"set":        {2, handleSet},
		"unset":      {2, handleUnset},
		"tags":       {1, handleTags},
	}
}

// handleAuth
func handleAuth(proxy *mist.Proxy, args []string) error {

	//
	if _, err := DefaultAuth.GetTagsForToken(args[0]); err != nil {
		return fmt.Errorf("Incorrect token\n")
	}

	// authorize the proxy
	// proxy.Authorized = true

	return nil
}

// handleRegister
func handleRegister(proxy *mist.Proxy, args []string) error {

	token := args[0]

	//
	if err := DefaultAuth.AddToken(token); err != nil {
		return fmt.Errorf("%s\n", err.Error())
	}

	//
	if err := DefaultAuth.AddTags(token, strings.Split(args[1], ",")); err != nil {
		return fmt.Errorf("%s\n", err.Error())
	}

	//
	return nil
}

// handleUnregister
func handleUnregister(proxy *mist.Proxy, args []string) error {
	return DefaultAuth.RemoveToken(args[0])
}

// handleSet
func handleSet(proxy *mist.Proxy, args []string) error {
	return DefaultAuth.AddTags(args[0], strings.Split(args[1], ","))
}

// handleUnset
func handleUnset(proxy *mist.Proxy, args []string) error {
	return DefaultAuth.RemoveTags(args[0], strings.Split(args[1], ","))
}

// handleTags
func handleTags(proxy *mist.Proxy, args []string) error {

	// tags, err := DefaultAuth.GetTagsForToken(args[0])
	// if err != nil {
	// 	return fmt.Errorf("%s\n", err.Error())
	// }

	// proxy.Pipe <- tags

	//
	return nil
}
