package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
)

// init adds "tcp" as an available mist server type
func init() {
	Register("tcp", StartTCP)
}

// StartTCP starts a tcp server listening on the specified address (default 127.0.0.1:1445)
// and then continually reads from the server handling any incoming connections
func StartTCP(uri string, errChan chan<- error) {

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

	// add basic TCP command handlers for this connection
	handlers = GenerateHandlers()

	//
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// publish mist messages to connected tcp clients (non-blocking)
	go func() {
		for msg := range proxy.Pipe {
			if err := encoder.Encode(msg); err != nil {
				errChan <- fmt.Errorf(err.Error())
				continue
			}
		}
	}()

	// connection loop (blocking); continually read off the connection. Once something
	// is read, check to see if it's a message the client understands to be one of
	// its commands. If so attempt to execute the command.
	for decoder.More() {

		//
		msg := mist.Message{}

		// decode an array value (Message)
		if err := decoder.Decode(&msg); err != nil {
			errChan <- fmt.Errorf(err.Error())
			continue
		}

		// read from the connection looking for an auth token; if anything that comes
		// across the connection at this point and it's not the auth token, nothing
		// proceeds
		if auth.DefaultAuth != nil {

			// if the next input does not match the token then bugout
			if msg.Data != token {
				return
			}

			// successful auth; allow auth command handlers on this connection
			for k, v := range auth.GenerateHandlers() {
				handlers[k] = v
			}
		}

		// no authentication wanted; authorize the proxy
		// proxy.Authorized = true

		// look for the command
		handler, found := handlers[msg.Command]

		// if the command isn't found, return an error
		if !found {
			encoder.Encode(&mist.Message{Command: msg.Command, Error: "Unknown Command"})
			continue
		}

		// attempt to run the command
		if err := handler(proxy, msg); err != nil {
			encoder.Encode(&mist.Message{Command: msg.Command, Error: err.Error()})
		}
	}
}
