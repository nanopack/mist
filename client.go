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
	"fmt"
	"io"
	"net"
	"strings"
	"time"
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

	maxRetries := 60

	// attempt to connect to the host:port
	for i := 0; i < maxRetries; i++ {
		if conn, err := net.Dial("tcp", host + ":" + port); err != nil {

			// max number of attempted retrys failed...
			if i >= maxRetries {
				fmt.Println("[MIST :: CLIENT] Failed to connect...")
				return nil, err
			}
			fmt.Printf("\r[MIST :: CLIENT] Unable to connect, retrying... (%v/%v attempts)", i, maxRetries)

			// connection successful
		} else {
			fmt.Println("\n[MIST :: CLIENT] Connected...")
			c.conn = conn
			break
		}

		//
		time.Sleep(1*time.Second)
	}

	//
	go func() {
		for {
			fmt.Println("for looping!")
			// read the first 4 bytes of the message so we know how long the message
			// is expected to be
			bsize := make([]byte, 4)
			if _, err := io.ReadFull(c.conn, bsize); err != nil {
				fmt.Println(err)
				c.Data <- Message{Tags: []string{"ERROR"}, Data: err.Error()}
				close(c.Data)
				// c.Close()
			}
			fmt.Println("bsize:", bsize)

			// create a buffer that is the length of the expected message
			n := binary.LittleEndian.Uint32(bsize)

			fmt.Println("len:", n)
			// read the length of the message up to the expected bytes
			b := make([]byte, n)
			if _, err := io.ReadFull(c.conn, b); err != nil {
				fmt.Println(err)
				c.Data <- Message{Tags: []string{"ERROR"}, Data: err.Error()}
				close(c.Data)
				// c.Close()
			}
			fmt.Println("b: ",b)
			//
			msg := Message{}

			//
			if err := json.Unmarshal(b, &msg); err != nil {
				fmt.Println(err)
				c.Data <- Message{Tags: []string{"ERROR"}, Data: err.Error()}
				close(c.Data)
				// c.Close()
			}
			fmt.Printf("msg: %#v\n", msg)

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
	if _, err := c.conn.Write([]byte("unsubscribe " + strings.Join(tags, ",") + "\n")); err != nil {
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

// Close
func (c *Client) Close() error {
	close(c.Data)
	// return c.conn.Close()
	return nil
}
