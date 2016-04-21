package clients

import (
	"os"
	"testing"

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

// TestTCPClientConnect tests to ensure a client can connect to a running server
func TestTCPClientConnect(t *testing.T) {
	client, err := New(testAddr)
	if err != nil {
		t.Fatalf("Client failed to connect - %v", err.Error())
	}
	defer client.Close()

	//
	if err := client.Ping(); err != nil {
		t.Fatalf("ping failed")
	}
	if msg := <-client.Messages(); msg.Data != "pong" {
		t.Fatalf("Unexpected data: Expecting 'pong' got %s", msg.Data)
	}
}

// TestTCPClient tests to ensure a client can run all of its expected commands;
// we don't have to actually test any of the results of the commands since those
// are already tested in other tests (proxy_test and subscriptions_test in the
// core package)
func TestTCPClient(t *testing.T) {

	//
	client, err := New(testAddr)
	if err != nil {
		t.Fatalf("failed to connect - %v", err.Error())
	}
	defer client.Close()

	// subscribe should fail with no tags
	if err := client.Subscribe([]string{}); err == nil {
		t.Fatalf("Subscription succeeded with missing tags!")
	}

	// test ability to subscribe
	if err := client.Subscribe([]string{"a"}); err != nil {
		t.Fatalf("client subscriptions failed %v", err.Error())
	}
	if msg := <-client.Messages(); msg.Data != "success" {
		t.Fatalf("Failed to 'subscribe' - %v", msg.Error)
	}

	// test ability to list (subscriptions)
	if err := client.List(); err != nil {
		t.Fatalf("listing subscriptions failed %v", err.Error())
	}
	if msg := <-client.Messages(); msg.Data == "" {
		t.Fatalf("Failed to 'list' - %v", msg.Error)
	}

	// test ability to unsubscribe
	if err := client.Unsubscribe([]string{"a"}); err != nil {
		t.Fatalf("client unsubscriptions failed %v", err.Error())
	}
	if msg := <-client.Messages(); msg.Data != "success" {
		t.Fatalf("Failed to 'unsubscribe' - %v", msg.Error)
	}

	// test ability to list (no subscriptions)
	if err := client.List(); err != nil {
		t.Fatalf("listing subscriptions failed %v", err.Error())
	}
	if msg := <-client.Messages(); msg.Data != "" {
		t.Fatalf("Failed to 'list' - %v", msg.Error)
	}
}
