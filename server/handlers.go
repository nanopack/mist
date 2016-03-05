package server

import (
	"strings"

	"github.com/nanopack/mist/core"
)

//
func GenerateHandlers() map[string]mist.Handler {
	return map[string]mist.Handler{
		"ping":        {0, handlePing},
		"list":        {0, handleList},
		"subscribe":   {1, handleSubscribe},
		"unsubscribe": {1, handleUnsubscribe},
		"publish":     {2, handlePublish},
	}
}

// handlePing
func handlePing(proxy *mist.Proxy, args []string) error {
	return nil
}

// handleSubscribe
func handleSubscribe(proxy *mist.Proxy, args []string) error {
	return proxy.Subscribe(strings.Split(args[0], ","))
}

// handleUnsubscribe
func handleUnsubscribe(proxy *mist.Proxy, args []string) error {
	return proxy.Unsubscribe(strings.Split(args[0], ","))
}

// handlePublish
func handlePublish(proxy *mist.Proxy, args []string) error {
	return proxy.Publish(strings.Split(args[0], ","), args[1])
}

// handleList
func handleList(proxy *mist.Proxy, args []string) error {
	return proxy.List()
}
