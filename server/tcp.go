package server

import (
	"fmt"
	"net"
	"strings"

	// "github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/util"
)

// QUESTION! how/where will other custom handlers be passed???

var handlers map[string]Handler

// start a mist server listening over TCP
func startTCP(uri string, errChan chan<- error) {

	// get basic TCP command handlers
	handlers = GenerateHandlers()

	// add any additional commands to existing tcp commands
	// for k, v := range additionalHandlers {
	// 	handlers[k] = v
	// }

	// if authenticator.DefaultAuth != nil {
	// 	// get auth handlers and merge with exsiting
	// }

	// don't close this because it's go routined and the error is handled from the
	// errChan
	if err := newTCP(uri, handlers); err != nil {
		errChan<- fmt.Errorf("Unable to start mist tcp listener %v", err)
	}
}

// newTCP starts a tcp server listening on the specified address (default 127.0.0.1:1445)
// and then continually reads from the server handling any incoming connections;
// this is intentionally non-blocking.
func newTCP(address string, additionalHandlers map[string]Handler) error {

	//
	if address == "" {
		address = mist.DEFAULT_ADDR
	}

	// start a TCP listener
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Printf("TCP server listening at '%s'...\n", address)

	// start continually listening for any incomeing tcp connections (non-blocking)
	go func() {
		for {

			// accept connections
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("TCPS BONK!", err) // what should we do with the error?
			}

			// handle each connection individually (non-blocking)
			go handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish Handler'
// that is used to publish messages to the data channel of the subscription
func handleConnection(conn net.Conn) {

	fmt.Printf("TCPS HANDLE CONNECTION! %q\n", conn)

	// create a new client for each connection
	proxy := mist.NewProxy(0)

	// clean up everything that we have setup
	defer func() {
		conn.Close()
		proxy.Close()
	}()

	// add a publisher that will publish across the connection (non-blocking)
	go publishHandler(proxy, conn)

	// add a reader that reads off the connection (blocking)
	readHandler(proxy, conn)

	fmt.Println("TCPS END!")
}

// publishHandler is used to...
func publishHandler(proxy mist.Proxy, conn net.Conn) {

	fmt.Println("TCPS PUBLISHING!")

	// make a done channel
	done := make(chan bool)
	defer close(done)

	for {

		//
		select {

		// return if we are done
		case <-done:
			fmt.Println("TCPS DONE!")
			return

		// when a message is recieved on the subscriptions channel write the message
		// to the connection
		case msg := <-proxy.Messages():
			fmt.Printf("TCPS MESSAGE! %#v\n", msg)
			if _, err := conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(msg.Tags, ","), msg.Data))); err != nil {
				break
			}
		}
	}
}

// readHandler is used to read off the open connection and execute any recongnized
// commands that come across
func readHandler(proxy mist.Proxy, conn net.Conn) {

	fmt.Println("TCPS READING!")

	// continually read off the connection; once something is read, check to see
	// if it is a message the client understands to be one of its commands. If so
	// execute the command.
	r := util.NewReader(conn)
	for r.Next() {

		fmt.Printf("TCPS NEXT! %#v\n", r)

		// what should we do with this error?
		if r.Err != nil {
			fmt.Println("ERROR!", r.Err)
		}

		//
		handler, found := handlers[r.Input.Cmd]

		//
		var response string
		switch {

		// no command found
		case !found:
			response = fmt.Sprintf("Error: Unknown Command '%s'", r.Input.Cmd)
			// continue

		//
		case handler.NumArgs != len(r.Input.Args):
			response = fmt.Sprintf("Error: Wrong number of arguments for '%s'. Expected %v got %v", r.Input.Cmd, handler.NumArgs, len(r.Input.Args))

		// execute command
		default:
			fmt.Println("TCPS EXECUTE! ", r.Input.Cmd)
			response = handler.Handle(proxy, r.Input.Args)
		}

		// write the response from the command back to the connection
		fmt.Println("TCPS WRITING RESPONSE! ", response)
		if _, err := conn.Write([]byte(response + "\n")); err != nil {
			break
		}
	}

	fmt.Println("READING DONE!")
}
