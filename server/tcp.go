package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/util"
)

// init
func init() {

	// add tcp as an available server type
	listeners["tcp"] = startTCP
}

// startTCP starts a tcp server listening on the specified address (default 127.0.0.1:1445)
// and then continually reads from the server handling any incoming connections
func startTCP(uri string, errChan chan<- error) {

	//
	if uri == "" {
		uri = mist.DEFAULT_ADDR
	}

	// start a TCP listener
	ln, err := net.Listen("tcp", uri)
	if err != nil {
		errChan <- fmt.Errorf("Failed to start tcp listener %v", err.Error())
		return
	}
	fmt.Printf("TCP server listening at '%s'...\n", uri)

	// start continually listening for any incomeing tcp connections (non-blocking)
	go func() {
		for {

			// accept connections
			conn, err := ln.Accept()
			if err != nil {
				errChan <- fmt.Errorf("Failed to accept connection %v", err.Error())
				return
			}

			// handle each connection individually (non-blocking)
			go handleConnection(conn, errChan)
		}
	}()
}

// handleConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish Handler'
// that is used to publish messages to the data channel of the subscription
func handleConnection(conn net.Conn, errChan chan<- error) {

	// close the connection when we're done here
	defer conn.Close()

	// create a new client for each connection
	proxy := mist.NewProxy()
	defer proxy.Close()

	// publish mist messages to connected tcp clients (non-blocking)
	go func() {
		for msg := range proxy.Messages() {
			if _, err := conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(msg.Tags, ","), msg.Data))); err != nil {
				errChan <- fmt.Errorf("Failed to write to connection %v", err)
			}
		}
	}()

	// add basic TCP command handlers for this connection
	handlers = GenerateHandlers()

	//
	r := util.NewReader(conn)

	// check for authentication
	switch {

	// authentication wanted...
	case auth.DefaultAuth != nil:

		// prompt for auth
		if _, err := conn.Write([]byte("auth: ")); err != nil {
			break
		}

		// read from the connection looking for an auth token; if anything that comes
		// across the connection at this point and it's not the auth token, nothing
		// proceeds
		for r.Next() {

			// if the next input does not match the token then promt again
			if r.Input.Cmd != auth.Token {
				if _, err := conn.Write([]byte("auth: ")); err != nil {
					break
				}
				continue
			}

			// successful auth; allow auth command handlers on this connection
			for k, v := range auth.GenerateHandlers() {
				handlers[k] = v
			}

			//
			break
		}

	// no authentication wanted; authorize the proxy
	default:
		proxy.Authorized = true
	}

	// connection loop (blocking); continually read off the connection. Once something
	// is read, check to see if it is a message the client understands to be one of
	// its commands. If so attempt to execute the command.
	for r.Next() {

		// what should we do with this error?
		if r.Err != nil {
			errChan <- fmt.Errorf("Read error %v", r.Err)
		}

		fmt.Printf("READ!!! %#v\n", r.Input)

		//
		handler, found := handlers[r.Input.Cmd]

		//
		var err error
		switch {

		// command not found
		case !found:
			err = fmt.Errorf("Unknown Command '%s'\n", r.Input.Cmd)

		// wrong number of args
		case handler.NumArgs != len(r.Input.Args):
			err = fmt.Errorf("Wrong number of arguments for '%s'. Expected %v got %v\n", r.Input.Cmd, handler.NumArgs, len(r.Input.Args))

		// execute the command
		default:
			err = handler.Handle(proxy, r.Input.Args)
		}

		// if something failed along the way, respond accordingly...
		if err != nil {
			if _, err := conn.Write([]byte(err.Error())); err != nil {
				errChan <- fmt.Errorf("Failed to write to connection %v", err.Error())
				break
			}

			//
			continue
		}

		// ...otherwise write a successful response
		if _, err := conn.Write([]byte(fmt.Sprintf("%v success\n", r.Input.Cmd))); err != nil {
			errChan <- fmt.Errorf("Failed to write to connection %v", err.Error())
			break
		}
	}
}
