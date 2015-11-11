// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"errors"
	"github.com/nanopack/mist/subscription"
	"sort"
	"sync"
)

type (
	Subscriptions interface {
		Add([]string)
		Remove([]string)
		Match([]string) bool
		ToSlice() [][]string
	}

	EnableReplication interface {
		EnableReplication() error
	}

	EnableInternal interface {
		EnableInternal()
	}

	localSubscriber struct {
		sync.Mutex

		check chan Message
		done  chan bool
		pipe  chan Message

		subscriptions Subscriptions
		mist          *Mist
		id            uint32
		internal      bool
		replicated    bool
	}
)

var (
	InternalErr = errors.New("Unable to perform action, internal mode enabled")
)

//
func NewLocalClient(mist *Mist, buffer int) Client {
	client := &localSubscriber{
		check:         make(chan Message, buffer),
		done:          make(chan bool),
		pipe:          make(chan Message),
		mist:          mist,
		id:            mist.nextId(),
		subscriptions: subscription.NewNode(),
		internal:      false,
	}

	// this gofunc handles matching messages to subscriptions for the client
	go func(client *localSubscriber) {

		defer func() {
			close(client.check)
			close(client.pipe)
		}()

		for {
			select {
			case msg := <-client.check:

				switch {
				case msg.internal && client.internal:
					client.pipe <- msg
				default:
					client.Lock()
					match := client.subscriptions.Match(msg.Tags)
					client.Unlock()

					if match {
						client.pipe <- msg
					}
				}
			case <-client.done:
				return
			}
		}
	}(client)

	// add the local client to mists list of subscribers
	mist.subscribers[client.id] = client

	return client
}

func (client *localSubscriber) EnableInternal() {
	// we don't want any already replicated messages to come across on this client
	// this will stop that
	delete(client.mist.subscribers, client.id)
	client.internal = true
	client.mist.internal[client.id] = client
}

func (client *localSubscriber) EnableReplication() error {
	// we need to flag that this client doesn't use publish any more.
	client.replicated = true

	// this client is no longer a subscriber.
	delete(client.mist.subscribers, client.id)
	client.mist.replicators[client.id] = client
	return nil
}

//
func (client *localSubscriber) List() ([][]string, error) {
	if client.internal {
		return nil, InternalErr
	}
	return client.subscriptions.ToSlice(), nil
}

//
func (client *localSubscriber) Subscribe(tags []string) error {
	if client.internal {
		return InternalErr
	}
	if len(tags) == 0 {
		return nil
	}
	sort.Sort(sort.StringSlice(tags))
	client.Lock()
	client.subscriptions.Add(tags)
	client.Unlock()

	// if this client is replicated, we don't need to inform other replicated
	// clients of this subscription
	if !client.replicated {
		// notify anyone who is interested about the new subscription
		client.mist.publish(tags, "subscribe")
	}

	return nil
}

// Unsubscribe iterates through each of mist clients subscriptions keeping all subscriptions
// that aren't the specified subscription
func (client *localSubscriber) Unsubscribe(tags []string) error {
	if client.internal {
		return InternalErr
	}
	if len(tags) == 0 {
		return nil
	}
	sort.Sort(sort.StringSlice(tags))
	client.Lock()
	client.subscriptions.Remove(tags)
	client.Unlock()

	// if this client is replicated, we don't need to inform other replicated
	// clients of this unsubscription
	if !client.replicated {
		client.mist.publish(tags, "unsubscribe")
	}
	return nil
}

// Sends a message across mist
func (client *localSubscriber) Publish(tags []string, data string) error {
	if client.internal {
		return InternalErr
	}
	switch client.replicated {
	case true:
		client.mist.Replicate(tags, data)
	default:
		client.mist.Publish(tags, data)
	}
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

	// remove the local client from mists list of subscribers/replicators/internal
	delete(client.mist.subscribers, client.id)
	delete(client.mist.replicators, client.id)
	delete(client.mist.internal, client.id)

	// send out the unsubscribe message to anyone listening
	for _, subscription := range client.subscriptions.ToSlice() {
		client.mist.publish(subscription, "unsubscribe")
	}

	return nil
}
