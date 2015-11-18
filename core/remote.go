// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"fmt"
	"github.com/nanopack/mist/subscription"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

//
type (

	// A remoteSubscriber represents a connection to the mist server
	remoteSubscriber struct {
		sync.Mutex

		subscriptions Subscriptions      // local copy of subscriptions
		conn          io.ReadWriteCloser // the connection the mist server
		done          chan error         // the channel to indicate that the connection is closed
		waiting       []chan remoteReply // all client waiting for a response
		data          chan Message       // the channel that mist server 'publishes' updates to
		open          bool               // flag that indicates that the conenction should reestablish
		replicated    bool               // is replication enabled on this connection
	}

	remoteReply struct {
		value interface{}
		err   error
	}
	nothing struct{}
)

// Connect attempts to connect to a running mist server at the clients specified
// host and port.
func NewRemoteClient(address string) (Client, error) {
	client := &remoteSubscriber{
		subscriptions: subscription.NewNode(),
		done:          make(chan error),
		waiting:       make([]chan remoteReply, 0),
		data:          make(chan Message),
		open:          true,
		conn:          nothing{},
	}
	return client, client.loop(address)
}

func (client *remoteSubscriber) loop(address string) error {
	// every time we disconnect, we want to reconnect
	conn, err := net.Dial("tcp", address)
	client.conn = conn

	// background the loop so we can return any inital connection errors
	go func() {
		for client.open {
			conn, err := net.Dial("tcp", address)
			client.conn = conn
			if err != nil || !client.open {
				<-time.After(time.Second)
				continue
			}
			// reenable replication
			if client.replicated {
				client.async("enable-replication\n")
			}

			// send all saved subscriptions across the channel
			client.Lock()
			for _, subscription := range client.subscriptions.ToSlice() {
				client.async("subscribe %v\n", strings.Join(subscription, ","))
			}
			client.Unlock()

			reader := newMistReader(conn)

			for reader.Next() {
				cmd := reader.Command()
				switch cmd[0] {
				case "publish":
					msg := Message{
						Tags: strings.Split(cmd[1], ","),
						Data: cmd[2],
					}
					select {
					case client.data <- msg:
					case <-client.done:
					}
				case "pong":
					client.Lock()
					wait := client.waiting[0]
					client.waiting = client.waiting[1:]
					client.Unlock()

					wait <- remoteReply{"pong", nil}
				case "list":
					client.Lock()
					wait := client.waiting[0]
					client.waiting = client.waiting[1:]
					client.Unlock()
					list := [][]string{strings.Split(cmd[1], ",")}
					if len(cmd) == 3 {
						cmd := strings.Split(cmd[2], " ")
						for _, subscription := range cmd {
							list = append(list, strings.Split(subscription, ","))
						}
					}
					wait <- remoteReply{list, nil}
				case "error":
					// close the connection as something is seriously wrong,
					// it will reconnect and and continue on
					conn.Close()

					waiting := make([]chan remoteReply, 0)

					client.Lock()
					waiting, client.waiting = client.waiting, waiting
					client.Unlock()

					for _, wait := range waiting {
						wait <- remoteReply{"", fmt.Errorf("%v", cmd[0])}
					}

				}
			}
		}
	}()
	return err
}

// List requests a list of current mist subscriptions from the server
func (client *remoteSubscriber) List() ([][]string, error) {
	remoteReply := client.sync("list\n")
	if remoteReply.value == nil {
		return nil, remoteReply.err
	}
	return remoteReply.value.([][]string), remoteReply.err
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (client *remoteSubscriber) Subscribe(tags []string) error {
	client.Lock()
	active := client.subscriptions.Match(tags)
	client.subscriptions.Add(tags)
	client.Unlock()
	if len(tags) == 0 {
		return nil
	}
	if !active {
		return client.async("subscribe %v\n", strings.Join(tags, ","))
	}
	return nil
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (client *remoteSubscriber) Unsubscribe(tags []string) error {
	client.Lock()
	client.subscriptions.Remove(tags)
	active := client.subscriptions.Match(tags)
	client.Unlock()
	if len(tags) == 0 {
		return nil
	}
	if !active {
		return client.async("unsubscribe %v\n", strings.Join(tags, ","))
	}
	return nil
}

// Publish sends a message to the mist server to be published to all subscribed clients
func (client *remoteSubscriber) Publish(tags []string, data string) error {
	if len(tags) == 0 {
		return nil
	}
	return client.async("publish %v %v\n", strings.Join(tags, ","), data)
}

// PublishDelay sends a message to the mist server to be published to all subscribed clients
// with delay
func (client *remoteSubscriber) PublishDelay(tags []string, data string, delay time.Duration) error {
	go func() {
		time.Sleep(delay)
		client.Publish(tags, data)
	}()
	return nil
}

// Ping pong the server
func (client *remoteSubscriber) Ping() error {
	remoteReply := client.sync("ping\n")
	return remoteReply.err
}

//
func (client *remoteSubscriber) Messages() <-chan Message {
	return client.data
}

func (client *remoteSubscriber) EnableReplication() error {
	client.replicated = true
	return client.async("enable-replication\n")
}

// Close closes the client data channel and the connection to the server
func (client *remoteSubscriber) Close() error {
	// we need to do it in this order in case the goroutine is stuck waiting for
	// more data from the socket
	client.open = false
	client.conn.Close()
	close(client.done)

	return nil
}

func (client *remoteSubscriber) sync(command string) remoteReply {
	wait := make(chan remoteReply, 1)
	client.Lock()
	if _, err := fmt.Fprintf(client.conn, command); err != nil {
		client.Unlock()
		return remoteReply{nil, err}
	}
	client.waiting = append(client.waiting, wait)
	client.Unlock()
	return <-wait
}

func (client *remoteSubscriber) async(format string, args ...interface{}) error {
	_, err := fmt.Fprintf(client.conn, format, args...)
	return err
}

func (nothing) Read([]byte) (int, error)  { return 0, fmt.Errorf("closed") }
func (nothing) Write([]byte) (int, error) { return 0, fmt.Errorf("closed") }
func (nothing) Close() error              { return nil }
