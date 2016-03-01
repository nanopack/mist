package handlers

import (
// "testing"
// "time"
//
// "github.com/golang/mock/gomock"
// "github.com/nanopack/mist/core"
// "github.com/nanopack/mist/handlers/mock"
)

var (
	testTag = "hello"
	testMsg = "world"
)

//
// func TestReplication(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	d1 := mock_golang_discovery.NewMockDiscover(ctrl)
// 	d2 := mock_golang_discovery.NewMockDiscover(ctrl)
//
// 	// we have 2 mists that we are forwarding between
// 	mist1 := mist.New()
// 	mist2 := mist.New()
//
// 	//
// 	l1, err := mist1.Listen("127.0.0.1:2223", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer l1.Close()
//
// 	//
// 	l2, err := mist2.Listen("127.0.0.1:2224", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer l2.Close()
//
// 	d1.EXPECT().Handle("mist", gomock.Any())
// 	d2.EXPECT().Handle("mist", gomock.Any())
//
// 	rep1 := EnableReplication(mist1, d1)
// 	rep2 := EnableReplication(mist2, d2)
//
// 	// we start up the replication monitor
// 	go rep1.Monitor()
// 	go rep2.Monitor()
//
// 	// inform both replication monitors of the other nodes
// 	repc2 := rep2.New("127.0.0.1:2223")
// 	defer repc2.Close()
// 	repc1 := rep1.New("127.0.0.1:2224")
// 	defer repc1.Close()
//
// 	// now we test the clients out.
// 	client1, _ := mist.NewProxy(mist1, 0)
// 	client2, _ := mist.NewProxy(mist2, 0)
//
// 	// subscribe and wait, because it could take a little bit of time
// 	// for everything to propagate correctly
// 	client2.Subscribe([]string{testTag})
// 	<-time.After(time.Millisecond * 10)
//
// 	// send a message and ensure that it is sent across correctly
// 	client1.Publish([]string{testTag}, testMsg)
// 	select {
// 	case msg := <-client2.Messages():
// 		if msg.Data != testMsg {
// 			t.Errorf("got the wrong message %v", msg)
// 		}
// 	case <-time.After(time.Second):
// 		t.Error("the message was not replicated correctly")
// 	}
//
// 	client2.Unsubscribe([]string{testTag})
// 	<-time.After(time.Millisecond * 10)
//
// 	// send a message and ensure that it is not sent across
// 	client1.Publish([]string{testTag}, testMsg)
// 	select {
// 	case msg := <-client2.Messages():
// 		t.Errorf("got the wrong message %v", msg)
// 	case <-time.After(time.Second):
// 	}
// }
