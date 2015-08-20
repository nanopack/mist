// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"bufio"
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
func (c *Client) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn

	// create a channel on which to publish messages received from mist server
	c.Data = make(chan Message)

	// continually read from conn, forwarding the data onto the clients data channel
	go func() {
		defer close(c.Data)

		r := bufio.NewReader(c.conn)
		for {
			var listChan chan [][]string
			var pongChan chan bool
			var dataChan chan Message

			line, err := r.ReadString('\n')
			if err != nil {
				// do we need to log the error?
				return
			}

			// create a new message
			var msg Message
			var list [][]string

			split := strings.SplitN(line, " ", 1)

			switch split[0] {
			case "publish":
				split := strings.SplitN(split[1], " ", 2)
				msg = Message{
					Tags: strings.Split(split[0], ","),
					Data: split[1],
				}
				dataChan = c.Data
			case "pong":
				pongChan = c.pong
			case "list":
				split := strings.Split(split[1], " ")
				list = make([][]string, len(split))
				for idx, subscription := range split {
					list[idx] = strings.Split(subscription, ",")
				}
				listChan = c.list
			}

			// send the message on the client channel, or close if this connection is done
			select {
			case listChan <- list:
			case pongChan <- true:
			case dataChan <- msg:
			case <-c.done:
				return
			}
		}
	}()

	return nil
}

// Publish sends a message to the mist server to be published to all subscribed clients
func (c *Client) Publish(tags []string, data string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := c.conn.Write([]byte("publish " + strings.Join(tags, ",") + "\n"))

	return err
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (c *Client) Subscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := c.conn.Write([]byte("subscribe " + strings.Join(tags, ",") + "\n"))

	return err
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (c *Client) Unsubscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	_, err := c.conn.Write([]byte("unsubscribe " + strings.Join(tags, ",") + "\n"))

	return err
}

// Subscriptions requests a list of current mist subscriptions from the server
func (c *Client) Subscriptions() ([][]string, error) {
	if _, err := c.conn.Write([]byte("list\n")); err != nil {
		return nil, err
	}
	return <-c.list, nil
}

// Ping pong the server
func (c *Client) Ping() error {
	if _, err := c.conn.Write([]byte("ping\n")); err != nil {
		return err
	}
	// wait for the pong to come back
	<-c.pong
	return nil
}

// Close closes the client data channel and the connection to the server
func (c *Client) Close() error {
	// we need to do it in this order incase the goproc is stuck waiting for
	// more data from the socket
	err := c.conn.Close()
	close(c.done)
	return err
}
