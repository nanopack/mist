// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package mist

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	NotSupported = errors.New("Command is not supported over websockets")

	// we use a map to avoid a nasty switch/case statement
	commandWebsocketMap = map[string]WebsocketHandler{
		"list":        handleWebsocketList,
		"subscribe":   handleWebsocketSubscribe,
		"unsubscribe": handleWebsocketUnubscribe,
		"ping":        handleWebsocketPing,
	}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type (
	WebsocketHandler func([]byte, chan<- string, Client) error
	subscribe        struct {
		Command string   `json:"command"`
		Tags    []string `json:"tags"`
	}

	command struct {
		Command string `json:"command"`
	}

	tagList struct {
		Tags []string `json:"tags"`
	}
	list struct {
		Subscriptions [][]string `json:"subscriptions"`
		Command       string     `json:"command"`
		Success       bool       `json:"success"`
	}
	reply struct {
		Success bool   `json:"success"`
		Command string `json:"command"`
		Error   string `json:"error"`
	}

	websocketClient struct {
		conn        *websocket.Conn
		messages    chan Message
		list        chan list
		ping        chan reply
		subscribe   chan reply
		unsubscribe chan reply
	}
)

func handleWebsocketList(frame []byte, write chan<- string, client Client) error {
	list := list{}
	var err error
	list.Subscriptions, err = client.List()
	if err != nil {
		return err
	}
	list.Command = "list"
	list.Success = true
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	write <- string(bytes)
	return nil
}

func handleWebsocketSubscribe(frame []byte, write chan<- string, client Client) error {
	tags := tagList{}
	// error would already be caught by unmarshalling the command
	json.Unmarshal(frame, &tags)
	client.Subscribe(tags.Tags)
	write <- "{\"success\":true,\"command\":\"subscribe\"}"
	return nil
}

func handleWebsocketUnubscribe(frame []byte, write chan<- string, client Client) error {
	tags := tagList{}
	// error would already be caught by unmarshalling the command
	json.Unmarshal(frame, &tags)
	client.Unsubscribe(tags.Tags)
	write <- "{\"success\":true,\"command\":\"unsubscribe\"}"
	return nil
}

func handleWebsocketPing(frame []byte, write chan<- string, client Client) error {
	write <- "{\"success\":true,\"command\":\"ping\"}"
	return nil
}

//
func GenerateWebsocketUpgrade(mist *Mist, additinal map[string]WebsocketHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// copy the original commands
		commands := make(map[string]WebsocketHandler)
		for key, value := range commandWebsocketMap {
			commands[key] = value
		}

		// add additional commands into the map
		for key, value := range additinal {
			commands[key] = value
		}

		// we don't want this to be buffered
		client := NewLocalClient(mist, 0)

		write := make(chan string)
		done := make(chan bool)
		defer func() {
			client.Close()
			close(done)
		}()

		// the gorilla websocket package must have all writes come from the
		// same process.
		go func() {
			for {
				select {
				case event := <-client.Messages():
					if msg, err := json.Marshal(event); err == nil {
						conn.WriteMessage(websocket.TextMessage, msg)
					}
				case msg := <-write:
					conn.WriteMessage(websocket.TextMessage, []byte(msg))
				case <-done:
					close(write)
					return
				}
			}
		}()

		for {
			messageType, frame, err := conn.ReadMessage()
			if err != nil {
				return
			}
			fmt.Println(string(frame))
			if messageType != websocket.TextMessage {
				write <- "{\"success\":false,\"error\":\"I don't understand binary messages\"}"
				continue
			}

			cmd := command{}
			if err := json.Unmarshal(frame, &cmd); err != nil {
				write <- "{\"success\":false,\"error\":\"Invalid json\"}"
				continue
			}
			command, ok := commandWebsocketMap[cmd.Command]
			if !ok {
				write <- "{\"success\":false,\"error\":\"unknown command\"}"
				continue
			}
			if err := command(frame, write, client); err != nil {
				// I should do something with this error..
				return
			}
		}
	}
}

func NewWebsocketClient(address string, header http.Header) (Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(address, header)
	if err != nil {
		return nil, err
	}

	client := &websocketClient{
		conn:        conn,
		messages:    make(chan Message, 0),
		list:        make(chan list, 0),
		ping:        make(chan reply, 0),
		subscribe:   make(chan reply, 0),
		unsubscribe: make(chan reply, 0),
	}
	go func(client *websocketClient) {
		defer func() {
			close(client.messages)
			close(client.list)
			close(client.ping)
			close(client.subscribe)
			close(client.unsubscribe)
		}()
		for {
			messageType, frame, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if messageType != websocket.TextMessage {
				// how do i report these?
				continue
			}

			cmd := reply{}
			if err := json.Unmarshal(frame, &cmd); err != nil {
				// how do i report these?
				continue
			}
			fmt.Println(cmd, string(frame))
			switch {
			case cmd.Command == "list":
				// TODO: what do we do if list failed?
				listResponse := list{}
				// error would have already been caught above
				json.Unmarshal(frame, &listResponse)
				client.list <- listResponse
			case cmd.Command == "ping":
				client.ping <- cmd
			case cmd.Command == "subscribe":
				client.subscribe <- cmd
			case cmd.Command == "unsubscribe":
				client.unsubscribe <- cmd
			default:
				message := Message{}
				json.Unmarshal(frame, &message)
				client.messages <- message
			}
		}

	}(client)
	return client, nil
}

func (client *websocketClient) List() ([][]string, error) {
	listReq := command{
		Command: "list",
	}
	bytes, err := json.Marshal(listReq)
	if err != nil {
		return nil, err
	}
	client.conn.WriteMessage(websocket.TextMessage, bytes)

	list := <-client.list
	return list.Subscriptions, nil
}

func (client *websocketClient) Subscribe(tags []string) error {
	unsubscribe := subscribe{
		Command: "subscribe",
		Tags:    tags,
	}
	bytes, err := json.Marshal(unsubscribe)
	if err != nil {
		return err
	}
	client.conn.WriteMessage(websocket.TextMessage, bytes)
	return isError(<-client.subscribe)
}

func (client *websocketClient) Unsubscribe(tags []string) error {
	unsubscribe := subscribe{
		Command: "unsubscribe",
		Tags:    tags,
	}
	bytes, err := json.Marshal(unsubscribe)
	if err != nil {
		return err
	}
	client.conn.WriteMessage(websocket.TextMessage, bytes)
	return isError(<-client.unsubscribe)
}

func (client *websocketClient) Publish(tags []string, data string) error {
	return NotSupported
}

func (client *websocketClient) PublishDelay(tags []string, data string, delay time.Duration) error {
	return NotSupported
}

func (client *websocketClient) Ping() error {
	ping := command{
		Command: "ping",
	}
	bytes, err := json.Marshal(ping)
	if err != nil {
		return err
	}
	client.conn.WriteMessage(websocket.TextMessage, bytes)
	return isError(<-client.ping)
}

func (client *websocketClient) Messages() <-chan Message {
	return client.messages
}

func (client *websocketClient) Close() error {
	return client.conn.Close()
}

func isError(reply reply) error {
	if !reply.Success {
		return errors.New(reply.Error)
	}
	return nil
}
