//
package mist

import (
	"errors"
	"fmt"
	"sort"
)

//
const (
	DEFAULT_ADDR = "127.0.0.1:1445"
)

var (
	NotSupported = errors.New("Unable to perform action: command not supported")
	subscribers = make(map[uint32]*Proxy)
	uid        uint32
)

//
type (

	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		Tags []string `json:"tags"`
		Data string   `json:"data"`
	}
)

// Publish publishes to both subscribers, and to replicators
func publish(pid uint32, tags []string, data string) error {
	fmt.Println("MIST publish!")

	// is this an error? or just something we need to ignore
	if len(tags) == 0 {
		return nil
	}

	//
	message := Message{Tags: tags, Data: data}

	// we do this here so that the tags come pre sorted for the clients
	sort.Sort(sort.StringSlice(message.Tags))

	// this should be more optimized, but it might not be an issue unless thousands
	// of clients are using mist.
	go func() {
		for _, subscriber := range subscribers {
			select {
			case <-subscriber.done:
			default:

				// dont sent this message to the publisher who just sent it
				if subscriber.id == pid {
					fmt.Println("SAME GUY!")
					continue
				}

				fmt.Println("publish")
				subscriber.check <- message
			}
		}
	}()

	return nil
}
