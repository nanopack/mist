package mist

import (
	"fmt"
	"strings"
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
		Pipe          chan Message

		// Authorized bool
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
		Pipe:          make(chan Message),

		// Authorized: false,
	}

	p.connect()

	return
}

// connect
func (p *Proxy) connect() {

	// add the proxy to mists list of subscribers
	p.Lock()
	subscribers[p.id] = p
	p.Unlock()

	// this gofunc handles matching messages to subscriptions for the proxy
	go p.handleMessages()
}

//
func (p *Proxy) handleMessages() {

	defer func() {
		close(p.check)
		close(p.Pipe)
	}()

	//
	for {
		select {

		// we need to ensure that this subscription actually has these tags before
		// sending anything to it; not doing this will cause everything to come
		// across the channel
		case msg := <-p.check:

			p.Lock()
			match := p.subscriptions.Match(msg.Tags)
			p.Unlock()

			// if there is a subscription for the tags publish the message
			if match {
				p.Pipe <- msg
			}

		//
		case <-p.done:
			return
		}
	}
}

// Ping
func (p *Proxy) Ping() error {
	return nil
}

// Subscribe
func (p *Proxy) Subscribe(tags []string) error {

	//
	// if !p.Authorized {
	// 	return ErrUnauthorized
	// }

	// is this an error?
	if len(tags) == 0 {
		return nil
	}

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
	// if !p.Authorized {
	// 	return ErrUnauthorized
	// }

	// is this an error?
	if len(tags) == 0 {
		return nil
	}

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
	// if !p.Authorized {
	// 	return ErrUnauthorized
	// }

	//
	return publish(p.id, tags, data)
}

// Sends a message with delay
func (p *Proxy) PublishAfter(tags []string, data string, delay time.Duration) error {

	//
	// if !p.Authorized {
	// 	return ErrUnauthorized
	// }

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
	// if !p.Authorized {
	// 	return ErrUnauthorized
	// }

	// convert the list into something friendlier
	p.Lock()
	var data []string
	for _, subscription := range p.subscriptions.ToSlice() {
		data = append(data, strings.Join(subscription, ","))
	}
	p.Unlock()

	msg := Message{
		Cmd:  "list",
		Tags: []string{},
		Data: fmt.Sprintf(strings.Join(data, " ")),
	}

	p.Pipe <- msg

	//
	return nil
}

//
func (p *Proxy) Close() {

	// this closes the goroutine that is matching messages to subscriptions
	close(p.done)

	// remove the local p from mists list of subscribers
	unsubscribe(p.id)
}
