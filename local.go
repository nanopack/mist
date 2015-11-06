// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	set "github.com/deckarep/golang-set"
	"strings"
	"sync"
)

type (
	tagSlice     []string
	subscription struct {
		present set.Set
		absent  set.Set
	}

	localSubscriber struct {
		sync.Mutex

		check chan Message
		done  chan bool
		pipe  chan Message

		subscriptions []subscription
		mist          *Mist
		id            uint32
	}

	Replicator interface {
		EnableReplication() Replicator
		Replicate([]string, string) error
	}
)

//
func NewLocalClient(mist *Mist, buffer int) Client {
	client := &localSubscriber{
		check: make(chan Message, buffer),
		done:  make(chan bool),
		pipe:  make(chan Message),
		mist:  mist,
		id:    mist.nextId()}

	// this gofunc handles matching messages to subscriptions for the client
	go func(client *localSubscriber) {

		defer func() {
			close(client.check)
			close(client.pipe)
		}()

		for {
			select {
			case msg := <-client.check:
				// we do this so that we don't need a mutex
				subscriptions := client.subscriptions
				for _, subscription := range subscriptions {
					if subscription.Check(msg.tags) {
						client.pipe <- msg
					}
				}
			case <-client.done:
				return
			}
		}
	}(client)

	// add the local client to mists list of subscribers
	mist.subscribers[client.id] = *client
	mist.replicators[client.id] = *client

	return client
}

func (client *localSubscriber) EnableReplication() Replicator {
	// we don't want any already replicated messages to come across on this client
	// this will stop that
	mist.replicators[client.id] = nil
	return client
}

func (client *localSubscriber) Replicate(tags []string, data string) error {
	return client.mist.Replicate(tags, data)
}

//
func (client *localSubscriber) List() ([][]string, error) {
	subscriptions := make([][]string, len(client.subscriptions))
	for i, subscription := range client.subscriptions {
		subscriptions[i] = tagSlice(subscription.ToSlice()).Clean()
	}
	return subscriptions, nil
}

//
func (client *localSubscriber) Subscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	subscription := makeSubscription(append(tags, "-__mist*internal__")...)

	client.Lock()
	client.subscriptions = append(client.subscriptions, subscription)
	client.Unlock()

	client.mist.Publish(append(tags, "__mist*internal__"), "subscribe")

	return nil
}

func (client *localSubscriber) TailSubscriptions() error {
	sub := makeSubscription("__mist*internal__")
	client.Lock()
	client.subscriptions = []subscription{sub}
	client.Unlock()

	return nil
}

// Unsubscribe iterates through each of mist clients subscriptions keeping all subscriptions
// that aren't the specified subscription
func (client *localSubscriber) Unsubscribe(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	//create a set for quick comparison
	test := makeSubscription(append(tags, "-__mist*internal__")...)

	// create a slice of subscriptions that are going to be kept
	keep := []subscription{}

	client.Lock()

	// iterate over all of mist clients subscriptions looking for ones that match the
	// subscription to unsubscribe
	for _, subscription := range client.subscriptions {

		// if they are not the same set (meaning they are a different subscription) then add them
		// to the keep set
		if !test.Equal(subscription) {
			keep = append(keep, subscription)
		}
	}

	client.subscriptions = keep

	client.Unlock()

	client.mist.Publish(append(tags, "__mist*internal__"), "unsubscribe")
	return nil
}

// Sends a message across mist
func (client *localSubscriber) Publish(tags []string, data string) error {
	// remove all internal tags from the tag list
	tags = tagSlice(tags).Clean()

	client.mist.Publish(tags, data)
	return nil
}

//
func (client *localSubscriber) Ping() error {
	return nil
}

// Returns all messages that have sucessfully matched the list of subscriptions that this
// client has subscribed to
func (client *localSubscriber) Messages() <-chan Message {
	return client.pipe
}

//
func (client *localSubscriber) Close() error {
	// this closes the goroutine that is matching messages to subscriptions
	close(client.done)

	// remove the local client from mists list of subscribers
	delete(client.mist.subscribers, client.id)

	// send out the unsubscribe message to anyone listening
	for _, subscription := range client.subscriptions {
		client.Publish(append(subscription.ToSlice(), "__mist*internal__"), "unsubscribe")
	}

	return nil
}

//
func makeSubscription(tags ...string) subscription {
	present := set.NewThreadUnsafeSet()
	absent := set.NewThreadUnsafeSet()
	for _, i := range tags {
		switch {
		case strings.HasPrefix(i, "-"):
			absent.Add(i[1:])
		default:
			present.Add(i)
		}
	}

	return subscription{
		absent:  absent,
		present: present,
	}
}

func makeBareSet(tags []string) set.Set {
	tagSlice := set.NewThreadUnsafeSet()
	for _, i := range tags {
		tagSlice.Add(i)
	}

	return tagSlice
}

func (tags tagSlice) Clean() tagSlice {
	for i := len(tags) - 1; i >= 0; i-- {
		switch {
		case tags[i] == "__mist*internal__":
			fallthrough
		case tags[i] == "-__mist*internal__":
			tags = append(tags[:i], tags[i+1:]...)
		}
	}
	return tags
}

func (sub subscription) Check(tags set.Set) bool {
	return sub.present.IsSubset(tags) && sub.absent.Intersect(tags).Cardinality() == 0
}

func (sub subscription) Equal(test subscription) bool {
	return sub.absent.Equal(test.absent) && sub.present.Equal(test.present)
}

func (sub subscription) ToSlice() []string {
	slice := make([]string, sub.absent.Cardinality()+sub.present.Cardinality())
	for j, tag := range sub.absent.ToSlice() {
		slice[j] = "-" + tag.(string)
	}
	i := sub.absent.Cardinality()
	for j, tag := range sub.present.ToSlice() {
		slice[i+j] = tag.(string)
	}

	return slice
}
