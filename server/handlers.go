package server

import (
	"fmt"
	"strings"

	mist "github.com/nanopack/mist/core"
)

// GenerateHandlers ...
func GenerateHandlers() map[string]mist.HandleFunc {
	return map[string]mist.HandleFunc{
		"auth":        handleAuth,
		"ping":        handlePing,
		"subscribe":   handleSubscribe,
		"unsubscribe": handleUnsubscribe,
		"publish":     handlePublish,
		"list":        handleList,
		"listall":     handleListAll, // listall related
		"who":         handleWho,     // who related
	}
}

// handleAuth only exists to avoid getting the message "Unknown command" when
// authing with a authenticated server
func handleAuth(proxy *mist.Proxy, msg mist.Message) error {
	return nil
}

// handlePing
func handlePing(proxy *mist.Proxy, msg mist.Message) error {
	// goroutining any of these would allow a client to spam and overwhelm the server. clients don't need the ability to ping indefinitely
	proxy.Pipe <- mist.Message{Command: "ping", Tags: []string{}, Data: "pong"}
	return nil
}

// handleSubscribe
func handleSubscribe(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Subscribe(msg.Tags)
	return nil
}

// handleUnsubscribe
func handleUnsubscribe(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Unsubscribe(msg.Tags)
	return nil
}

// handlePublish
func handlePublish(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Publish(msg.Tags, msg.Data)
	return nil
}

// handleList
func handleList(proxy *mist.Proxy, msg mist.Message) error {
	var subscriptions string
	for _, v := range proxy.List() {
		subscriptions += strings.Join(v, ",")
	}
	proxy.Pipe <- mist.Message{Command: "list", Tags: msg.Tags, Data: subscriptions}
	return nil
}

// handleListAll - listall related
func handleListAll(proxy *mist.Proxy, msg mist.Message) error {
	subscriptions := mist.Subscribers()
	proxy.Pipe <- mist.Message{Command: "listall", Tags: msg.Tags, Data: subscriptions}
	return nil
}

// handleWho - who related
func handleWho(proxy *mist.Proxy, msg mist.Message) error {
	who, max := mist.Who()
	subscribers := fmt.Sprintf("Lifetime  connections: %d\nSubscribers connected: %d", max, who)
	proxy.Pipe <- mist.Message{Command: "who", Tags: msg.Tags, Data: subscribers}
	return nil
}
