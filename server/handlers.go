package server

import (
  "fmt"
  "strings"

  "github.com/nanopack/mist/core"
)

type (
	Handler struct {
		NumArgs int
		Handle  func(mist.Proxy, []string) string
	}
)

//
func GenerateHandlers() map[string]Handler {
	return map[string]Handler{
		"ping":               {0, handlePing},
		"list":               {0, handleList},
		"subscribe":          {1, handleSubscribe},
		"unsubscribe":        {1, handleUnsubscribe},
		"publish":            {2, handlePublish},
	}
}

//
func handlePing(client mist.Proxy, args []string) string {
	fmt.Println("HANDLE  PING")
	return "pong"
}

//
func handleSubscribe(client mist.Proxy, args []string) string {
	fmt.Println("HANDLE  SUB")
	tags := strings.Split(args[0], ",")
	client.Subscribe(tags)

	return fmt.Sprintf("subscribed %v", tags)
}

//
func handleUnsubscribe(client mist.Proxy, args []string) string {
	fmt.Println("HANDLE  UNSUB")
	tags := strings.Split(args[0], ",")
	client.Unsubscribe(tags)

	return fmt.Sprintf("unsubscribed %v", tags)
}

//
func handlePublish(client mist.Proxy, args []string) string {
	fmt.Println("HANDLE  PUB", args)
	tags := strings.Split(args[0], ",")
	msg := args[1]
	client.Publish(tags, msg)

	return fmt.Sprintf("published '%v' to %v", msg, tags)
}

//
func handleList(client mist.Proxy, args []string) string {
	fmt.Println("HANDLE  LIST")

  //
	return fmt.Sprintf("subscribed to %v", client.List())
}
