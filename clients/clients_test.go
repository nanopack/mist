package clients

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/server"
)

var (
	testAddr = "127.0.0.1:1445"
	testTag  = "hello"
	testMsg  = "world"
)

// TestMain
func TestMain(m *testing.M) {

	//
	server.StartTCP(testAddr, nil)

	//
	os.Exit(m.Run())
}

// TestTCPClient tests
func TestTCPClientSubscriptions(t *testing.T) {

	//
	client, err := New(testAddr)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer client.Close()

	// ping...
	if err := client.Ping(); err != nil {
		t.Fatalf("ping failed")
	}

	// ...should return pong
	if msg := <-client.Messages(); msg.Data != "pong" {
		t.Fatalf("Unexpected data: Expecting 'pong' got %s", msg.Data)
	}

	// subscribe...

	// ...should fail with no tags
	if err := client.Subscribe([]string{}); err == nil {
		t.Fatalf("Subscription succeeded with missing tags!")
	}

	// ...should subscribe a single tag
	if err := client.Subscribe([]string{testTag}); err != nil {
		t.Fatalf("client subscriptions failed %v", err.Error())
	}

	// ...should subscribe multiple tags
	if err := client.Subscribe([]string{testTag, testMsg}); err != nil {
		t.Fatalf("client subscriptions failed %v", err.Error())
	}

	// ...should list two sets of tags
	if err := client.List(); err != nil {
		t.Fatalf("listing subscriptions failed %v", err.Error())
	}
	if msg := <-filterChannel(client.Messages(), t); msg.Data != fmt.Sprintf("%s %s,%s", testTag, testTag, testMsg) {
		t.Fatalf("Incorrect subscriptions: Expected %v got %v", fmt.Sprintf("%s %s,%s", testTag, testTag, testMsg), msg.Data)
	}

	// unsubscribe...

	// ...should remove a single tag
	if err := client.Unsubscribe([]string{testTag}); err != nil {
		t.Fatalf("client unsubscriptions failed %v", err.Error())
	}

	// ...should remove multiple tags
	if err := client.Unsubscribe([]string{testTag, testMsg}); err != nil {
		t.Fatalf("client unsubscriptions failed %v", err.Error())
	}

	// ...should list no tags
	if err := client.List(); err != nil {
		t.Fatalf("listing subscriptions failed %v", err.Error())
	}
	if msg := <-filterChannel(client.Messages(), t); msg.Data != "" {
		t.Fatalf("Unexpected subscriptions: %v", msg.Data)
	}
}

//
func filterChannel(msgChan <-chan mist.Message, t *testing.T) <-chan mist.Message {
	rtn := make(chan mist.Message)
	go func() {
		for msg := range msgChan {
			fmt.Printf("MESSAGE! %#v\n", msg)
			if msg.Data != "success" {
				fmt.Printf("DATA! %#v\n", msg)
				rtn <- msg
			}
			if msg.Error != "" {
				t.Fatalf("Unexpected error: %v", msg.Error)
			}
		}
		close(rtn)
	}()

	return rtn
}

// TestSameTCPClient tests to ensure that mist will not send message to the
// same client who publishes them
// func TestSameTCPClient(t *testing.T) {
//
// 	//
// 	sender, err := New(testAddr)
// 	if err != nil {
// 		t.Fatalf("failed to connect - %v", err.Error())
// 	}
// 	defer sender.Close()
//
// 	// sender subscribes to tags and then tries to publish to those same tags...
// 	if err := sender.Subscribe([]string{testTag}); err != nil {
// 		t.Fatalf("client subscription failed %v", err.Error())
// 	}
// 	if err := sender.Publish([]string{testTag}, testMsg); err != nil {
// 		t.Fatalf("client publish failed %v", err.Error())
// 	}
//
// 	// sender should NOT get a message because mist shouldnt send a message to the
// 	// same proxy that publishes them.
// 	for msg := range sender.Messages() {
// 		fmt.Printf("MESSAGE! %#v\n", msg)
// 	}
// 	for {
// 		select {
//
// 		// wait for a message...
// 		case msg := <-sender.Messages():
// 			if msg.Data != "success" {
// 				t.Fatalf("Received own message! %#v\n", msg)
// 			}
//
// 		// after 1 second assume no message is coming
// 		case <-time.After(time.Second * 1):
// 			break
// 		}
// 	}
// }

