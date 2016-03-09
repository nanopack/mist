package mist

import (
	"testing"
	"time"
)

var (
	testTag = "hello"
	testMsg = "world"
)

// TestPublish tests that the publish Publish method publishes to all subscribers
func TestPublish(t *testing.T) {

	//
	p1 := NewProxy()
	defer p1.Close()

	p2 := NewProxy()
	defer p2.Close()

	//
	err := p1.Subscribe([]string{testTag})
	err = p2.Subscribe([]string{testTag})
	if err != nil {
		t.Fatalf("one or more proxy subscribes failed %v", err.Error())
	}

	// have mist publish the message
	PublishAfter([]string{testTag}, testMsg, 1)

	//
	waitMessage(p1, t)
	waitMessage(p2, t)

	err = p1.Unsubscribe([]string{testTag})
	err = p2.Unsubscribe([]string{testTag})
	if err != nil {
		t.Fatalf("one or more proxy unsubscribes failed %v", err.Error())
	}

	// have mist publish the message
	PublishAfter([]string{testTag}, testMsg, 1)

	// proxies should NOT get a message this time
	waitNoMessage(p1, t)
	waitNoMessage(p2, t)
}

// BenchmarkMist
func BenchmarkMist(b *testing.B) {

	//
	p := NewProxy()
	defer p.Close()

	//
	p.Subscribe([]string{testTag})

	//
	b.ResetTimer()

	//
	for i := 0; i < b.N; i++ {
		p.Publish([]string{testTag}, testMsg)
		_ = <-p.Pipe
	}
}

// waitMessage waits for a message to come to a proxy then tests to see if it is
// the expected message
func waitMessage(p *Proxy, t *testing.T) {

	//
	select {

	// wait for a message then test to make sure it's the expected message...
	case msg := <-p.Pipe:
		if len(msg.Tags) != 1 {
			t.Fatalf("Wrong number of tags: Expected '%v' received '%v'\n", 1, len(msg.Tags))
		}
		if msg.Data != testMsg {
			t.Fatalf("Incorrect data: Expected '%v' received '%v'\n", testMsg, msg.Data)
		}
		break

	// after 1 second assume no messages are coming
	case <-time.After(time.Second * 1):
		t.Errorf("Expecting messages, received none!")
	}
}

// waitNoMessage waits to NOT receive a message
func waitNoMessage(p *Proxy, t *testing.T) {

	//
	select {

	// wait for a message...
	case msg := <-p.Pipe:
		t.Fatalf("Received a message from unsubscribed tags: %#v", msg)

	// after 1 second assume no message is coming
	case <-time.After(time.Second * 1):
		break
	}
}
