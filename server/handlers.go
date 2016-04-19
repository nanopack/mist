package server

import (
	"strings"

	"github.com/nanopack/mist/core"
)

// GenerateHandlers ...
func GenerateHandlers() map[string]mist.HandleFunc {
	return map[string]mist.HandleFunc{
		"auth":        handleAuth,
		"ping":        handlePing,
		"subscribe":   handleSubscribe,
		"unsubscribe": handleUnsubscribe,
		"publish":     handlePublish,
		// "publishAfter":     handlePublishAfter,
		"list": handleList,
	}
}

// handleAuth only exists to avoid getting the message "Unknown command" when
// authing with a authenticated server
func handleAuth(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Pipe <- mist.Message{Command: "auth", Tags: []string{}, Data: "success"}
	return nil
}

// handlePing
func handlePing(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Pipe <- mist.Message{Command: "ping", Tags: []string{}, Data: "pong"}
	return nil
}

// handleSubscribe
func handleSubscribe(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Subscribe(msg.Tags)
	proxy.Pipe <- mist.Message{Command: "subscribe", Tags: []string{}, Data: "pong"}
	return nil
}

// handleUnsubscribe
func handleUnsubscribe(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Unsubscribe(msg.Tags)
	proxy.Pipe <- mist.Message{Command: "unsubscribe", Tags: []string{}, Data: "pong"}
	return nil
}

// handlePublish
func handlePublish(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Publish(msg.Tags, msg.Data)
	proxy.Pipe <- mist.Message{Command: "publish", Tags: []string{}, Data: "pong"}
	return nil
}

// handlePublishAfter - how do we get the [delay] here?
// func handlePublishAfter(proxy *mist.Proxy, msg mist.Message) error {
// 	proxy.PublishAfter(msg.Tags, msg.Data, ???)
// 	proxy.Pipe <- mist.Message{Command: "publish after", Tags: msg.Tags, Data: "success"}
// 	return nil
// }

// handleList
func handleList(proxy *mist.Proxy, msg mist.Message) error {
	var subscriptions string
	for _, v := range proxy.List() {
		subscriptions += strings.Join(v, ",")
	}
	proxy.Pipe <- mist.Message{Command: "list", Tags: msg.Tags, Data: subscriptions}
	return nil
}