// TestDifferentTCPClient tests to ensure that mist will send messages
// to another subscribed client, and then not send when unsubscribed.
// func TestDifferentTCPClient(t *testing.T) {
//
// 	//
// 	sender, err := New(testAddr)
// 	if err != nil {
// 		t.Fatalf("failed to connect - %v", err.Error())
// 	}
// 	defer sender.Close()
//
// 	//
// 	receiver, err := New(testAddr)
// 	if err != nil {
// 		t.Fatalf("failed to connect - %v", err.Error())
// 	}
// 	defer receiver.Close()
//
// 	// receiver subscribes to tags and then sender publishes to those tags...
// 	err = receiver.Subscribe([]string{testTag})
// 	if err != nil {
// 		t.Fatalf("client subscription failed %v", err.Error())
// 	}
//
// 	// publish after to allow time for communication across TCP connection
// 	err = sender.PublishAfter([]string{testTag}, testMsg, 1)
// 	if err != nil {
// 		t.Fatalf("client publish failed %v", err.Error())
// 	}
//
// 	// allow time for communication across TCP connection
// 	<-time.After(time.Second * 1)
//
// 	//
// 	verifyMessage(receiver.Messages(), t)
//
// 	// receiver unsubscribes from the tags and sender publishes again to the same
// 	// tags
// 	err = receiver.Unsubscribe([]string{testTag})
// 	if err != nil {
// 		t.Fatalf("client unsubscribe failed %v", err.Error())
// 	}
//
// 	// publish after to allow time for communication across TCP connection
// 	err = sender.PublishAfter([]string{testTag}, testMsg, 1)
// 	if err != nil {
// 		t.Fatalf("client publish failed %v", err.Error())
// 	}
//
// 	// allow time for communication across TCP connection
// 	<-time.After(time.Second * 1)
//
// 	// receiver should NOT get a message this time
// 	verifyNoMessage(receiver.Messages(), t)
// }

// TestManyTCPClients tests to ensure that mist will send messages to many
// subscribers of the same tags, and then not send once unsubscribed
// func TestManyTCPClients(t *testing.T) {
//
// 	//
// 	sender, err := New(testAddr)
// 	if err != nil {
// 		t.Fatalf("failed to connect - %v", err.Error())
// 	}
// 	defer sender.Close()
//
// 	//
// 	r1, err := New(testAddr)
// 	if err != nil {
// 		t.Fatalf("failed to connect - %v", err.Error())
// 	}
// 	defer r1.Close()
//
// 	//
// 	r2, err := New(testAddr)
// 	if err != nil {
// 		t.Fatalf("failed to connect - %v", err.Error())
// 	}
// 	defer r2.Close()
//
// 	//
// 	r3, err := New(testAddr)
// 	if err != nil {
// 		t.Fatalf("failed to connect - %v", err.Error())
// 	}
// 	defer r3.Close()
//
// 	// receivers subscribe to tags and then sender publishes to those tags...
// 	err = r1.Subscribe([]string{testTag})
// 	err = r2.Subscribe([]string{testTag})
// 	err = r3.Subscribe([]string{testTag})
// 	if err != nil {
// 		t.Fatalf("one or more client subscription failed %v", err.Error())
// 	}
//
// 	// publish after to allow time for communication across TCP connection
// 	err = sender.PublishAfter([]string{testTag}, testMsg, 1)
// 	if err != nil {
// 		t.Fatalf("client publish failed %v", err.Error())
// 	}
//
// 	// allow time for communication across TCP connection
// 	<-time.After(time.Second * 1)
//
// 	//
// 	verifyMessage(r1.Messages(), t)
// 	verifyMessage(r2.Messages(), t)
// 	verifyMessage(r3.Messages(), t)
//
// 	// receiver unsubscribes from the tags and sender publishes again to the same
// 	// tags
// 	err = r1.Unsubscribe([]string{testTag})
// 	err = r2.Unsubscribe([]string{testTag})
// 	err = r3.Unsubscribe([]string{testTag})
// 	if err != nil {
// 		t.Fatalf("one or more client unsubscription failed %v", err.Error())
// 	}
//
// 	// publish after to allow time for communication across TCP connection
// 	err = sender.PublishAfter([]string{testTag}, testMsg, 1)
// 	if err != nil {
// 		t.Fatalf("client publish failed %v", err.Error())
// 	}
//
// 	// allow time for communication across TCP connection
// 	<-time.After(time.Second * 1)
//
// 	// receivers should NOT get a message this time
// 	verifyNoMessage(r1.Messages(), t)
// 	verifyNoMessage(r2.Messages(), t)
// 	verifyNoMessage(r3.Messages(), t)
// }

// verifyMessage waits for a message to come to a proxy then tests to see if it is
// the expected message
func verifyMessage(messages <-chan mist.Message, t *testing.T) {

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

// verifyNoMessage waits to NOT receive a message
func verifyNoMessage(messages <-chan mist.Message, t *testing.T) {

	//
	select {

	// wait for a message...
	case msg := <-messages:
		t.Fatalf("Received a message from unsubscribed tags: %#v", msg)

	// after 1 second assume no message is coming
	case <-time.After(time.Second * 1):
		break
	}
}
