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
	proxy struct {
		sync.Mutex

		check chan Message
		done  chan bool
		pipe  chan Message

		subscriptions subscription.Subscriptions
		mist          *Mist
		id            uint32
		internal      bool
		replicated    bool
	}
)

//
func NewProxy(buffer int) (Client, error) {

	//
	p := &proxy{
		subscriptions: subscription.NewNode(),
		done:          make(chan bool),
		check:         make(chan Message, buffer),
		pipe:          make(chan Message),
		mist:          Self,
		id:            atomic.AddUint32(&Self.next, 1),
		internal:      false,
	}

	return p, p.connect()
}

// connect
func (p *proxy) connect() error {

	fmt.Println("PROXY CONNECT!", p.id)

	// this gofunc handles matching messages to subscriptions for the p
	go func() {

		defer func() {
			close(p.check)
			close(p.pipe)
		}()

		//
		for {
			select {

			//
			case msg := <-p.check:

				fmt.Println("PROXY MSG!", msg)

				switch {

				//
				case msg.internal && p.internal:
					fmt.Println("PROXY INTERNAL?!")
					p.pipe <- msg

				//
				default:
					fmt.Println("PROXY DEFAULT!")
					p.Lock()
					match := p.subscriptions.Match(msg.Tags)
					p.Unlock()

					if match {
						p.pipe <- msg
					}
				}

				//
			case <-p.done:
				fmt.Println("PROXY DONE!")
				return
			}
		}
	}()

	// add the local p to mists list of subscribers
	p.mist.subscribers[p.id] = p

	return nil
}

//
func (p *proxy) Ping() error {
	fmt.Println("PROXY PING!")
	return nil
}

// Subscribe
func (p *proxy) Subscribe(tags []string) error {

	fmt.Println("PROXY SUBSCRIBE!")

	//
	if p.internal {
		return InternalErr
	}

	//
	if len(tags) == 0 {
		return nil
	}

	sort.Sort(sort.StringSlice(tags))

	//
	p.Lock()
	p.subscriptions.Add(tags)
	p.Unlock()

	// if this p is replicated, we don't need to inform other replicated
	// ps of this subscription
	if !p.replicated {
		p.mist.publish(tags, "subscribe")
	}

	return nil
}

// Unsubscribe
func (p *proxy) Unsubscribe(tags []string) error {

	fmt.Println("PROXY UNSUBSCRIBE!")

	//
	if p.internal {
		return InternalErr
	}

	//
	if len(tags) == 0 {
		return nil
	}

	sort.Sort(sort.StringSlice(tags))

	//
	p.Lock()
	p.subscriptions.Remove(tags)
	p.Unlock()

	// if this p is replicated, we don't need to inform other replicated
	// ps of this unsubscription
	if !p.replicated {
		p.mist.publish(tags, "unsubscribe")
	}

	return nil
}

// Publish
func (p *proxy) Publish(tags []string, data string) error {

	fmt.Println("PROXY PUBLISH!")

	//
	if p.internal {
		return InternalErr
	}

	//
	switch p.replicated {
	case true:
		p.mist.Replicate(tags, data)
	default:
		p.mist.Publish(tags, data)
	}

	return nil
}

// Sends a message with delay
func (p *proxy) PublishAfter(tags []string, data string, delay time.Duration) error {

	//
	if p.internal {
		return InternalErr
	}

	//
	go func() {
		<-time.After(delay)
		p.Publish(tags, data)
	}()

	//
	return nil
}

// List
func (p *proxy) List() ([][]string, error) {

	fmt.Println("PROXY LIST!")

	//
	if p.internal {
		return nil, InternalErr
	}

	//
	return p.subscriptions.ToSlice(), nil
}

//
func (p *proxy) Close() error {

	fmt.Println("PROXY CLOSE!")

	// this closes the goroutine that is matching messages to subscriptions
	close(p.done)

	// remove the local p from mists list of subscribers/replicators/internal
	delete(p.mist.subscribers, p.id)
	delete(p.mist.replicators, p.id)
	delete(p.mist.internal, p.id)

	// send out the unsubscribe message to anyone listening
	for _, subscription := range p.subscriptions.ToSlice() {
		p.mist.publish(subscription, "unsubscribe")
	}

	return nil
}

// Returns all messages that have sucessfully matched the list of subscriptions
// that this p has subscribed to
func (p *proxy) Messages() <-chan Message {
	fmt.Println("PROXY MESSAGES!")
	return p.pipe
}

//
func (p *proxy) EnableInternal() {
	// we don't want any already replicated messages to come across on this p
	// this will stop that
	delete(p.mist.subscribers, p.id)
	p.internal = true
	p.mist.internal[p.id] = p
}

//
func (p *proxy) EnableReplication() error {
	// we need to flag that this p doesn't use publish any more.
	p.replicated = true

	// this p is no longer a subscriber.
	delete(p.mist.subscribers, p.id)
	p.mist.replicators[p.id] = p
	return nil
}
