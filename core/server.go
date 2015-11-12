// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// build a small applicationController so that we don't have to play with a large
// switch statement
type (
	Handler struct {
		ArgCount int
		Handle   func(Client, []string) string
	}

	mistReader struct {
		reader *bufio.Reader
		err    error
		cmd    []string
	}
)

var (
	commandMap = map[string]Handler{
		"list":               {0, handleList},
		"subscribe":          {1, handleSubscribe},
		"unsubscribe":        {1, handleUnubscribe},
		"publish":            {2, handlePublish},
		"ping":               {0, handlePing},
		"enable-replication": {0, handleEnableReplication},
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
	return ""
}

//
func handleUnubscribe(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Unsubscribe(tags)
	return ""
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
func handlePublish(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Publish(tags, args[1])
	return ""
}

func handleEnableReplication(client Client, args []string) string {
	client.(EnableReplication).EnableReplication()
	return ""
}

// start starts a tcp server listening on the specified address (default 127.0.0.1:1445),
// it then continually reads from the server handling any incoming connections
func (m *Mist) Listen(address string, additinal map[string]Handler) (net.Listener, error) {
	if address == "" {
		address = "127.0.0.1:1445"
	}
	serverSocket, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	// copy the original commands
	commands := make(map[string]Handler)
	for key, value := range commandMap {
		commands[key] = value
	}

	// add additional commands into the map
	for key, value := range additinal {
		commands[key] = value
	}

	go func() {
		defer serverSocket.Close()
		// Continually listen for any incoming connections.
		for {
			conn, err := serverSocket.Accept()
			if err != nil {
				// what should we do with the error?
				return
			}

			// handle each connection individually (non-blocking)
			go m.handleConnection(conn, commands)
		}
	}()
	return serverSocket, nil
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish Handler'
// that is used to publish messages to the data channel of the subscription
func (m *Mist) handleConnection(conn net.Conn, commands map[string]Handler) {

	// create a new client to match with this connection

	client := NewLocalClient(m, 0)

	// make a done channel
	done := make(chan bool)

	// clean up everything that we have setup
	defer func() {
		conn.Close()
		client.Close()
		close(done)
	}()

	// create a 'publish Handler'
	go func() {
		for {

			// when a message is recieved on the subscriptions channel write the message
			// to the connection
			select {
			case msg := <-client.Messages():

				if _, err := conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(msg.Tags, ","), msg.Data))); err != nil {
					break
				}

			// return if we are done
			case <-done:
				return
			}
		}
	}()

	reader := newMistReader(conn)

	for reader.Next() {
		cmd := reader.Command()

		handler, found := commands[cmd[0]]

		var response string
		args := cmd[1:]

		switch {
		case !found:
			response = fmt.Sprintf("error Unknown command '%v'", cmd)
		case handler.ArgCount != len(args):
			response = fmt.Sprintf("error Incorrect number of arguments for '%v'", cmd)
		default:
			response = handler.Handle(client, args)
		}

		if response != "" {
			// Is it safe to send from 2 gorountines at the same time?
			if _, err := conn.Write([]byte(response + "\n")); err != nil {
				break
			}
		}
	}
	// what should we do with the error?
	reader.Error()
}

func newMistReader(conn net.Conn) *mistReader {
	reader := &mistReader{
		reader: bufio.NewReader(conn),
	}

	return reader
}

func (r *mistReader) Next() bool {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		r.err = err
		return false
	}
	line = strings.TrimSuffix(line, "\n")

	r.cmd = strings.SplitN(line, " ", 3)
	return true
}

func (r *mistReader) Command() []string {
	return r.cmd
}

func (r *mistReader) Error() error {
	return r.err
}
