// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"sort"
	"sync/atomic"
)

type (

	//
	Client interface {
		List() ([][]string, error)
		Subscribe(tags []string) error
		Unsubscribe(tags []string) error
		Publish(tags []string, data string) error
		Ping() error
		Messages() <-chan Message
		Close() error
	}

	//
	Mist struct {
		subscribers map[uint32]*localSubscriber
		replicators map[uint32]*localSubscriber
		internal    map[uint32]*localSubscriber
		next        uint32
	}

	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		internal bool
		Tags     []string `json:"tags"`
		Data     string   `json:"data"`
	}
)

// creates a new mist
func New() *Mist {

	return &Mist{
		subscribers: make(map[uint32]*localSubscriber),
		replicators: make(map[uint32]*localSubscriber),
		internal:    make(map[uint32]*localSubscriber),
	}
}

// Publish takes a list of tags and iterates through mist's list of subscribers,
// sending to each if they are available.
func (mist *Mist) Publish(tags []string, data string) error {

	// is this an error? or just something we need to ignore
	if len(tags) == 0 {
		return nil
	}

	message := Message{
		Tags: tags,
		Data: data,
	}
	// publishes go to both subscribers, and to replicators
	forward(message, mist.subscribers)
	forward(message, mist.replicators)

	return nil
}

func (mist *Mist) Replicate(tags []string, data string) error {

	// is this an error? or just something we need to ignore
	if len(tags) == 0 {
		return nil
	}

	message := Message{
		Tags: tags,
		Data: data,
	}
	// replicate only goes to subscribers
	forward(message, mist.subscribers)

	return nil
}

func (mist *Mist) publish(tags []string, data string) error {
	// is this an error? or just something we need to ignore
	if len(tags) == 0 {
		return nil
	}

	message := Message{
		Tags:     tags,
		internal: true,
		Data:     data,
	}
	forward(message, mist.internal)

	return nil
}

func forward(msg Message, subscribers map[uint32]*localSubscriber) {
	// we do this here so that the tags come pre sorted for the clients
	sort.Sort(sort.StringSlice(msg.Tags))
	// this should be more optimized, but it might not be an issue unless thousands of clients
	// are using mist.
	for _, localReplicator := range subscribers {
		select {
		case <-localReplicator.done:
		case localReplicator.check <- msg:
			// default:
			// do we drop the message? enqueue it? pull one off the front and then add this one?
		}
	}
}

//
func (mist *Mist) nextId() uint32 {
	return atomic.AddUint32(&mist.next, 1)
}
