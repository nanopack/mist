package server

import (
	"testing"

	"github.com/nanopack/mist/auth"
)

// TestHTTPStart tests to ensure a server will start
func TestHTTPStart(t *testing.T) {

	// ensure authentication is disabled
	auth.DefaultAuth = nil

	if err := Start([]string{"http://127.0.0.1:8080"}, ""); err != nil {
		t.Fatalf("Unexpected error - %v", err.Error())
	}
}
