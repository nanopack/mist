package clients

import (
  "fmt"
	"net"
	"net/http"
	"testing"
	"time"

  "github.com/nanopack/mist/core"
  "github.com/nanopack/mist/server"
)

var (
	testTag = "hello"
	testMsg = "world"
  testAddress = "127.0.0.1:1234"
)

//
func TestTCPClient(t *testing.T) {

  //
  mist := mist.New()
  fmt.Println("MIST!", mist)

  //
	server, err := server.ListenTCP(testAddress, nil)
  defer server.Close()

  if err != nil {
    t.Errorf("failed to listen - %v", err.Error())
  }

  //
	tcp, err := NewTCP(testAddress)
  if err != nil {
    t.Errorf("failed to conenct - %v", err.Error())
  }
  defer tcp.Close()

  // connection delay
	// <-time.After(time.Millisecond * 10)

  if err := tcp.Ping(); err != nil {
    t.Errorf("ping failed")
  }

  //
	tcp.Subscribe([]string{testTag})
	tcp.Subscribe([]string{testTag, testMsg})

  //
	// list, err := tcp.List()
  //
  // if err != nil {
  //   t.Errorf("listing subsctiptions failed %v", err.Error())
  // }
  // if len(list) != 2 {
  //   t.Errorf("wrong number of subscriptions were returned %v", list)
  // }
  // if len(list[0]) != 1 {
  //   t.Errorf("wrong number of tags %v", list[0])
  // }
  // if len(list[1]) != 2 {
  //   t.Errorf("wrong number of tags %v", list[1])
  // }

  //
  tcp.Publish([]string{testTag}, testMsg)
	message, ok := <-tcp.Messages()

  if !ok {
    t.Errorf("got a nil message")
  }
  if message.Data != testMsg {
    t.Errorf("got the wrong message %v", message.Data)
  }

  //
	tcp.PublishAfter([]string{testTag}, testMsg, time.Second*1)
	message, ok = <-tcp.Messages()

  if !ok {
    t.Errorf("got a nil message")
  }
  if message.Data != testMsg {
    t.Errorf("got the wrong message %v", message.Data)
  }
}

//
func TestWebsocketClient(t *testing.T) {

  //
  mist := mist.New()

  //
	ln, err := net.Listen("tcp", testAddress)
  if err != nil {
    t.Errorf("unable to listen to websockets %v", err)
  }
	defer ln.Close()

  //
	go http.Serve(ln, server.ListenWS(nil))

  //
	ws, err := NewWS("ws://127.0.0.1:1234/", nil)
  if err != nil {
    t.Errorf("unable to connect %v", err)
  }
	defer ws.Close()

  //
	err = ws.Subscribe([]string{testTag})

  if err != nil {
    t.Errorf("subscription failed %v", err)
  }

  //
	mist.Publish([]string{testTag}, testMsg)
	<-ws.Messages()

  //
	list, err := ws.List()

  if err != nil {
    t.Errorf("unable to list %v", err)
  }
  if len(list) != 1 {
    t.Errorf("list of subscriptions is wrong %v", list)
  }
  if len(list[0]) != 1 {
    t.Errorf("wrong number of tags in subscription %v", list[0])
  }

  //
  err = ws.Unsubscribe([]string{testTag})

  //
	list, err = ws.List()

  if err != nil {
    t.Errorf("unable to list %v", err)
  }
  if len(list) != 0 {
    t.Errorf("list of subscriptions is wrong %v", list)
  }
}
