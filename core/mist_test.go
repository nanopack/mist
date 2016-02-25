package mist

import (
	// "fmt"
	"testing"
	"time"
)

var (
	testTag = "hello"
	testMsg = "world"
)

//
func setup() (*Mist, Client) {
	m := New()
	p, _ := NewProxy(0)
	return m, p
}

//
func BenchmarkMist(b *testing.B) {

	//
	mist, proxy := setup()
	defer proxy.Close()

	proxy.Subscribe([]string{testTag})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mist.Publish([]string{testTag}, testMsg)
		_ = <-proxy.Messages()
	}
}

//
func TestMist(t *testing.T) {

	//
	mist, proxy := setup()
	defer proxy.Close()

	//
	proxy.Subscribe([]string{testTag})

	//
	for count := 0; count < 2; count++ {

		//
		mist.Publish([]string{testTag}, testMsg)
		message := <-proxy.Messages()

		if len(message.Tags) != 1 {
			t.Errorf("wrong number of tags")
		}
		if message.Data != testMsg {
			t.Errorf("incorrect data")
		}
	}

	//
	proxy.Unsubscribe([]string{testTag})

	//
	mist.Publish([]string{testTag}, testMsg)

	select {
	case <-proxy.Messages():
		t.Errorf("the message should not have been received")

	// wait 1 second before assuming that nothing is coming across the wire
	case <-time.After(time.Second*1):
		break
	}
}

//
// func TestMistReplication(t *testing.T) {
//
// 	//
// 	_, proxy := setup()
// 	defer proxy.Close()
//
// 	_, replication1 := setup()
// 	defer replication1.Close()
//
// 	_, replication2 := setup()
// 	defer replication2.Close()
//
// 	// two clients will represent remote replicated nodes
// 	replication1.(Replicatable).EnableReplication()
// 	replication2.(Replicatable).EnableReplication()
//
// 	proxy.Subscribe([]string{"foo"})
// 	replication1.Subscribe([]string{"foo"})
// 	replication2.Subscribe([]string{"foo"})
//
// 	// when a normal client publishes, both replicated clients receive the message
// 	proxy.Publish([]string{"foo"}, "data")
// 	<-replication1.Messages()
// 	<-replication2.Messages()
// 	<-proxy.Messages()
//
// 	//
// 	replication1.Publish([]string{"foo"}, "data")
// 	select {
// 	case <-proxy.Messages():
// 		// a proxy client should get messages from a replicated client
// 	case <-replication2.Messages():
// 		t.Error("a replicated client should not get a message from another replicated client")
// 	}
//
// 	//
// 	replication2.Publish([]string{"foo"}, "data")
// 	select {
// 	case <-proxy.Messages():
// 		// a proxy client should get messages from a replicated client
// 	case <-replication1.Messages():
// 		t.Error("a replicated client should not get a message from another replicated client")
// 	}
// }
