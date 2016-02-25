package server

import (
  "fmt"
  "net"
  "strings"

  "github.com/nanopack/mist/core"
	"github.com/nanopack/mist/server/handlers"
	"github.com/nanopack/mist/util"
)

//
var tcpCommands map[string]mist.TCPHandler

//
func init() {

  // add TCP handlers
  tcpCommands = handlers.GenerateTCPCommands()
}

// ListenTCP starts a tcp server listening on the specified address (default 127.0.0.1:1445)
// and then continually reads from the server handling any incoming connections;
// this is intentionally non-blocking.
func ListenTCP(address string, mixins map[string]mist.TCPHandler) (net.Listener, error) {

	//
	if address == "" {
		address = mist.DEFAULT_ADDR
	}

	// start a TCP listener
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
  fmt.Printf("TCP server listening at '%s'...\n", address)

	// non-blocking
	go func() {

    // add any additional handlers
    for k, v := range mixins {
  		tcpCommands[k] = v
  	}

    //
    // defer client.Close()

		// Continually listen for any incoming connections.
		for {

			// accept connections
			conn, err := ln.Accept()
			if err != nil {
        fmt.Println("BONK!", err) // what should we do with the error?
			}

			// handle each connection individually (non-blocking)
			go handleConnection(conn)
		}
	}()

	return ln, nil
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish Handler'
// that is used to publish messages to the data channel of the subscription
func handleConnection(conn net.Conn) {

  fmt.Printf("NEW CONNECTION! %q\n", conn)

  // create a new client for each connection
  proxy, err := mist.NewProxy(0)
  if err != nil {
    fmt.Println("BONK!", err) // what should we do with the error?
  }

	// clean up everything that we have setup
	defer func() {
    fmt.Println("CALLED?????????")
		conn.Close()
		proxy.Close()
	}()

	// add a "publisher" for this connection (non-blocking)
	go publisher(proxy, conn)

	// add a "reader" for the connection (blocking)
	reader(proxy, conn)

  fmt.Println("DONE???")
}

//
func publisher(proxy mist.Client, conn net.Conn) {

  // make a done channel
	done := make(chan bool)
  defer close(done)

	for {

		//
		select {

		// when a message is recieved on the subscriptions channel write the message
		// to the connection
    case msg := <-proxy.Messages():
      fmt.Printf("TCP MESSAGE! %#v\n", msg)
			if _, err := conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(msg.Tags, ","), msg.Data))); err != nil {
				break
			}

		// return if we are done
		case <-done:
      fmt.Println("TCP DONE!")
			return
		}
	}
}

//
func reader(proxy mist.Client, conn net.Conn) {

	//
	r := util.NewReader(conn)

	//
	for r.Next() {

    fmt.Printf("TCP NEXT! %#v\n", r)

		// what should we do with this error?
		if r.Err != nil {
			fmt.Println("ERROR!", r.Err)
		}

		//
		handler, found := tcpCommands[r.Input.Cmd]

		//
		var response string
		switch {

		// no command found
		case !found:
			response = fmt.Sprintf("Error: Unknown Command '%s'", r.Input.Cmd)

    //
    case handler.NumArgs != len(r.Input.Args):
      response = fmt.Sprintf("Error: Wrong number of arguments for '%s'. Expected %v got %v", r.Input.Cmd, handler.NumArgs, len(r.Input.Args))

		// execute command
		default:
      fmt.Println("EXECUTE CLIENT! ", r.Input.Cmd)
			response = handler.Handle(proxy, r.Input.Args)
		}

		// write the response from the command back to the connection
    fmt.Println("TCP WRITING RESPONSE! ", response)
		if _, err := conn.Write([]byte(response + "\n")); err != nil {
			break
		}
	}
}
