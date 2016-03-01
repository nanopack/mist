package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/nanopack/mist/core"
)

type (

	//
	wsClient struct {
		conn        *websocket.Conn
		messages    chan mist.Message
		list        chan list
		ping        chan wsReply
		subscribe   chan wsReply
		unsubscribe chan wsReply
	}

	//
	wsReply struct {
		Success bool   `json:"success"`
		Command string `json:"command"`
		Error   string `json:"error"`
	}

	//
	list struct {
		Subscriptions [][]string `json:"subscriptions"`
		Command       string     `json:"command"`
		Success       bool       `json:"success"`
	}
)

//
func NewWS(address string, header http.Header) (mist.Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(address, header)
	if err != nil {
		return nil, err
	}

	client := &wsClient{
		conn:        conn,
		messages:    make(chan mist.Message, 0),
		list:        make(chan list, 0),
		ping:        make(chan wsReply, 0),
		subscribe:   make(chan wsReply, 0),
		unsubscribe: make(chan wsReply, 0),
	}

	return client, client.connect(address)
}

//
func (client *wsClient) connect(address string) error {

	fmt.Println("WS CONNECT!")

	//
	go func() {

		//
		defer func() {
			close(client.messages)
			close(client.list)
			close(client.ping)
			close(client.subscribe)
			close(client.unsubscribe)
		}()

		for {
			messageType, frame, err := client.conn.ReadMessage()
			if err != nil {
				return
			}

			// how/should these be reported?
			if messageType != websocket.TextMessage {
				continue
			}

			//
			cmd := wsReply{}
			if err := json.Unmarshal(frame, &cmd); err != nil {
				// how do i report these?
				continue
			}

			//
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
				message := mist.Message{}
				json.Unmarshal(frame, &message)
				client.messages <- message
			}
		}
	}()

	return nil
}

//
func (client *wsClient) List() ([][]string, error) {
	cmd := struct {
		Command string `json:"command"`
	}{
		Command: "list",
	}

	bytes, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	client.conn.WriteMessage(websocket.TextMessage, bytes)

	list := <-client.list
	return list.Subscriptions, nil
}

//
func (client *wsClient) Subscribe(tags []string) error {

	sub := struct {
		Command string   `json:"command"`
		Tags    []string `json:"tags"`
	}{
		Command: "subscribe",
		Tags:    tags,
	}

	bytes, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	client.conn.WriteMessage(websocket.TextMessage, bytes)
	return isError(<-client.subscribe)
}

//
func (client *wsClient) Unsubscribe(tags []string) error {

	unsub := struct {
		Command string   `json:"command"`
		Tags    []string `json:"tags"`
	}{
		Command: "unsubscribe",
		Tags:    tags,
	}

	bytes, err := json.Marshal(unsub)
	if err != nil {
		return err
	}
	client.conn.WriteMessage(websocket.TextMessage, bytes)
	return isError(<-client.unsubscribe)
}

//
func (client *wsClient) Publish(tags []string, data string) error {
	return mist.NotSupported
}

//
func (client *wsClient) PublishAfter(tags []string, data string, delay time.Duration) error {
	return mist.NotSupported
}

//
func (client *wsClient) Ping() error {

	cmd := struct {
		Command string `json:"command"`
	}{
		Command: "ping",
	}

	bytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	client.conn.WriteMessage(websocket.TextMessage, bytes)
	return isError(<-client.ping)
}

//
func (client *wsClient) Close() error {
	return client.conn.Close()
}

//
func (client *wsClient) Messages() <-chan mist.Message {
	return client.messages
}

//
func isError(reply wsReply) error {
	if !reply.Success {
		return errors.New(reply.Error)
	}
	return nil
}
