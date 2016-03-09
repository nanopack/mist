//
package mist

import (
	"fmt"
	"sync"
)

//
const (
	DEFAULT_ADDR = "127.0.0.1:1445"
)

var (

	//
	ErrUnauthorized = fmt.Errorf("Error: Unauthorized action\n")

	mutex       = &sync.Mutex{}
	subscribers = make(map[uint32]*Proxy)
	uid         uint32
)

//
type (

	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		Command string   `json:"command"`
		Tags    []string `json:"tags"`
		Data    string   `json:"data,omitemtpy"`
		Error   string   `json:"error,omitempty"`
	}

	//
	HandleFunc func(*Proxy, Message) error
)

// publish publishes to both subscribers, and to replicators
func publish(pid uint32, tags []string, data string) error {

	//
	if len(tags) == 0 {
		return fmt.Errorf("Failed to publish. Missing tags...")
	}

	// this should be more optimized, but it might not be an issue unless thousands
	// of clients are using mist.
	go func() {
		mutex.Lock()
		for _, subscriber := range subscribers {
			select {
			case <-subscriber.done:
				fmt.Println("Does this ever happen?")

			//
			default:

				// dont sent this message to the publisher who just sent it
				if subscriber.id == pid {
					continue
				}

				//
				msg := Message{Command: "publish", Tags: tags, Data: data}

				// we don't want this operation blocking the range of other subscribers
				// waiting to get messages
				go func(p *Proxy, msg Message) {
					p.check <- msg
				}(subscriber, msg)
			}
		}
		mutex.Unlock()
	}()

	return nil
}

// subscribe adds a proxy to the list of mist subscribers; we need this so that
// we can lock this process incase multiple proxies are subscribing at the same
// time
func subscribe(p *Proxy) {
	mutex.Lock()
	subscribers[p.id] = p
	mutex.Unlock()
}

// unsubscribe removes a proxy from the list of mist subscribers; we need this
// so that we can lock this process incase multiple proxies are unsubscribing at
// the same time
func unsubscribe(pid uint32) {
	mutex.Lock()
	delete(subscribers, pid)
	mutex.Unlock()
}
