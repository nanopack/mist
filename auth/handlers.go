package auth

import (
	"fmt"

	"github.com/nanopack/mist/core"
)

//
func GenerateHandlers() map[string]mist.HandleFunc {
	return map[string]mist.HandleFunc{
		"register":   handleRegister,
		"unregister": handleUnregister,
		"authorize":  handleAuth,
		"set":        handleSet,
		"unset":      handleUnset,
		"tags":       handleTags,
	}
}

// handleAuth
func handleAuth(proxy *mist.Proxy, msg mist.Message) error {

	//
	if _, err := DefaultAuth.GetTagsForToken(msg.Data); err != nil {
		return fmt.Errorf("Incorrect token\n")
	}

	// authorize the proxy
	// proxy.Authorized = true

	return nil
}

// handleRegister
func handleRegister(proxy *mist.Proxy, msg mist.Message) error {

	//
	if err := DefaultAuth.AddToken(msg.Data); err != nil {
		return fmt.Errorf("%s\n", err.Error())
	}

	//
	if err := DefaultAuth.AddTags(msg.Data, msg.Tags); err != nil {
		return fmt.Errorf("%s\n", err.Error())
	}

	//
	return nil
}

// handleUnregister
func handleUnregister(proxy *mist.Proxy, msg mist.Message) error {
	return DefaultAuth.RemoveToken(msg.Data)
}

// handleSet
func handleSet(proxy *mist.Proxy, msg mist.Message) error {
	return DefaultAuth.AddTags(msg.Data, msg.Tags)
}

// handleUnset
func handleUnset(proxy *mist.Proxy, msg mist.Message) error {
	return DefaultAuth.RemoveTags(msg.Data, msg.Tags)
}

// handleTags
func handleTags(proxy *mist.Proxy, msg mist.Message) error {

	// tags, err := DefaultAuth.GetTagsForToken(args[0])
	// if err != nil {
	// 	return fmt.Errorf("%s\n", err.Error())
	// }

	// proxy.Pipe <- tags

	//
	return nil
}
