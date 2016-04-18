package mist

import (
	"testing"
	"time"
)

// TestSameSubscriber tests to ensure that mist will not send message to the
// same proxy who publishes them
func TestSameSubscriber(t *testing.T) {

	//
	sender := NewProxy()
	defer sender.Close()

	// sender subscribes to tags and then tries to publish to those same tags...
	sender.Subscribe([]string{testTag})
	sender.Publish([]string{testTag}, testMsg)

	// sender should NOT get a message because mist shouldnt send a message to the
	// same proxy that publishes them.
	select {

	// wait for a message...
	case <-sender.Pipe:
		t.Fatalf("Received own message!")

	// after 1 second assume no message is coming
	case <-time.After(time.Second * 1):
		break
	}

	//
	sender.Unsubscribe([]string{testTag})
}

// TestDifferentSubscriber tests to ensure that mist will send messages
// to another subscribed proxy, and then not send when unsubscribed.
func TestDifferentSubscriber(t *testing.T) {

	//
	sender := NewProxy()
	defer sender.Close()

	//
	receiver := NewProxy()
	defer receiver.Close()

	// receiver subscribes to tags and then sender publishes to those tags...
	receiver.Subscribe([]string{testTag})
	sender.Publish([]string{testTag}, testMsg)

	//
	waitMessage(receiver, t)

	// receiver unsubscribes from the tags and sender publishes again to the same
	// tags
	receiver.Unsubscribe([]string{testTag})
	sender.Publish([]string{testTag}, testMsg)

	// receiver should NOT get a message this time
	select {

	// wait for a message...
	case msg := <-receiver.Pipe:
		t.Fatalf("Received a message from unsubscribed tags: %#v", msg)

	// after 1 second assume no message is coming
	case <-time.After(time.Second * 1):
		break
	}
}

// TestManySubscribers tests to ensure that mist will send messages to many
// subscribers of the same tags, and then not send once unsubscribed
func TestManySubscribers(t *testing.T) {

	//
	sender := NewProxy()
	defer sender.Close()

	//
	r1 := NewProxy()
	defer r1.Close()

	//
	r2 := NewProxy()
	defer r2.Close()

	//
	r3 := NewProxy()
	defer r3.Close()

	// receivers subscribe to tags and then sender publishes to those tags...
	r1.Subscribe([]string{testTag})
	r2.Subscribe([]string{testTag})
	r3.Subscribe([]string{testTag})

	//
	sender.Publish([]string{testTag}, testMsg)

	//
	waitMessage(r1, t)
	waitMessage(r2, t)
	waitMessage(r3, t)

	// receiver unsubscribes from the tags and sender publishes again to the same
	// tags
	r1.Unsubscribe([]string{testTag})
	r2.Unsubscribe([]string{testTag})
	r3.Unsubscribe([]string{testTag})

	//
	sender.Publish([]string{testTag}, testMsg)

	// receivers should NOT get a message this time
	waitNoMessage(r1, t)
	waitNoMessage(r2, t)
	waitNoMessage(r3, t)
}
