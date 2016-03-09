package clients

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/server"
)

var (
	testTag = "hello"
	testMsg = "world"
)

// TestMain
func TestMain(m *testing.M) {

	//
	server.StartTCP(mist.DEFAULT_ADDR, nil)

	//
	os.Exit(m.Run())
}

// TestTCPClient tests
func TestTCPClientSubscriptions(t *testing.T) {

	//
	client, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer client.Close()

	// test ping
	if err := client.Ping(); err != nil {
		t.Fatalf("ping failed")
	}

	// test subscribe
	err = client.Subscribe([]string{testTag})
	err = client.Subscribe([]string{testTag, testMsg})
	if err != nil {
		t.Fatalf("client subscriptions failed %v", err.Error())
	}

	// test list
	err = client.List()
	if err != nil {
		t.Fatalf("listing subscriptions failed %v", err.Error())
	}

	//
	select {

	//
	case msg := <-client.Messages():

		fmt.Println("MESG!", msg)

		//
		subs := strings.Split(msg.Data, " ")
		if len(subs) != 2 {
			t.Fatalf("Incorrect number of subscriptions returned. Expecting 2 got %v", len(subs))
		}
		break

	//
	case <-time.After(time.Second * 1):
		t.Fatalf("Expecting messages, received none!")
	}

	// test unsubscribe
	err = client.Unsubscribe([]string{testTag})
	err = client.Unsubscribe([]string{testTag, testMsg})
	if err != nil {
		t.Fatalf("client unsubscriptions failed %v", err.Error())
	}
}

// TestSameTCPClient tests to ensure that mist will not send message to the
// same client who publishes them
func TestSameTCPClient(t *testing.T) {

	//
	sender, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer sender.Close()

	// sender subscribes to tags and then tries to publish to those same tags...
	err = sender.Subscribe([]string{testTag})
	if err != nil {
		t.Fatalf("client subscription failed %v", err.Error())
	}
	defer sender.Unsubscribe([]string{testTag})

	err = sender.Publish([]string{testTag}, testMsg)
	if err != nil {
		t.Fatalf("client publish failed %v", err.Error())
	}

	// sender should NOT get a message because mist shouldnt send a message to the
	// same proxy that publishes them.
	select {

	// wait for a message...
	case <-sender.Messages():
		t.Fatalf("Received own message!")

	// after 1 second assume no message is coming
	case <-time.After(time.Second * 1):
		break
	}
}

// TestDifferentTCPClient tests to ensure that mist will send messages
// to another subscribed client, and then not send when unsubscribed.
func TestDifferentTCPClient(t *testing.T) {

	//
	sender, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer sender.Close()

	//
	receiver, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer receiver.Close()

	// receiver subscribes to tags and then sender publishes to those tags...
	err = receiver.Subscribe([]string{testTag})
	if err != nil {
		t.Fatalf("client subscription failed %v", err.Error())
	}

	//
	err = sender.Publish([]string{testTag}, testMsg)
	if err != nil {
		t.Fatalf("client publish failed %v", err.Error())
	}

	//
	waitMessage(receiver.Messages(), t)

	// receiver unsubscribes from the tags and sender publishes again to the same
	// tags
	err = receiver.Unsubscribe([]string{testTag})
	if err != nil {
		t.Fatalf("client unsubscribe failed %v", err.Error())
	}

	// allow time for communication across TCP connection
	<-time.After(time.Second * 1)

	//
	err = sender.Publish([]string{testTag}, testMsg)
	if err != nil {
		t.Fatalf("client publish failed %v", err.Error())
	}

	// receiver should NOT get a message this time
	select {

	// wait for a message...
	case msg := <-receiver.Messages():
		t.Fatalf("Received a message from unsubscribed tags: %#v", msg)

	// after 1 second assume no message is coming
	case <-time.After(time.Second * 1):
		break
	}
}

// TestManyTCPClients tests to ensure that mist will send messages to many
// subscribers of the same tags, and then not send once unsubscribed
func TestManyTCPClients(t *testing.T) {

	//
	sender, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer sender.Close()

	//
	r1, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer r1.Close()

	//
	r2, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer r2.Close()

	//
	r3, err := New(mist.DEFAULT_ADDR)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer r3.Close()

	// receivers subscribe to tags and then sender publishes to those tags...
	err = r1.Subscribe([]string{testTag})
	err = r2.Subscribe([]string{testTag})
	err = r3.Subscribe([]string{testTag})
	if err != nil {
		t.Fatalf("one or more client subscription failed %v", err.Error())
	}

	//
	err = sender.Publish([]string{testTag}, testMsg)
	if err != nil {
		t.Fatalf("client publish failed %v", err.Error())
	}

	//
	waitMessage(r1.Messages(), t)
	waitMessage(r2.Messages(), t)
	waitMessage(r3.Messages(), t)

	// receiver unsubscribes from the tags and sender publishes again to the same
	// tags
	err = r1.Unsubscribe([]string{testTag})
	err = r2.Unsubscribe([]string{testTag})
	err = r3.Unsubscribe([]string{testTag})
	if err != nil {
		t.Fatalf("one or more client unsubscription failed %v", err.Error())
	}

	// allow time for communication across TCP connection
	<-time.After(time.Second * 1)

	//
	err = sender.Publish([]string{testTag}, testMsg)
	if err != nil {
		t.Fatalf("client publish failed %v", err.Error())
	}

	// receivers should NOT get a message this time
	select {

	// wait for a messages...
	case msg := <-r1.Messages():
		t.Fatalf("Received a message from unsubscribed tags: %#v", msg)

	case msg := <-r2.Messages():
		t.Fatalf("Received a message from unsubscribed tags: %#v", msg)

	case msg := <-r3.Messages():
		t.Fatalf("Received a message from unsubscribed tags: %#v", msg)

	// after 1 second assume no message is coming
	case <-time.After(time.Second * 1):
		break
	}
}

// waitMessage waits for a message to come to a proxy then tests to see if it is
// the expected message
func waitMessage(messages <-chan mist.Message, t *testing.T) {

	//
	select {

	// wait for a message then test to make sure it's the expected message...
	case msg := <-messages:
		if len(msg.Tags) != 1 {
			t.Fatalf("Wrong number of tags: Expected '%v' received '%v'\n", 1, len(msg.Tags))
		}
		if msg.Data != testMsg {
			t.Fatalf("Incorrect data: Expected '%v' received '%v'\n", testMsg, msg.Data)
		}
		break

	// after 1 second assume no messages are coming
	case <-time.After(time.Second * 1):
		t.Fatalf("Expecting messages, received none!")
	}
}
