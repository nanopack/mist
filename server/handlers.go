package server

import (
	"github.com/nanopack/mist/core"
)

//
func GenerateHandlers() map[string]mist.HandleFunc {
	return map[string]mist.HandleFunc{
		"ping":        handlePing,
		"list":        handleList,
		"subscribe":   handleSubscribe,
		"unsubscribe": handleUnsubscribe,
		"publish":     handlePublish,
	}
}

// handlePing
func handlePing(proxy *mist.Proxy, msg mist.Message) error {
	return proxy.Ping()
}

// handleSubscribe
func handleSubscribe(proxy *mist.Proxy, msg mist.Message) error {
	return proxy.Subscribe(msg.Tags)
}

// handleUnsubscribe
func handleUnsubscribe(proxy *mist.Proxy, msg mist.Message) error {
	return proxy.Unsubscribe(msg.Tags)
}

// handlePublish
func handlePublish(proxy *mist.Proxy, msg mist.Message) error {
	return proxy.Publish(msg.Tags, msg.Data)
}

// handleList
func handleList(proxy *mist.Proxy, msg mist.Message) error {
	return proxy.List()
}
