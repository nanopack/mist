package mist

import (
	"fmt"
	"sync"

	"github.com/nanobox-core/hatchet"
)

//
const Version = "0.0.1"

//
type (

	//
	Mist struct {
		sync.Mutex

		log           hatchet.Logger
		Subscriptions map[string]map[chan Message]string // Subscriptions represent...
	}

	//
	Subscription struct {
		Tags []string     `json:"tags"`
		Sub  chan Message `json:"channel"`
	}

	//
	Message struct {
		Tags []string `json:"tags"`
		Data string   `json:"data"`
	}
)

//
func New(port string, logger hatchet.Logger) *Mist {
	fmt.Println("Initializing 'Mist'...")

  //
  if logger == nil {
    logger = hatchet.DevNullLogger{}
  }

	mist := &Mist{
    Subscriptions: make(map[string]map[chan Message]string),
    log: logger,
  }

	server := &Server{}
	server.start(port, mist)

	return mist
}

// Publish takes a list of tags and iterates through Mist's list of subscriptions,
// looking for matching subscriptions to publish messages too. It ensures that the
// list of recipients is a unique set, so as not to publish the same message more
// than once over a channel
func (m *Mist) Publish(tags []string, data string) {

	// a unique list of recipients (may contain duplicate channels from multiple
	// subscriptions)
	recipients := make(map[chan Message]int)

	// a *unique* list of recipients that will receive broadcasts
	// var recipients []chan Message

	// iterate through each provided tag looking for subscriptions to publish to
	for _, t := range tags {

		// keep track of how many times a subscription is requested
		used := 0

		// iterate through any matching subscriptions and add all of that subscriptions
		// channels to the list of recipients
		if sub, ok := m.Subscriptions[t]; ok {
			for ch, _ := range sub {

				// ensure that we keep the list of recipients unique, by checking each
				// match against a temporary map of found channels.
				if _, ok := recipients[ch]; !ok {
					used++

					// update our list of found channels, with a value of how many times
					// that channel has been subscribed to
					recipients[ch] = used

					// add the channel to our unique list of channels
					// recipients = append(recipients, ch)
				}
			}
		}
	}

	// format the data and send it on each unique recipient's channel
	msg := Message{Tags: tags, Data: data}

	//
	fmt.Printf("Publishing: %+v\n", msg)

	//
	for ch, _ := range recipients {
		go func() { ch <- msg }()
	}
}

// Subscribe
func (m *Mist) Subscribe(sub Subscription) {
	m.Lock()

	//
	fmt.Printf("Subscribing to: %+v\n", sub.Tags)

	// iterate over each subscription, adding it to our list of subscriptions (if
	// not already found), and then adding the channel into the subscription's list
	// of subscribers.
	for _, t := range sub.Tags {

		// if we don't find a subscription, make one (type []chan Message), and add
		// it to our list of subscriptions
		if _, ok := m.Subscriptions[t]; !ok {
			m.Subscriptions[t] = make(map[chan Message]string)
			fmt.Printf("Created new subscription '%+v'\n", t)
		}

		// add the channel to each subscription...
		m.Subscriptions[t][sub.Sub] = ""
		fmt.Printf("Subscribed '%+v' to '%+v'\n", sub.Sub, t)
	}
	m.Unlock()
}

// Unsubscribe
func (m *Mist) Unsubscribe(sub Subscription) {
	m.Lock()

	//
	fmt.Printf("Unsubscribing '%+v' from '%+v'\n", sub.Sub, sub.Tags)

	//
	for _, t := range sub.Tags {

		//
		if s, ok := m.Subscriptions[t]; ok {
			delete(s, sub.Sub)
			fmt.Printf("Unsubscribed '%+v' from '%+v'\n", sub.Sub, s)
		}

		//
		if len(m.Subscriptions[t]) <= 0 {
			delete(m.Subscriptions, t)
			fmt.Printf("Removed empty subscription '%+v'\n", t)
		}
	}
	m.Unlock()
}

//
func (m *Mist) List() error {
	fmt.Println(m.Subscriptions)
	return nil
}
