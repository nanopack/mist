package clients

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/util"
)

//
type (

	// A tcp represents a TCP connection to the mist server
	tcp struct {
		host     string
		conn     net.Conn          // the connection the mist server
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

	// listen for incoming mist messages
	go func() {
		r := util.NewReader(c.conn)
		for r.Next() {

			//
			var msg mist.Message
			switch r.Input.Cmd {

			// this is any message published from mist
			case "publish", "list":
				msg = mist.Message{Tags: strings.Split(r.Input.Args[0], ","), Data: r.Input.Args[1]}

			// this would be any message sent from this client (ping, subscribe, unsubscribe
			// publish, or list) which we don't care about here
			default:
				continue
			}

			c.messages <- msg
		}
	}()

	return nil
}

// Ping the server
func (c *tcp) Ping() error {
	_, err := fmt.Fprintf(c.conn, "ping\n")
	return err
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (c *tcp) Subscribe(tags []string) error {

	//
	if len(tags) == 0 {
		return fmt.Errorf("Unable to subscribe - missing tags")
	}

	_, err := fmt.Fprintf(c.conn, fmt.Sprintf("subscribe %v\n", strings.Join(tags, ",")))
	return err
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (c *tcp) Unsubscribe(tags []string) error {

	//
	if len(tags) == 0 {
		return fmt.Errorf("Unable to unsubscribe - missing tags")
	}

	_, err := fmt.Fprintf(c.conn, fmt.Sprintf("unsubscribe %v\n", strings.Join(tags, ",")))
	return err
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

	_, err := fmt.Fprintf(c.conn, fmt.Sprintf("publish %v %v\n", strings.Join(tags, ","), data))
	return err
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
	_, err := fmt.Fprintf(c.conn, "list\n")
	return err
}

// Close closes the client data channel and the connection to the server
func (c *tcp) Close() error {

	// we need to do it in this order in case the goroutine is stuck waiting for
	// more data from the socket
	c.conn.Close()
	close(c.messages)

	return nil
}

//
func (c *tcp) Messages() <-chan mist.Message {
	return c.messages
}
