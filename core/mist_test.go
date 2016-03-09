package mist

import (
	"testing"
)

var (
	testTag = "hello"
	testMsg = "world"
)

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
