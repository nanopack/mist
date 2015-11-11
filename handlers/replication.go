// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package handlers

import (
	"github.com/nanobox-io/golang-discovery"
	"github.com/nanobox-io/golang-mist/core"
	"github.com/nanobox-io/golang-mist/subscription"
	"io"
	"time"
)

type (
	Looper interface {
		Loop(time.Duration) error
	}

	replicate struct {
		mist          *mist.Mist
		clients       []mist.Client
		newClient     chan mist.Client
		doneClient    chan mist.Client
		subscriptions mist.Subscriptions
	}
)

func EnableReplication(server *mist.Mist, discover discovery.Discover) *replicate {

	replicate := &replicate{
		mist:          server,
		newClient:     make(chan mist.Client),
		doneClient:    make(chan mist.Client),
		subscriptions: subscription.NewNode(),
	}

	discover.Handle("mist", replicate)

	return replicate
}

func (rep *replicate) Monitor() {
	// we want to catch all subscription/unsubscription changes.
	// this should at least give a good safety zone.
	client := mist.NewLocalClient(rep.mist, 100)
	defer client.Close()

	// set the client to be in internal mode
	// now only internal message will be received
	client.(mist.EnableInternal).EnableInternal()

	for {
		select {
		case msg, ok := <-client.Messages():
			if !ok {
				return
			}
			switch msg.Data {
			case "subscribe":
				rep.subscriptions.Add(msg.Tags)
			case "unsubscribe":
				rep.subscriptions.Remove(msg.Tags)
			default:
				// we ignore all other messages
				continue
			}
			rep.forwardAll(msg.Data, msg.Tags)
		case remote := <-rep.doneClient:
			// very innefficeint, really shouldn't be a slice
			// shouldn't matter unless we have a cluster of over 100
			// machines
			for i := len(rep.clients) - 1; i >= 0; i-- {
				if rep.clients[i] == remote {
					rep.clients = append(rep.clients[:i], rep.clients[i+1:]...)
				}
			}
		case remote := <-rep.newClient:
			rep.clients = append(rep.clients, remote)

			// forward all published messages to the local mist server
			go func() {
				for msg := range remote.Messages() {
					rep.mist.Replicate(msg.Tags, msg.Data)
				}
				rep.doneClient <- remote
			}()

			// send all subscriptions across the connection.
			for _, subscription := range rep.subscriptions.ToSlice() {
				forward(remote, "subscribe", subscription)
			}
		}
	}
}

func (rep replicate) forwardAll(fun string, subscription []string) {
	perform := getFunc(fun)
	for _, client := range rep.clients {
		if err := perform(client, subscription); err != nil {
			// should we log this error?
		}
	}
}

func forward(client mist.Client, fun string, subscription []string) {
	perform := getFunc(fun)
	if err := perform(client, subscription); err != nil {
		// should we log this error?
	}
}

func getFunc(fun string) func(mist.Client, []string) error {
	if fun == "subscribe" {
		return mist.Client.Subscribe
	}
	return mist.Client.Unsubscribe
}

func (rep *replicate) New(address string) io.Closer {
	client, err := mist.NewRemoteClient(address)
	if err != nil {
		return nil
	}

	if client.(mist.EnableReplication).EnableReplication() != nil {
		return nil
	}
	// add this client to the list of all clients
	rep.newClient <- client

	return client
}
