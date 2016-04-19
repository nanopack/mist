package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/jcelliott/lumber"

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

	// start a TCP listener
	ln, err := net.Listen("tcp", uri)
	if err != nil {
		errChan <- fmt.Errorf("Failed to start tcp listener %v", err.Error())
		return
	}
	lumber.Info("TCP server listening at '%s'...\n", uri)

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

			// if the message fails to encode its probably a syntax issue and needs to
			// break the loop here because it will never be able to encode it; this will
			// disconnect the client.
			if err := encoder.Encode(msg); err != nil {
				errChan <- fmt.Errorf(err.Error())
				break
			}
		}
	}()

	// connection loop (blocking); continually read off the connection. Once something
	// is read, check to see if it's a message the client understands to be one of
	// its commands. If so attempt to execute the command.
	for decoder.More() {

		//
		msg := mist.Message{}

		// if the message fails to decode its probably a syntax issue and needs to
		// break the loop here because it will never be able to decode it; this will
		// disconnect the client.
		if err := decoder.Decode(&msg); err != nil {
			errChan <- fmt.Errorf(err.Error())
			break
		}

		// if an authenticator was passed, check for a token on connect to see if
		// auth commands are allowed
		if auth.DefaultAuth != nil && !authenticated {

			// if the next input does not match the token then...
			if msg.Data == authtoken {

				// if the next input matches the token then add auth commands
				for k, v := range auth.GenerateHandlers() {
					handlers[k] = v
				}

				// set this connection to "authenticated" so it wont need to do
				authenticated = true
			}
		}

		// look for the command
		handler, found := handlers[msg.Command]

		// if the command isn't found, return an error and wait for the next command
		if !found {
			encoder.Encode(&mist.Message{Command: msg.Command, Tags: msg.Tags, Data: msg.Data, Error: "Unknown Command"})
			continue
		}

		// attempt to run the command; if the command fails return the error and wait
		// for the next command
		if err := handler(proxy, msg); err != nil {
			encoder.Encode(&mist.Message{Command: msg.Command, Error: err.Error()})
			continue
		}
	}
}
