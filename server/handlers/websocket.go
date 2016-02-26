package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/nanopack/mist/core"
)

//
func GenerateWSCommands() map[string]mist.WSHandler {
	return map[string]mist.WSHandler{
		"ping":        {0, handleWSPing},
		"subscribe":   {0, handleWSSubscribe},
		"unsubscribe": {0, handleWSUnubscribe},
		"list":        {0, handleWSList},
	}
}

//
func handleWSPing(client mist.Client, frame []byte, write chan<- string) error {
	write <- "{\"success\":true,\"command\":\"ping\"}"
	return nil
}

//
func handleWSSubscribe(client mist.Client, frame []byte, write chan<- string) error {
	tags := struct {
		Tags []string `json:"tags"`
	}{}

	// error would already be caught by unmarshalling the command
	if err := json.Unmarshal(frame, &tags); err != nil {
		fmt.Println("BUNK!", err)
	}

	//
	client.Subscribe(tags.Tags)

	write <- "{\"success\":true,\"command\":\"subscribe\"}"

	return nil
}

//
func handleWSUnubscribe(client mist.Client, frame []byte, write chan<- string) error {
	tags := struct {
		Tags []string `json:"tags"`
	}{}

	// error would already be caught by unmarshalling the command
	if err := json.Unmarshal(frame, &tags); err != nil {
		fmt.Println("BUNK!", err)
	}

	//
	client.Unsubscribe(tags.Tags)

	write <- "{\"success\":true,\"command\":\"unsubscribe\"}"

	return nil
}

//
func handleWSList(client mist.Client, frame []byte, write chan<- string) (err error) {

	//
	list := struct {
		Subscriptions [][]string `json:"subscriptions"`
		Command       string     `json:"command"`
		Success       bool       `json:"success"`
	}{}

	if list.Subscriptions, err = client.List(); err != nil {
		return err
	}

	list.Command = "list"
	list.Success = true

	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}

	//
	write <- string(bytes)

	return
}
