package clients

import (
	"fmt"
	// "net"
	// "net/http"
	"testing"
	// "time"

	// "github.com/nanopack/mist/core"
	"github.com/nanopack/mist/server"
)

var (
	testTag     = "hello"
	testMsg     = "world"
	testAddress = "127.0.0.1:1234"
)

//
func TestTCPClient(t *testing.T) {

	//
	ln, err := server.NewTCP(testAddress, nil)
	if err != nil {
		t.Errorf("failed to listen - %v", err.Error())
	}
	defer ln.Close()

	//
	client, err := NewTCP(testAddress)
	if err != nil {
		t.Errorf("failed to conenct - %v", err.Error())
	}
	defer client.Close()

	//
	// if err := client.Ping(); err != nil {
	// 	t.Errorf("ping failed")
	// }

	//
	// client.Subscribe([]string{testTag})
	// client.Subscribe([]string{testTag, testMsg})

	//
	// list, err := client.List()
	// if err != nil {
	// 	t.Errorf("listing subsctiptions failed %v", err.Error())
	// }
	// if len(list) != 2 {
	// 	t.Errorf("wrong number of subscriptions were returned %v", list)
	// }
	// if len(list[0]) != 1 {
	// 	t.Errorf("wrong number of tags %v", list[0])
	// }
	// if len(list[1]) != 2 {
	// 	t.Errorf("wrong number of tags %v", list[1])
	// }

	//
	// client.Publish([]string{testTag}, testMsg)
	// message, ok := <-client.Messages()
	// if !ok {
	// 	t.Errorf("got a nil message")
	// }
	// if message.Data != testMsg {
	// 	t.Errorf("got the wrong message %v", message.Data)
	// }

	//
	// client.PublishAfter([]string{testTag}, testMsg, time.Second*1)
	// message, ok = <-client.Messages()
	// if !ok {
	// 	t.Errorf("got a nil message")
	// }
	// if message.Data != testMsg {
	// 	t.Errorf("got the wrong message %v", message.Data)
	// }
	//
	// fmt.Println("TEST ENDING!!")

	fmt.Println("CLOSING THINGS!")
}

//
// func TestWebsocketClient(t *testing.T) {
//
// 	fmt.Println("TEST STARTING!")
//
// 	//
// 	ln, err := net.Listen("tcp", "127.0.0.1:2345/")
// 	if err != nil {
// 		t.Errorf("unable to listen to websockets %v", err)
// 	}
// 	defer ln.Close()
//
// 	//
// 	go http.Serve(ln, server.ListenWS(nil))
//
// 	//
// 	ws, err := NewWS("ws://127.0.0.1:2345/", nil)
// 	if err != nil {
// 		t.Errorf("unable to connect %v", err)
// 	}
// 	defer ws.Close()
//
// 	//
// 	if err := ws.Subscribe([]string{testTag}); err != nil {
// 		t.Errorf("subscription failed %v", err)
// 	}
//
// 	//
// 	mist.Self.Publish([]string{testTag}, testMsg)
// 	<-ws.Messages()
//
// 	//
// 	list, err := ws.List()
// 	if err != nil {
// 		t.Errorf("unable to list %v", err)
// 	}
// 	if len(list) != 1 {
// 		t.Errorf("list of subscriptions is wrong %v", list)
// 	}
// 	if len(list[0]) != 1 {
// 		t.Errorf("wrong number of tags in subscription %v", list[0])
// 	}
//
// 	//
// 	err = ws.Unsubscribe([]string{testTag})
//
// 	//
// 	list, err = ws.List()
// 	if err != nil {
// 		t.Errorf("unable to list %v", err)
// 	}
// 	if len(list) != 0 {
// 		t.Errorf("list of subscriptions is wrong %v", list)
// 	}
// }
