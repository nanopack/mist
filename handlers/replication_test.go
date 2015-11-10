// -*- mode: go; tab-width: 2; indent-tabs-mode: 1; st-rulers: [70] -*-
// vim: ts=4 sw=4 ft=lua noet
//--------------------------------------------------------------------
// @author Daniel Barney <daniel@nanobox.io>
// Copyright (C) Pagoda Box, Inc - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly
// prohibited. Proprietary and confidential
//
// @doc
//
// @end
// Created :   12 August 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
package handlers_test

import (
	"github.com/golang/mock/gomock"
	"github.com/nanobox-io/golang-mist/core"
	"github.com/nanobox-io/golang-mist/handlers"
	"github.com/nanobox-io/golang-mist/handlers/mock"
	"testing"
	"time"
)

func TestReplication(test *testing.T) {
	ctrl := gomock.NewController(test)
	defer ctrl.Finish()

	discover1 := mock_golang_discovery.NewMockDiscover(ctrl)
	discover2 := mock_golang_discovery.NewMockDiscover(ctrl)

	// we have 2 mists that we are forwarding between
	mist1 := mist.New()
	mist2 := mist.New()

	listen1, err := mist1.Listen("127.0.0.1:2223", nil)
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
	defer listen1.Close()
	listen2, err := mist2.Listen("127.0.0.1:2224", nil)
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
	defer listen2.Close()

	discover1.EXPECT().Handle("mist", gomock.Any())
	discover2.EXPECT().Handle("mist", gomock.Any())

	rep1 := handlers.EnableReplication(mist1, discover1)
	rep2 := handlers.EnableReplication(mist2, discover2)

	// we start up the replication monitor
	go rep1.Monitor()
	go rep2.Monitor()

	// inform both replication monitors of the other nodes
	repc1 := rep2.New("127.0.0.1:2223")
	defer repc1.Close()
	repc2 := rep1.New("127.0.0.1:2224")
	defer repc2.Close()

	// now we test the clients out.
	client1 := mist.NewLocalClient(mist1, 0)
	client2 := mist.NewLocalClient(mist2, 0)

	// subscribe and wait, because it could take a little bit of time
	// for everything to propagate correctly
	client2.Subscribe([]string{"tag"})
	<-time.After(time.Millisecond * 10)

	// send a message and ensure that it is sent across correctly
	client1.Publish([]string{"tag"}, "what")
	select {
	case msg := <-client2.Messages():
		if msg.Data != "what" {
			test.Log("got the wrong message", msg)
			test.FailNow()
		}
	case <-time.After(time.Second):
		test.Log("the message was not replicated correctly")
		test.FailNow()
	}

	client2.Unsubscribe([]string{"tag"})
	<-time.After(time.Millisecond * 10)

	// send a message and ensure that it is not sent across
	client1.Publish([]string{"tag"}, "what")
	select {
	case msg := <-client2.Messages():
		test.Log("got the wrong message", msg)
		test.FailNow()
	case <-time.After(time.Second):
	}
}
