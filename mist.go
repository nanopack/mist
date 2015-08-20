// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	set "github.com/deckarep/golang-set"
	"sync"
	"sync/atomic"
)

//
type (

	//
	Mist struct {
		clients map[uint32]MistClient
		next    uint32
	}

	//
	MistClient struct {
		sync.Mutex

		check chan Message
		done  chan bool
		pipe  chan Message

		subscriptions []set.Set
		mist          *Mist
		id            uint32
	}

	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		tags set.Set
		Tags []string
		Data interface{}
	}
)

// creates a new mist
func New() *Mist {

	return &Mist{
		clients: make(map[uint32]MistClient),
	}
}

func makeSet(tags []string) set.Set {
	set := set.NewThreadUnsafeSet()
	for _, i := range tags {
		set.Add(i)
	}

	return set
}

// Publish takes a list of tags and iterates through mist's list of clients,
// sending to each if they are available.
func (mist *Mist) Publish(tags []string, data interface{}) {

	message := Message{
		Tags: tags,
		tags: makeSet(tags),
		Data: data}

	for _, client := range mist.clients {
		select {
		case <-client.done:
		case client.check <- message:
			// default:
			// 	// do we drop the message? enqueue it? pull one off the front and then add this one?
		}
	}
}

func (mist *Mist) Client(buffer int) *MistClient {
	client := MistClient{
		check: make(chan Message, buffer),
		done:  make(chan bool),
		pipe:  make(chan Message),
		mist:  mist,
		id:    atomic.AddUint32(&mist.next, 1)}

	// this gofunc handles matching messages to subscriptions for the client
	go func(client *MistClient) {

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
					if subscription.IsSubset(msg.tags) {
						client.pipe <- msg
					}
				}
			case <-client.done:
				return
			}
		}
	}(&client)

	mist.clients[client.id] = client

	return &client
}

func (client *MistClient) Subscribe(tags []string) {
	subscription := makeSet(tags)

	client.Lock()
	client.subscriptions = append(client.subscriptions, subscription)
	client.Unlock()
}

// Unsubscribe iterates through each of mist clients subscriptions keeping all subscriptions
// that aren't the specified subscription
func (client *MistClient) Unsubscribe(tags []string) {
	client.Lock()

	//create a set for quick comparison
	test := makeSet(tags)

	// create a slice of subscriptions that are going to be kept
	keep := []set.Set{}

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
}

func (client *MistClient) List() [][]string {
	subscriptions := make([][]string, len(client.subscriptions))
	for i, subscription := range client.subscriptions {
		sub := make([]string, subscription.Cardinality())
		for j, tag := range subscription.ToSlice() {
			sub[j] = tag.(string)
		}
		subscriptions[i] = sub
	}
	return subscriptions
}

func (client *MistClient) Close() {
	// this closes the goproc that is matching messages to subscriptions
	close(client.done)

	// remove this client from mist
	delete(client.mist.clients, client.id)
}

// Returns all messages that have sucessfully matched the list of subscriptions that this
// client has subscribed to
func (client *MistClient) Messages() <-chan Message {
	return client.pipe
}

// Sends a message across mist
func (client *MistClient) Publish(tags []string, data interface{}) {
	client.mist.Publish(tags, data)
}
