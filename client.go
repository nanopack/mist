// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"strings"
)

//
type (

	// Client...
	Client struct {
		conn net.Conn
		Data chan Message
	}
)

//
func (c *Client) Connect(host, port string) (*Client, error) {

	//
	c.Data = make(chan Message)

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}

	c.conn = conn

	go func() {
		for {
			bsize := make([]byte, 4)
			if _, err := io.ReadFull(c.conn, bsize); err != nil {
				c.Data <- Message{Tags: []string{"ERROR"}, Data: err.Error()}
				close(c.Data)
			}

			n := binary.LittleEndian.Uint32(bsize)

			b := make([]byte, n)
			if _, err := io.ReadFull(c.conn, b); err != nil {
				c.Data <- Message{Tags: []string{"ERROR"}, Data: err.Error()}
				close(c.Data)
			}

			msg := Message{}

			if err := json.Unmarshal(b, &msg); err != nil {
				c.Data <- Message{Tags: []string{"ERROR"}, Data: err.Error()}
				close(c.Data)
			}

			c.Data <- msg
		}
	}()

	return c, nil
}

// Subscribe
func (c *Client) Subscribe(tags []string) ([]string, error) {
	if _, err := c.conn.Write([]byte("subscribe " + strings.Join(tags, ",") + "\n")); err != nil {
		return nil, err
	}

	return tags, nil
}

// Unsubscribe
func (c *Client) Unsubscribe(tags []string) error {
	if _, err := c.conn.Write([]byte("subscribe " + strings.Join(tags, ",") + "\n")); err != nil {
		return err
	}

	return nil
}

// Subscriptions
func (c *Client) Subscriptions() error {
	if _, err := c.conn.Write([]byte("subscriptions\n")); err != nil {
		return err
	}

	return nil
}
