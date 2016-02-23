//
package mist

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"time"
)

//
const (
	DEFAULT_ADDR = "127.0.0.1:1445"
)

var (
	NotSupported = errors.New("Unable to perform action: command not supported")
	InternalErr  = errors.New("Unable to perform action: internal mode enabled")
)

// interfaces
type (

	//
	Client interface {
		Ping() error
		List() ([][]string, error)
		Subscribe(tags []string) error
		Unsubscribe(tags []string) error
		Publish(tags []string, data string) error
		PublishAfter(tags []string, data string, delay time.Duration) error
		Close() error
		Messages() <-chan Message
	}

	Subscriptions interface {
		Add([]string)
		Remove([]string)
		Match([]string) bool
		ToSlice() [][]string
	}

	Replicatable interface {
		EnableReplication() error
	}

	Internalizable interface {
		EnableInternal()
	}
)

//
type (

	//
	Mist struct {
		subscribers map[uint32]*localClient
		replicators map[uint32]*localClient
		internal    map[uint32]*localClient
		next        uint32
	}

	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		internal bool

		Tags []string `json:"tags"`
		Data string   `json:"data"`
	}
)

// creates a new mist
func New() *Mist {
	return &Mist{
		subscribers: make(map[uint32]*localClient),
		replicators: make(map[uint32]*localClient),
		internal:    make(map[uint32]*localClient),
	}
}

// Listen starts a tcp server listening on the specified address (default 127.0.0.1:1445)
// and then continually reads from the server handling any incoming connections;
// this is intentionally non-blocking.
func (mist *Mist) Listen(address string, additional map[string]Handler) (net.Listener, error) {

	//
	if address == "" {
		address = DEFAULT_ADDR
	}

	//
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	// copy the original commands
	commands := make(map[string]Handler)
	for key, value := range serverCommands {
		commands[key] = value
	}

	// add additional commands into the map
	for key, value := range additional {
		commands[key] = value
	}

	// non-blocking
	go func() {

		defer ln.Close()

		// Continually listen for any incoming connections.
		for {

			// accept connections
			conn, err := ln.Accept()
			if err != nil {
				return // what should we do with the error?
			}

			// create a new client for this connection
			client, err := NewLocalClient(mist, 0)
			if err != nil {
				fmt.Println("BONK!")
			}

			// handle each connection individually (non-blocking)
			go newConnection(client, conn, commands)
		}
	}()

	return ln, nil
}

// Publish publishes to both subscribers, and to replicators
func (mist *Mist) Publish(tags []string, data string) error {

	// is this an error? or just something we need to ignore
	if len(tags) == 0 {
		return nil
	}

	message := Message{
		Tags: tags,
		Data: data,
	}

	forward(message, mist.subscribers)
	forward(message, mist.replicators)

	return nil
}

// Replicate publishes to subscribers only
func (mist *Mist) Replicate(tags []string, data string) error {

	// is this an error? or just something we need to ignore
	if len(tags) == 0 {
		return nil
	}

	message := Message{
		Tags: tags,
		Data: data,
	}

	forward(message, mist.subscribers)

	return nil
}

// publish
func (mist *Mist) publish(tags []string, data string) error {

	// is this an error? or just something we need to ignore
	if len(tags) == 0 {
		return nil
	}

	message := Message{
		Tags:     tags,
		internal: true,
		Data:     data,
	}

	forward(message, mist.internal)

	return nil
}

// forward
func forward(msg Message, subscribers map[uint32]*localClient) {

	// we do this here so that the tags come pre sorted for the clients
	sort.Sort(sort.StringSlice(msg.Tags))

	// this should be more optimized, but it might not be an issue unless thousands
	// of clients are using mist.
	for _, localReplicator := range subscribers {
		select {
		case <-localReplicator.done:
		case localReplicator.check <- msg:
			// default:
			// do we drop the message? enqueue it? pull one off the front and then add this one?
		}
	}
}
