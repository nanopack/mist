// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"testing"
)

func TestMistCore(test *testing.T) {
	mist := New()
	client := mist.Client(0)
	defer client.Close()

	client.Subscribe([]string{"tag0"})

	for count := 0; count < 2; count++ {
		mist.Publish([]string{"tag0"}, []byte("this is my data"))
		message := <-client.Messages()
		assert(test, len(message.Tags) == 1, "wrong number of tags")
		// assert(test, message.Data == []byte("this is my data"), "data was incorrect")
	}

	client.Unsubscribe([]string{"tag0"})
	mist.Publish([]string{"tag0"}, []byte("this is my data"))
	select {
	case <-client.Messages():
		assert(test, false, "the message should not have been received")
	default:
	}

}

func BenchmarkMistCore(b *testing.B) {
	mist := New()
	client := mist.Client(0)
	defer client.Close()
	client.Subscribe([]string{"tag0"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mist.Publish([]string{"tag0"}, []byte("this is my data"))
		_ = <-client.Messages()
	}
}

func assert(test *testing.T, check bool, fmt string, args ...interface{}) {
	if !check {
		test.Logf(fmt, args...)
		test.FailNow()
	}
}
