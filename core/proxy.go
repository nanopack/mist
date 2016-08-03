package mist

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/jcelliott/lumber"
)

type (
	// Proxy ...
	Proxy struct {
		sync.RWMutex

		Authenticated bool
		Pipe          chan Message
		check         chan Message
		done          chan bool
		id            uint32
		subscriptions subscriptions
	}
)

// NewProxy ...
func NewProxy() (p *Proxy) {

	// create new proxy
	p = &Proxy{
		Pipe:          make(chan Message),
		check:         make(chan Message),
		done:          make(chan bool),
		id:            atomic.AddUint32(&uid, 1),
		subscriptions: newNode(),
	}

	p.connect()

	return
}

// connect
func (p *Proxy) connect() {

	lumber.Trace("Proxy connecting...")
	// add the proxy to mists list of subscribers
	// p.Lock() // locked in subscribe() function
	subscribe(p)
	// p.Unlock()

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

			p.RLock()
			match := p.subscriptions.Match(msg.Tags)
			p.RUnlock()

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

// Subscribe ...
func (p *Proxy) Subscribe(tags []string) {
	lumber.Trace("Proxy subscribing to '%v'...", tags)

	// if len(tags) == 0 {
	// 	// is this an error?
	// }

	// add tags to subscription
	p.Lock()
	p.subscriptions.Add(tags)
	p.Unlock()
}

// Unsubscribe ...
func (p *Proxy) Unsubscribe(tags []string) {
	lumber.Trace("Proxy unsubscribing from '%v'...", tags)

	// if len(tags) == 0 {
	// 	// is this an error?
	// }

	// remove tags from subscription
	p.Lock()
	p.subscriptions.Remove(tags)
	p.Unlock()
}

// Publish ...
func (p *Proxy) Publish(tags []string, data string) error {
	lumber.Trace("Proxy publishing to %v...", tags)

	return publish(p.id, tags, data)
}

// PublishAfter sends a message after [delay]
func (p *Proxy) PublishAfter(tags []string, data string, delay time.Duration) {
	go func() {
		<-time.After(delay)
		if err := publish(p.id, tags, data); err != nil {
			// log this error and continue
			lumber.Error("Proxy failed to PublishAfter - %v", err)
		}
	}()
}

// List returns a list of all current subscriptions
func (p *Proxy) List() (data [][]string) {
	lumber.Trace("Proxy listing subscriptions...")
	p.RLock()
	data = p.subscriptions.ToSlice()
	p.RUnlock()

	return
}

// Close ...
func (p *Proxy) Close() {
	lumber.Trace("Proxy closing...")

	// this closes the goroutine that is matching messages to subscriptions
	close(p.done)

	// remove the local p from mists list of subscribers
	unsubscribe(p.id)
}
