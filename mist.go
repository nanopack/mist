// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"fmt"
	"sync"

	"github.com/pagodabox/golang-hatchet"
)

//
const (
	DefaultPort = "1445"
	Version     = "0.0.1"
)

//
type (

	//
	Mist struct {
		sync.Mutex

		log           hatchet.Logger //
		port          string
		Subscriptions []Subscription //
	}

	//
	Subscription struct {
		Tags []string     `json:"tags"`
		Sub  chan Message `json:"channel"`
	}

	//
	Message struct {
		Tags []string `json:"tags"`
		Data string   `json:"data,string"`
	}
)

//
func New(port string, logger hatchet.Logger) *Mist {

	//
	if logger == nil {
		logger = hatchet.DevNullLogger{}
	}

	mist := &Mist{
		log:           logger,
		port:          port,
		Subscriptions: []Subscription{},
	}

	mist.start()

	return mist
}

// Publish takes a list of tags and iterates through Mist's list of subscriptions,
// looking for matching subscriptions to publish messages too. It ensures that the
// list of recipients is a unique set, so as not to publish the same message more
// than once over a channel
func (m *Mist) Publish(tags []string, data string) {

	m.log.Info("PUBLISH!! %v, %v", tags, data)

	// create a message
	msg := Message{Tags: tags, Data: data}

	m.log.Info("MSG: %#v\n", msg)

	for _, subscription := range m.Subscriptions {
		if contains(tags, subscription.Tags) {
			go func(ch chan Message, msg Message) { ch <- msg }(subscription.Sub, msg)
		}
	}
}

// Subscribe
func (m *Mist) Subscribe(sub Subscription) {
	m.Lock()
	m.Subscriptions = append(m.Subscriptions, sub)
	m.Unlock()
}

// Unsubscribe
func (m *Mist) Unsubscribe(sub Subscription) {
	m.Lock()

	newSubscriptions := []Subscription{}
	for _, subscription := range m.Subscriptions {
		if !sameSub(subscription, sub) {
			newSubscriptions = append(newSubscriptions, sub)
		}
	}
	m.Subscriptions = newSubscriptions

	m.Unlock()
}

// List
func (m *Mist) List() {
	fmt.Println(m.Subscriptions)
}

func contains(full, subset []string) bool {
	fullMap := map[string]interface{}{}
	for _, str := range full {
		fullMap[str] = ""
	}
	for _, str := range subset {
		if fullMap[str] == nil {
			return false
		}
	}

	return true
}

func sameSub(x, y Subscription) bool {
	if len(x.Tags) != len(y.Tags) || x.Sub != y.Sub {
		return false
	}
	for i := 0; i < len(x.Tags); i++ {
		if x.Tags[i] != y.Tags[i] {
			return false
		}
	}

	return true
}