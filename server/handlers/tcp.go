package handlers

import (
	"fmt"
	"strings"

	"github.com/nanopack/mist/core"
)

//
func GenerateTCPHandlers() map[string]mist.TCPHandler {
	return map[string]mist.TCPHandler{
		"ping":               {0, handleTCPPing},
		"list":               {0, handleTCPList},
		"subscribe":          {1, handleTCPSubscribe},
		"unsubscribe":        {1, handleTCPUnsubscribe},
		"publish":            {2, handleTCPPublish},
		"enable-replication": {0, handleTCPReplicate},
	}
}

//
func handleTCPPing(client mist.Client, args []string) string {
	fmt.Println("HANDLE TCP PING")
	return "pong"
}

//
func handleTCPSubscribe(client mist.Client, args []string) string {
	fmt.Println("HANDLE TCP SUB")
	tags := strings.Split(args[0], ",")
	client.Subscribe(tags)

	return fmt.Sprintf("subscribed %v", tags)
}

//
func handleTCPUnsubscribe(client mist.Client, args []string) string {
	fmt.Println("HANDLE TCP UNSUB")
	tags := strings.Split(args[0], ",")
	client.Unsubscribe(tags)

	return fmt.Sprintf("unsubscribed %v", tags)
}

//
func handleTCPPublish(client mist.Client, args []string) string {
	fmt.Println("HANDLE TCP PUB", args)
	tags := strings.Split(args[0], ",")
	msg := args[1]
	client.Publish(tags, msg)

	return fmt.Sprintf("published '%v' to %v", msg, tags)
}

//
func handleTCPList(client mist.Client, args []string) string {
	fmt.Println("HANDLE TCP LIST")
	subscriptions, err := client.List()
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("subscribed to %v", subscriptions)
}

//
func handleTCPReplicate(client mist.Client, args []string) string {
	client.(mist.Replicatable).EnableReplication()
	return "replication enabled"
}
