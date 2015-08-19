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
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/pagodabox/nanobox-golang-stylish"
)

//
type (

	// A Client has a connection to the mist server, a data channel from which it
	// receives messages from the server, ad a host and port to use when connecting
	// to the server
	Client struct {
		conn net.Conn     // the connection the mist server
		done chan bool    // the channel to indicate that the connection is closed
		Data chan Message // the channel that mist server 'publishes' updates to
		Host string       // the connection host for where mist server is running
		Port string       // the connection port for where mist server is running
	}
)

// Connect attempts to connect to a running mist server at the clients specified
// host and port.
func (c *Client) Connect() error {
	fmt.Printf(stylish.Bullet("Attempting to connect to mist..."))

	// number of seconds/attempts to try when failing to conenct to mist server
	maxRetries := 60

	// attempt to connect to the host:port
	for i := 0; i < maxRetries; i++ {
		if conn, err := net.Dial("tcp", c.Host+":"+c.Port); err != nil {

			// max number of attempted retrys failed...
			if i >= maxRetries {
				fmt.Printf(stylish.Error("mist connection failed", "The attempted connection to mist failed. This shouldn't effect any running processes, however no output should be expected"))
				return err
			}
			fmt.Printf("\r   Connection failed! Retrying (%v/%v attempts)...", i, maxRetries)

			// upon successful connection, set the clients connection (conn) to the tcp
			// connection that was established with the server
		} else {
			fmt.Printf(stylish.SubBullet("- Connection established"))
			fmt.Printf(stylish.Success())
			c.conn = conn
			break
		}

		//
		time.Sleep(1 * time.Second)
	}

	// create a channel on which to publish messages received from mist server
	c.Data = make(chan Message)

	// continually read from conn, forwarding the data onto the clients data channel
	go func() {
		defer close(c.Data)

		r := bufio.NewReader(c.conn)
		for {
			bytes, err := r.ReadBytes('\n')
			if err != nil {
				// do we need to log the error?
				return
			}

			// create a new message
			msg := Message{}

			// unmarshal the raw message into a mist message
			if err := json.Unmarshal(bytes, &msg); err != nil {
				// or send the error if there is one
				// msg = Message{Tags: []string{"err"}, Data: err.Error()}
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
	_, err := c.conn.Write([]byte("subscriptions\n"))

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
