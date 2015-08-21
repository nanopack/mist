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

//
type (

	// A Client has a connection to the mist server, a data channel from which it
	// receives messages from the server, ad a host and port to use when connecting
	// to the server
	Client struct {
		conn net.Conn        // the connection the mist server
		done chan bool       // the channel to indicate that the connection is closed
		pong chan bool       // the channel for ping responses
		list chan [][]string // the channel for subscription listing
		Data chan Message    // the channel that mist server 'publishes' updates to
	}
)

// Connect attempts to connect to a running mist server at the clients specified
// host and port.
func (m *Mist) Connect(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	client := Client{
		done: make(chan bool),
		pong: make(chan bool),
		list: make(chan [][]string),
	}
	client.conn = conn

	// create a channel on which to publish messages received from mist server
	client.Data = make(chan Message)

	// continually read from conn, forwarding the data onto the clients data channel
	go func() {
		defer close(client.Data)

		r := bufio.NewReader(client.conn)
		for {
			var listChan chan [][]string
			var pongChan chan bool
			var dataChan chan Message

			line, err := r.ReadString('\n')
			if err != nil {
				// do we need to log the error?
				return
			}
			line = strings.TrimSuffix(line, "\n")

			// create a new message
			var msg Message
			var list [][]string

			split := strings.SplitN(line, " ", 2)

			switch split[0] {
			case "publish":
				split := strings.SplitN(split[1], " ", 2)
				msg = Message{
					Tags: strings.Split(split[0], ","),
					Data: split[1],
				}
				dataChan = client.Data
			case "pong":
				pongChan = client.pong
			case "list":
				split := strings.Split(split[1], " ")
				list = make([][]string, len(split))
				for idx, subscription := range split {
					list[idx] = strings.Split(subscription, ",")
				}
				listChan = client.list
			case "error":
				// need to report the error somehow
				// close the connection as something is seriously wrong
				client.Close()
				return
			}

			// send the message on the client channel, or close if this connection is done
			select {
			case listChan <- list:
			case pongChan <- true:
			case dataChan <- msg:
			case <-client.done:
				return
			}
		}
	}()

	return &client, nil
}

// Publish sends a message to the mist server to be published to all subscribed clients
func (client *Client) Publish(tags []string, data string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := client.conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(tags, ","), data)))

	return err
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (client *Client) Subscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := client.conn.Write([]byte("subscribe " + strings.Join(tags, ",") + "\n"))

	return err
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (client *Client) Unsubscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := client.conn.Write([]byte("unsubscribe " + strings.Join(tags, ",") + "\n"))

	return err
}

// Subscriptions requests a list of current mist subscriptions from the server
func (client *Client) Subscriptions() ([][]string, error) {
	if _, err := client.conn.Write([]byte("list\n")); err != nil {
		return nil, err
	}
	return <-client.list, nil
}

// Ping pong the server
func (client *Client) Ping() error {
	if _, err := client.conn.Write([]byte("ping\n")); err != nil {
		return err
	}
	// wait for the pong to come back
	<-client.pong
	return nil
}

// Close closes the client data channel and the connection to the server
func (client *Client) Close() error {
	// we need to do it in this order in case the goroutine is stuck waiting for
	// more data from the socket
	err := client.conn.Close()
	close(client.done)
	return err
}
