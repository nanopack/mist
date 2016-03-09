package clients

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/nanopack/mist/core"
)

//
type (

	// A tcp represents a TCP connection to the mist server
	tcp struct {
		host     string
		conn     net.Conn          // the connection the mist server
		encoder  *json.Encoder     //
		messages chan mist.Message // the channel that mist server 'publishes' updates to
	}
)

// New attempts to connect to a running mist server at the clients specified
// host and port.
func New(host string) (*tcp, error) {
	client := &tcp{
		host:     host,
		messages: make(chan mist.Message),
	}

	return client, client.connect()
}

// connect dials the remote mist server and handles any incoming responses back
// from mist
func (c *tcp) connect() error {

	// attempt to connect to the server
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return err
	}

	//
	c.conn = conn

	//
	c.encoder = json.NewEncoder(c.conn)

	// connection loop (blocking); continually read off the connection. Once something
	// is read, check to see if it's a message the client understands to be one of
	// its commands. If so attempt to execute the command.
	go func() {
		decoder := json.NewDecoder(conn)

		for decoder.More() {

			//
			msg := mist.Message{}

			// decode an array value (Message)
			if err := decoder.Decode(&msg); err != nil {
				// log this error and continue?
				continue
			}

			//
			c.messages <- msg
		}
	}()

	return nil
}

// Ping the server
func (c *tcp) Ping() error {
	return c.encoder.Encode(&mist.Message{Command: "ping"})
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (c *tcp) Subscribe(tags []string) error {

	//
	if len(tags) == 0 {
		return fmt.Errorf("Unable to subscribe - missing tags")
	}

	//
	return c.encoder.Encode(&mist.Message{Command: "subscribe", Tags: tags})
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (c *tcp) Unsubscribe(tags []string) error {

	//
	if len(tags) == 0 {
		return fmt.Errorf("Unable to unsubscribe - missing tags")
	}

	//
	return c.encoder.Encode(&mist.Message{Command: "unsubscribe", Tags: tags})
}

// Publish sends a message to the mist server to be published to all subscribed
// clients
func (c *tcp) Publish(tags []string, data string) error {

	//
	if len(tags) == 0 {
		return fmt.Errorf("Unable to publish - missing tags")
	}

	//
	if data == "" {
		return fmt.Errorf("Unable to publish - missing data")
	}

	//
	return c.encoder.Encode(&mist.Message{Command: "publish", Tags: tags, Data: data})
}

// PublishAfter sends a message to the mist server to be published to all subscribed
// clients after a specified delay
func (c *tcp) PublishAfter(tags []string, data string, delay time.Duration) error {
	go func() {
		<-time.After(delay)
		c.Publish(tags, data)
	}()
	return nil
}

// List requests a list from the server of the tags this client is subscribed to
func (c *tcp) List() error {
	return c.encoder.Encode(&mist.Message{Command: "list"})
}

// Close closes the client data channel and the connection to the server
func (c *tcp) Close() {

	// we need to do it in this order in case the goroutine is stuck waiting for
	// more data from the socket
	c.conn.Close()
	close(c.messages)
}

// Messages
func (c *tcp) Messages() <-chan mist.Message {
	return c.messages
}
