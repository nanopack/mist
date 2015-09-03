// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"net"
	"net/http"
	"testing"
)

func TestMistCore(test *testing.T) {
	mist := New()
	client := NewLocalClient(mist, 0)
	defer client.Close()

	client.Subscribe([]string{"tag0"})
	for count := 0; count < 2; count++ {
		mist.Publish([]string{"tag0"}, "this is my data")
		message := <-client.Messages()
		assert(test, len(message.Tags) == 1, "wrong number of tags")
		// assert(test, message.Data == []byte("this is my data"), "data was incorrect")
	}

	client.Unsubscribe([]string{"tag0"})
	mist.Publish([]string{"tag0"}, "this is my data")
	select {
	case <-client.Messages():
		assert(test, false, "the message should not have been received")
	default:
	}

}

func BenchmarkMistCore(b *testing.B) {
	mist := New()
	client := NewLocalClient(mist, 0)
	defer client.Close()
	client.Subscribe([]string{"tag0"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mist.Publish([]string{"tag0"}, "this is my data")
		_ = <-client.Messages()
	}
}

func TestMistApi(test *testing.T) {
	mist := New()
	serverSocket, err := mist.Listen("127.0.0.1:1234")
	defer serverSocket.Close()
	assert(test, err == nil, "listen errored: %v", err)

	client, err := NewRemoteClient("127.0.0.1:1234")
	defer client.Close()
	assert(test, err == nil, "connect errored: %v", err)

	assert(test, client.Ping() == nil, "ping failed")

	client.Subscribe([]string{"tag"})
	client.Subscribe([]string{"other", "what", "is", "going", "on"})

	client.Publish([]string{"tag"}, "message")

	list, err := client.List()
	assert(test, err == nil, "listing subsctiptions failed %v", err)
	assert(test, len(list) == 2, "wrong number of subscriptions were returned %v", list)
	assert(test, len(list[0]) == 1, "wrong number of tags %v", list[0])
	assert(test, len(list[1]) == 5, "wrong number of tags %v", list[1])

	msg, ok := <-client.Messages()

	assert(test, ok, "got a nil message")
	assert(test, msg.Data == "message", "got the wrong message %v", msg.Data)

}

func TestMistWebsocket(test *testing.T) {
	mist := New()
	httpListener, err := net.Listen("tcp", "127.0.0.1:2345")
	assert(test, err == nil, "unable to listen to websockets %v", err)

	defer httpListener.Close()

	go http.Serve(httpListener, GenerateWebsocketUpgrade(mist))

	header := make(http.Header, 0)
	client, err := NewWebsocketClient("ws://127.0.0.1:2345/", header)
	assert(test, err == nil, "unable to connect %v", err)
	defer client.Close()
	err = client.Subscribe([]string{"test"})
	assert(test, err == nil, "subscription failed %v", err)
	mist.Publish([]string{"test"}, "some data")
	<-client.Messages()
	list, err := client.List()
	assert(test, err == nil, "unable to list %v", err)
	assert(test, len(list) == 1, "list of subscriptions is wrong %v", list)
	assert(test, len(list[0]) == 1, "wrong number of tags in subscription %v", list[0])
	err = client.Unsubscribe([]string{"test"})
	list, err = client.List()
	assert(test, err == nil, "unable to list %v", err)
	assert(test, len(list) == 0, "list of subscriptions is wrong %v", list)
	mist.Publish([]string{"test"}, "more data")
	client.Close()
}

func assert(test *testing.T, check bool, fmt string, args ...interface{}) {
	if !check {
		test.Logf(fmt, args...)
		test.FailNow()
	}
}
