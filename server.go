// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"strings"

	"github.com/pagodabox/nanobox-golang-stylish"
)

// build a small applicationController so that we don't have to play with
type (
	handler struct {
		argCount int
		handle   func(*Mist, []string) error
	}
)

var (
	commandMap = map[string]handler{
		"subscribe":     {1, handleSubscribe},
		"unsubscribe":   {1, handleUnubscribe},
		"subscriptions": {0, handleSubscriptions},
		"publish":       {2, handlePublish},
	}
)

func handleSubscribe(m *Mist, args []string) error {
	m.Subscribe(args[0])
	return nil
}
func handleUnubscribe(m *Mist, args []string) error {
	m.Unsubscribe(args[0])
	return nil
}
func handleSubscriptions(m *Mist, args []string) error {
	// don't know how this works yet
	m.List()
	return nil
}
func handlePublish(m *Mist, args []string) error {
	m.Publish(args[0], args[1])
	return nil
}

// start starts a tcp server listening on the specified port (default 1445), it
// then continually reads from the server handling any incoming connections
func (m *Mist) start() {
	m.log.Info(stylish.Bullet("Starting mist server..."))

	serverSocket, err := net.Listen("tcp", ":"+m.port)
	if err != nil {
		m.log.Error("%+v\n", err)
		return
	}
	m.log.Info(stylish.Bullet("Mist listening on port " + m.port))

	//
	go func(serverSocket net.Listener) {
		//
		defer serverSocket.Close()

		// Continually listen for any incoming connections.
		for {
			conn, err := serverSocket.Accept()
			if err != nil {
				m.log.Error("%+v\n", err)
				return
			}

			// handle each connection individually (non-blocking)
			go m.handleConnection(conn)
		}
	}(serverSocket)
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish handler'
// that is used to publish messages to the data channel of the subscription
func (m *Mist) handleConnection(conn net.Conn) {
	m.log.Debug("[MIST :: SERVER] New connection detected: %+v\n", conn)

	// create a new subscription
	sub := Subscription{
		Sub: make(chan Message),
	}

	// make a done channel
	done := make(chan bool)

	// clean up everything that we have setup
	defer func() {
		conn.Close()
		m.Unsubscribe(sub)
		close(done)
		// the channel is not closed here, because this is left up to the client
		// close(sub.Sub)
	}()

	// create a 'publish handler'
	go func() {
		for {

			// when a message is recieved on the subscriptions channel write the message
			// to the connection
			select {
			case msg := <-sub.Sub:

				bytes, err := json.Marshal(msg)
				if err != nil {
					m.log.Error("[MIST :: SERVER] Failed to marshal message: %v\n", err)
					continue
				}

				// 15 is '\n' or a newline character
				if _, err := conn.Write(append(bytes, 15)); err != nil {
					break
				}

			// return if we are done
			case <-done:
				return
			}
		}
	}()

	//
	r := bufio.NewReader(conn)

	//
	for {

		// read messages coming across the tcp channel
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			// some unexpected error happened
			m.log.Error("[MIST :: SERVER] Error reading stream: %+v\n", err.Error())
			return
		}

		// this is the general format of the commands that are accepted
		// ["cmd" ,"tag,tag2", "all the rest"]
		split := strings.SplitN(line, " ", 3)
		cmd := split[0]

		handler, found := commandMap[cmd]

		if !found {
			m.log.Error("[MIST :: SERVER] Unknown command: %+v\n", cmd)
			continue
		}

		args := split[1:]
		if handler.argCount != len(args) {
			m.log.Error("[MIST :: SERVER] incorrect number of arguments for `%v`", cmd)
			continue
		}

		if err := handler.handle(m, args); err != nil {
			m.log.Error("[MIST :: SERVER] %v", err)
		}
	}
}
