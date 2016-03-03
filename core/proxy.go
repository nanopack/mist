package mist

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nanopack/mist/subscription"
)

type (

	//
	Proxy struct {
		sync.Mutex

		subscriptions subscription.Subscriptions
		done  chan bool
		check chan Message
		pipe  chan Message
		id            uint32
	}
)

//
func NewProxy(buffer int) (p Proxy) {

	fmt.Println("NEW PROXY!")

	//
	p = Proxy{
		subscriptions: subscription.NewNode(),
		done:          make(chan bool),
		check:         make(chan Message, buffer),
		pipe:          make(chan Message),
		id:            atomic.AddUint32(&uid, 1),
	}

	p.connect()

	return
}

// connect
func (p *Proxy) connect() {

	fmt.Println("PROXY CONNECT!")

	// this gofunc handles matching messages to subscriptions for the proxy
	go p.handleMessages()

	// add the proxy to mists list of subscribers
	subscribers[p.id] = p
}

//
func (p *Proxy) handleMessages() {

	defer func() {
		close(p.check)
		close(p.pipe)
	}()

	//
	for {
		select {

		//
		case msg := <-p.check:

			p.Lock()
			match := p.subscriptions.Match(msg.Tags)
			p.Unlock()

			if match {
				p.pipe <- msg
			}

			//
		case <-p.done:
			fmt.Println("PROXY SHOULD STOP!")
			return
		}
	}
}

//
func (p *Proxy) Ping() {
	fmt.Println("PROXY PING!")
}

// Subscribe
func (p *Proxy) Subscribe(tags []string) error {
	fmt.Println("PROXY SUBSCRIBE!")

	//
	if len(tags) == 0 {
		return nil
	}

	sort.Sort(sort.StringSlice(tags))

	//
	p.Lock()
	p.subscriptions.Add(tags)
	p.Unlock()

	//
	return publish(p.id, tags, "subscribe")
}

// Unsubscribe
func (p *Proxy) Unsubscribe(tags []string) error {
	fmt.Println("PROXY UNSUBSCRIBE!")

	//
	if len(tags) == 0 {
		return nil
	}

	sort.Sort(sort.StringSlice(tags))

	//
	p.Lock()
	p.subscriptions.Remove(tags)
	p.Unlock()

	//
	return publish(p.id, tags, "unsubscribe")
}

// Publish
func (p *Proxy) Publish(tags []string, data string) error {
	fmt.Println("PROXY PUBLISH!")

	//
	return publish(p.id, tags, data)
}

// Sends a message with delay
func (p *Proxy) PublishAfter(tags []string, data string, delay time.Duration) {

	//
	go func() {
		<-time.After(delay)
		if err := publish(p.id, tags, data); err != nil {
			// write this to a log?
		}
	}()
}

// List
func (p *Proxy) List() [][]string {
	fmt.Println("PROXY LIST!")

	//
	return p.subscriptions.ToSlice()
}

//
func (p *Proxy) Close() {

	fmt.Println("PROXY CLOSE!")

	// this closes the goroutine that is matching messages to subscriptions
	close(p.done)

	// remove the local p from mists list of subscribers/replicators/internal
	delete(subscribers, p.id)

	// send out the unsubscribe message to anyone listening
	for _, subscription := range p.subscriptions.ToSlice() {
		publish(p.id, subscription, "unsubscribe")
	}
}

// Returns all messages that have sucessfully matched the list of subscriptions
// that this proxy has subscribed to
func (p *Proxy) Messages() <-chan Message {
	fmt.Println("PROXY MESSAGES!")
	return p.pipe
}
