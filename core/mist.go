//
package mist

import (
	"errors"
	"sort"
	"time"
)

//
const (
	DEFAULT_ADDR = "127.0.0.1:1445"
)

var (
	InternalErr  = errors.New("Unable to perform action: internal mode enabled")
	NotSupported = errors.New("Unable to perform action: command not supported")
	Self         *Mist
)

//
type (
	//
	Client interface {
		Ping() error
		Subscribe(tags []string) error
		Unsubscribe(tags []string) error
		Publish(tags []string, data string) error
		PublishAfter(tags []string, data string, delay time.Duration) error
		List() ([][]string, error)
		Close() error
		Messages() <-chan Message
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
		subscribers map[uint32]*proxy
		replicators map[uint32]*proxy
		internal    map[uint32]*proxy
		next        uint32
	}

	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		internal bool

		Tags []string `json:"tags"`
		Data string   `json:"data"`
	}

	// tcp handler
	TCPHandler struct {
		NumArgs int
		Handle  func(Client, []string) string
	}

	// websocket handler
	WSHandler struct {
		NumArgs int
		Handle  func(Client, []byte, chan<- string) error
	}
)

// creates a new mist
func New() *Mist {

	Self = &Mist{
		subscribers: make(map[uint32]*proxy),
		replicators: make(map[uint32]*proxy),
		internal:    make(map[uint32]*proxy),
	}

	return Self
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
func forward(msg Message, subscribers map[uint32]*proxy) {

	// we do this here so that the tags come pre sorted for the clients
	sort.Sort(sort.StringSlice(msg.Tags))

	// this should be more optimized, but it might not be an issue unless thousands
	// of clients are using mist.
	for _, subscriber := range subscribers {
		select {
		case <-subscriber.done:
		case subscriber.check <- msg:
			// default:
			// do we drop the message? enqueue it? pull one off the front and then add this one?
		}
	}
}
