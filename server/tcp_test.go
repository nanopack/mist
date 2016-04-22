package server

import (
	"testing"

	"github.com/nanopack/mist/auth"
)

// TestTCPStart tests to ensure a server will start
func TestTCPStart(t *testing.T) {

	// ensure authentication is disabled
	auth.DefaultAuth = nil

	if err := Start([]string{"tcp://127.0.0.1:1445"}, ""); err != nil {
		t.Fatalf("Unexpected error - %v", err.Error())
	}
}
