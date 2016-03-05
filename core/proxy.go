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
		check         chan Message
		done          chan bool
		id            uint32
		pipe          chan Message

		Authorized bool
	}
)

//
func NewProxy() (p *Proxy) {

	//
	p = &Proxy{
		subscriptions: subscription.NewNode(),
		check:         make(chan Message),
		done:          make(chan bool),
		id:            atomic.AddUint32(&uid, 1),
		pipe:          make(chan Message),

		Authorized: false,
	}

	p.connect()

	return
}

// connect
func (p *Proxy) connect() {

	// add the proxy to mists list of subscribers
	subscribers[p.id] = p

	// this gofunc handles matching messages to subscriptions for the proxy
	go p.handleMessages()
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
	fmt.Println("PROXY PING!", p.id)
}

// Subscribe
func (p *Proxy) Subscribe(tags []string) error {

	//
	if !p.Authorized {
		return ErrUnauthorized
	}

	// is this an error?
	if len(tags) == 0 {
		return nil
	}

	sort.Sort(sort.StringSlice(tags))

	//
	p.Lock()
	p.subscriptions.Add(tags)
	p.Unlock()

	//
	return nil
}

// Unsubscribe
func (p *Proxy) Unsubscribe(tags []string) error {

	//
	if !p.Authorized {
		return ErrUnauthorized
	}

	// is this an error?
	if len(tags) == 0 {
		return nil
	}

	sort.Sort(sort.StringSlice(tags))

	//
	p.Lock()
	p.subscriptions.Remove(tags)
	p.Unlock()

	//
	return nil
}

// Publish
func (p *Proxy) Publish(tags []string, data string) error {

	//
	if !p.Authorized {
		return ErrUnauthorized
	}

	//
	return publish(p.id, tags, data)
}

// Sends a message with delay
func (p *Proxy) PublishAfter(tags []string, data string, delay time.Duration) error {

	//
	if !p.Authorized {
		return ErrUnauthorized
	}

	//
	go func() {
		<-time.After(delay)
		if err := publish(p.id, tags, data); err != nil {
			// log this error and continue?
		}
	}()

	return nil
}

// List
func (p *Proxy) List() error {

	//
	if !p.Authorized {
		return ErrUnauthorized
	}

	msg := Message{
		Tags: []string{},
		Data: fmt.Sprint(p.subscriptions.ToSlice()),
	}

	p.pipe <- msg

	//
	return nil
}

// Returns all messages that have sucessfully matched the list of subscriptions
// that this proxy has subscribed to
func (p *Proxy) Messages() <-chan Message {
	return p.pipe
}

//
func (p *Proxy) Close() {

	// this closes the goroutine that is matching messages to subscriptions
	close(p.done)

	// remove the local p from mists list of subscribers/replicators/internal
	delete(subscribers, p.id)
}
