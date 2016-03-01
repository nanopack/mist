package mist

import (
// "testing"
// "time"
)

var (
	testTag = "hello"
	testMsg = "world"
)

//
// func BenchmarkMist(b *testing.B) {
//
// 	//
// 	p, _ := NewProxy(0)
// 	defer p.Close()
//
// 	p.Subscribe([]string{testTag})
//
// 	b.ResetTimer()
//
// 	for i := 0; i < b.N; i++ {
// 		Self.Publish([]string{testTag}, testMsg)
// 		_ = <-p.Messages()
// 	}
// }

//
// func TestMist(t *testing.T) {
//
// 	//
// 	p, _ := NewProxy(0)
// 	defer p.Close()
//
// 	//
// 	p.Subscribe([]string{testTag})
//
// 	//
// 	for count := 0; count < 2; count++ {
//
// 		//
// 		Self.Publish([]string{testTag}, testMsg)
// 		message := <-p.Messages()
// 		if len(message.Tags) != 1 {
// 			t.Errorf("wrong number of tags")
// 		}
// 		if message.Data != testMsg {
// 			t.Errorf("incorrect data")
// 		}
// 	}
//
// 	//
// 	p.Unsubscribe([]string{testTag})
//
// 	//
// 	Self.Publish([]string{testTag}, testMsg)
// 	select {
// 	case <-p.Messages():
// 		t.Errorf("the message should not have been received")
//
// 	// wait 1 second before assuming that nothing is coming across the wire
// 	case <-time.After(time.Second * 1):
// 		break
// 	}
// }

//
// func TestMistReplication(t *testing.T) {
//
// 	//
// 	p, _ := NewProxy(0)
// 	defer p.Close()
//
// 	r1, _ := NewProxy(0)
// 	defer r1.Close()
//
// 	r2, _ := NewProxy(0)
// 	defer r2.Close()
//
// 	// two clients will represent remote replicated nodes
// 	r1.(Replicatable).EnableReplication()
// 	r2.(Replicatable).EnableReplication()
//
// 	p.Subscribe([]string{"foo"})
// 	r1.Subscribe([]string{"foo"})
// 	r2.Subscribe([]string{"foo"})
//
// 	// when a normal client publishes, both replicated clients receive the message
// 	p.Publish([]string{"foo"}, "data")
// 	<-r1.Messages()
// 	<-r2.Messages()
// 	<-p.Messages()
//
// 	//
// 	r1.Publish([]string{"foo"}, "data")
// 	select {
// 	case <-p.Messages():
//
// 	// a p client should get messages from a replicated client
// 	case <-r2.Messages():
// 		t.Error("a replicated client should not get a message from another replicated client")
// 	}
//
// 	//
// 	r2.Publish([]string{"foo"}, "data")
// 	select {
// 	case <-p.Messages():
//
// 	// a p client should get messages from a replicated client
// 	case <-r1.Messages():
// 		t.Error("a replicated client should not get a message from another replicated client")
// 	}
// }
