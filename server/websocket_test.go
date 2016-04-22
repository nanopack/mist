package server

import (
	"testing"

	"github.com/nanopack/mist/auth"
)

// TestWSStart tests to ensure a server will start
func TestWSStart(t *testing.T) {

	// ensure authentication is disabled
	auth.DefaultAuth = nil

	if err := Start([]string{"ws://127.0.0.1:8888"}, ""); err != nil {
		t.Fatalf("Unexpected error - %v", err.Error())
	}
}
