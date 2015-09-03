// -*- mode: go; tab-width: 2; indent-tabs-mode: 1; st-rulers: [70] -*-
// vim: ts=4 sw=4 ft=lua noet
//--------------------------------------------------------------------
// @author Daniel Barney <daniel@nanobox.io>
// @copyright 2015, Pagoda Box Inc.
// @doc
//
// @end
// Created :   12 August 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
package mist

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	NotSupported = errors.New("Command is not supported over websockets")

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type (
	subscribe struct {
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

//
func GenerateWebsocketUpgrade(mist *Mist) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
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
			switch cmd.Command {
			case "subscribe":
				tags := tagList{}
				// error would already be caught by unmarshalling the command
				json.Unmarshal(frame, &tags)
				client.Subscribe(tags.Tags)
				write <- "{\"success\":true,\"command\":\"subscribe\"}"
			case "unsubscribe":
				tags := tagList{}
				// error would already be caught by unmarshalling the command
				json.Unmarshal(frame, &tags)
				client.Unsubscribe(tags.Tags)
				write <- "{\"success\":true,\"command\":\"unsubscribe\"}"
			case "list":
				list := list{}
				list.Subscriptions, err = client.List()
				if err != nil {
					// do we need to do something with this error?
					return
				}
				list.Command = "list"
				list.Success = true
				bytes, err := json.Marshal(list)
				if err != nil {
					// Do I need to do something more here?
					return
				}
				write <- string(bytes)
			case "ping":
				write <- "{\"success\":true,\"command\":\"ping\"}"
			default:
				write <- "{\"success\":false,\"error\":\"unknown command\"}"
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
