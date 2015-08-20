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
	"net"
	"strings"
)

//
type (

	// A Client has a connection to the mist server, a data channel from which it
	// receives messages from the server, ad a host and port to use when connecting
	// to the server
	Client struct {
		conn  net.Conn     // the connection the mist server
		done  chan bool    // the channel to indicate that the connection is closed
		error chan error   // the channel for error messages
		Data  chan Message // the channel that mist server 'publishes' updates to
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
			line, err := r.ReadString('\n')
			if err != nil {
				// do we need to log the error?
				return
			}

			// create a new message
			msg := Message{}

			split := strings.SplitN(line, " ", 1)

			switch split[0] {
			case "publish":
				if err := json.Unmarshal([]byte(split[1]), &msg); err != nil {
					// or send the error if there is one
					msg = Message{Tags: []string{"err"}, Data: err.Error()}
				}
			case "ok":
				if len(split) > 1 {
					// I don't know how to get this to the right place yet.
				}
			case "error":
				// I don't know how to get this to the right place yet.
				// msg = Message{Tags: []string{"err"}, Data: split[1]}
			}

			// send the message on the client data channel, or close if this connection is done
			select {
			case c.Data <- msg:
			case <-c.done:
				return
			}
		}
	}()

	return nil
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (c *Client) Subscribe(tags []string) ([]string, error) {
	if _, err := c.conn.Write([]byte("subscribe " + strings.Join(tags, ",") + "\n")); err != nil {
		return nil, err
	}

	return tags, nil
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (c *Client) Unsubscribe(tags []string) error {
	_, err := c.conn.Write([]byte("unsubscribe " + strings.Join(tags, ",") + "\n"))

	return err
}

// Subscriptions requests a list of current mist subscriptions from the server
func (c *Client) Subscriptions() error {
	_, err := c.conn.Write([]byte("list\n"))

	return err
}

// Close closes the client data channel and the connection to the server
func (c *Client) Close() error {
	// we need to do it in this order incase the goproc is stuck waiting for
	// more data from the socket
	err := c.conn.Close()
	close(c.done)
	return err
}
