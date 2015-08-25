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

type (
	Client interface {
		List() ([][]string, error)
		Subscribe(tags []string) error
		Unsubscribe(tags []string) error
		Publish(tags []string, data string) error
		Ping() error
		Close() error
		Messages() <-chan Message
	}

	//
	Mist struct {
		subscribers map[uint32]subscriber
		next        uint32
	}

	//
	subscriber struct {
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
		Tags []string    `json:"tags"`
		Data interface{} `json:"data"`
	}
)

// creates a new mist
func New() *Mist {

	return &Mist{
		subscribers: make(map[uint32]subscriber),
	}
}
func (mist *Mist) nextId() uint32 {
	return atomic.AddUint32(&mist.next, 1)
}

func (mist *Mist) addSubscriber(subscriber *subscriber) {
	mist.subscribers[subscriber.id] = *subscriber
}

func (mist *Mist) removeSubscriber(id uint32) {
	// remove this subscriber from mist
	delete(mist.subscribers, id)
}

// Publish takes a list of tags and iterates through mist's list of subscribers,
// sending to each if they are available.
func (mist *Mist) Publish(tags []string, data interface{}) {

	message := Message{
		Tags: tags,
		tags: makeSet(tags),
		Data: data}

	for _, subscriber := range mist.subscribers {
		select {
		case <-subscriber.done:
		case subscriber.check <- message:
			// default:
			// do we drop the message? enqueue it? pull one off the front and then add this one?
		}
	}
}
