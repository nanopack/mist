package mist

import (
	"fmt"
	"strings"
)

var (

	// a map of the available commands that the server will respond to; format is
	// "name":{#args, handler}
	serverCommands = map[string]Handler{
		"ping":               {0, handlePing},
		"list":               {0, handleList},
		"subscribe":          {1, handleSubscribe},
		"unsubscribe":        {1, handleUnubscribe},
		"publish":            {2, handlePublish},
		"enable-replication": {0, handleEnableReplication},
	}
)

type (
	//
	Handler struct {
		ArgCount int
		Handle   func(Client, []string) string
	}
)

//
func handlePing(client Client, args []string) string {
	return "pong"
}

//
func handleSubscribe(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Subscribe(tags)
	return fmt.Sprintf("subscribed '%v'", tags)
}

//
func handleUnubscribe(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Unsubscribe(tags)
	return fmt.Sprintf("unsubscribed '%v'", tags)
}

//
func handlePublish(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	sub := args[1]
	client.Publish(tags, sub)
	return fmt.Sprintf("published '%v' to '%v'", tags, sub)
}

//
func handleList(client Client, args []string) string {
	list, err := client.List()
	if err != nil {
		return err.Error()
	}
	tmp := make([]string, len(list))

	for idx, subscription := range list {
		tmp[idx] = strings.Join(subscription, ",")
	}

	response := strings.Join(tmp, " ")
	return fmt.Sprintf("list %v", response)
}

//
func handleEnableReplication(client Client, args []string) string {
	client.(EnableReplication).EnableReplication()
	return ""
}
