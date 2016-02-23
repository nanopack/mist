package mist

import (
	"fmt"
	// "io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/nanopack/mist/subscription"
	"github.com/nanopack/mist/util"
)

//
type (

	// A remoteClient represents a connection to the mist server
	remoteClient struct {
		sync.Mutex

		address       string
		subscriptions Subscriptions      // local copy of subscriptions
		conn          net.Conn           // the connection the mist server
		done          chan error         // the channel to indicate that the connection is closed
		waiting       []chan remoteReply // all client waiting for a response
		data          chan Message       // the channel that mist server 'publishes' updates to
		open          bool               // flag that indicates that the conenction should reestablish
		attempts      int
		replicated    bool // is replication enabled on this connection
	}

	remoteReply struct {
		value interface{}
		err   error
	}

	// nothing struct{}
)

// Connect attempts to connect to a running mist server at the clients specified
// host and port.
func NewRemoteClient(address string) (Client, error) {
	client := &remoteClient{
		subscriptions: subscription.NewNode(),
		done:          make(chan error),
		waiting:       make([]chan remoteReply, 0),
		data:          make(chan Message),
		open:          false,
		attempts:      0,
		address:       address,
	}

	return client, client.connect(address)
}

// connect
func (client *remoteClient) connect(address string) error {

	// attempt an initial connection to the server
	conn, err := net.Dial("tcp", client.address)
	if err != nil {
		return err
	}

	client.conn = conn
	client.open = true

	// keep the connection open
	go client.loop()

	return nil
}

//
func (client *remoteClient) reconnect() {

	// attempt to reconnect
	conn, err := net.Dial("tcp", client.address)

	//
	switch {

	//
	case err != nil && client.attempts < 10:
		fmt.Printf("connection failed... attempting to reconnect %v/10\r", client.attempts)
		<-time.After(time.Second)
		client.attempts += 1
		client.reconnect()

	//
	case err != nil:
		fmt.Printf("unable to connect to server...")

	//
	default:
		client.conn = conn
		client.open = true

		//
		go client.loop()
	}
}

//
func (client *remoteClient) loop() {

	//
	for client.open {

		// reenable replication
		// if client.replicated {
		// 	client.async("enable-replication\n")
		// }

		// send all saved subscriptions across the channel
		client.Lock()
		for _, subscription := range client.subscriptions.ToSlice() {
			client.async("subscribe %v\n", strings.Join(subscription, ","))
		}
		client.Unlock()

		reader := util.NewReader(client.conn)
		for reader.Next() {

			cmd := reader.Input

			//
			switch cmd[0] {

			//
			case "publish":
				msg := Message{
					Tags: strings.Split(cmd[1], ","),
					Data: cmd[2],
				}
				select {
				case client.data <- msg:
				case <-client.done:
				}

			//
			case "pong":
				client.Lock()
				wait := client.waiting[0]
				client.waiting = client.waiting[1:]
				client.Unlock()

				wait <- remoteReply{"pong", nil}

			//
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

			//
			case "error":
				// close the connection as something is seriously wrong,
				// it will reconnect and and continue on
				client.conn.Close()

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

	// if the client ever disconnects, attempt to reconnect; may want to put in a
	// limit here so that it doesn't try and connect forever...
	if !client.open {
		client.reconnect()
	}
}

// List requests a list of current mist subscriptions from the server
func (client *remoteClient) List() ([][]string, error) {
	remoteReply := client.sync("list\n")
	if remoteReply.value == nil {
		return nil, remoteReply.err
	}
	return remoteReply.value.([][]string), remoteReply.err
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (client *remoteClient) Subscribe(tags []string) error {
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
func (client *remoteClient) Unsubscribe(tags []string) error {
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
func (client *remoteClient) Publish(tags []string, data string) error {
	if len(tags) == 0 {
		return nil
	}
	return client.async("publish %v %v\n", strings.Join(tags, ","), data)
}

// PublishAfter sends a message to the mist server to be published to all subscribed clients
// with delay
func (client *remoteClient) PublishAfter(tags []string, data string, delay time.Duration) error {
	go func() {
		time.After(delay)
		client.Publish(tags, data)
	}()
	return nil
}

// Ping pong the server
func (client *remoteClient) Ping() error {
	remoteReply := client.sync("ping\n")
	return remoteReply.err
}

// Close closes the client data channel and the connection to the server
func (client *remoteClient) Close() error {
	// we need to do it in this order in case the goroutine is stuck waiting for
	// more data from the socket
	client.open = false
	client.conn.Close()
	close(client.done)

	return nil
}

//
func (client *remoteClient) Messages() <-chan Message {
	return client.data
}

func (client *remoteClient) EnableReplication() error {
	client.replicated = true
	return client.async("enable-replication\n")
}

//
func (client *remoteClient) sync(command string) remoteReply {
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

//
func (client *remoteClient) async(format string, args ...interface{}) error {
	_, err := fmt.Fprintf(client.conn, format, args...)
	return err
}

//
// func (nothing) Read([]byte) (int, error) { return 0, fmt.Errorf("closed") }
// func (nothing) Write([]byte) (int, error) { return 0, fmt.Errorf("closed") }
// func (nothing) Close() error { return nil }
