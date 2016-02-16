package handlers_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/handlers"
	"github.com/nanopack/mist/handlers/mock"
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
	repc2 := rep2.New("127.0.0.1:2223")
	defer repc2.Close()
	repc1 := rep1.New("127.0.0.1:2224")
	defer repc1.Close()

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
