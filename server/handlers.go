package server

import (
	"github.com/nanopack/mist/core"
)

// GenerateHandlers ...
func GenerateHandlers() map[string]mist.HandleFunc {
	return map[string]mist.HandleFunc{
		"auth":        handleAuth,
		"ping":        handlePing,
		"list":        handleList,
		"subscribe":   handleSubscribe,
		"unsubscribe": handleUnsubscribe,
		"publish":     handlePublish,
		// "publishAfter":     handlePublishAfter,
	}
}

// handleAuth
func handleAuth(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Pipe <- mist.Message{Command: "auth", Tags: []string{}, Data: "success"}
	return nil
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

// // handlePublishAfter - do we want a publish after here?
// func handlePublishAfter(proxy *mist.Proxy, msg mist.Message) error {
// 	return proxy.PublishAfter(msg.Tags, msg.Data)
// }

// handleList
func handleList(proxy *mist.Proxy, msg mist.Message) error {
	return proxy.List()
}
