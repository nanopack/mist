package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
)

// init adds ws/wss as available mist server types
func init() {
	Register("ws", StartWS)
	Register("wss", StartWSS)
}

// StartWS starts a mist server listening over a websocket
func StartWS(uri string, errChan chan<- error) {
	fmt.Printf("WS server listening at '%s'...\n", uri)

	router := pat.New()
	router.Get("/subscribe/websocket", func(w http.ResponseWriter, r *http.Request) {

		//
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		//
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			errChan <- fmt.Errorf("Failed to upgrade connection %v", err.Error())
			return
		}
		defer conn.Close()

		//
		proxy := mist.NewProxy()
		defer proxy.Close()

		// add basic WS handlers for this socket
		handlers = GenerateHandlers()

		// read and publish mist messages to connected clients (non-blocking)
		go func() {
			for msg := range proxy.Pipe {
				if err := conn.WriteJSON(msg); err != nil {
					errChan <- fmt.Errorf(err.Error())
					continue
				}
			}
		}()

		// check for authentication
		switch {

		// authentication wanted...
		case auth.DefaultAuth != nil:

			//
			var xtoken string
			switch {
			case r.Header.Get("x-auth-token") != "":
				xtoken = r.Header.Get("x-auth-token")
			case r.FormValue("x-auth-token") != "":
				xtoken = r.FormValue("x-auth-token")
			}

			// if the websocket is connected with the required token, add auth command
			// handlers
			if xtoken == token {
				for k, v := range auth.GenerateHandlers() {
					handlers[k] = v
				}
			}

		// no authentication wanted; authorize the proxy
		default:
			// proxy.Authorized = true
		}

		// connection loop (blocking); continually read off the connection. Once something
		// is read, check to see if it's a message the client understands to be one of
		// its commands. If so attempt to execute the command.
		for {

			msg := mist.Message{}

			// decode an array value (Message)
			if err := conn.ReadJSON(&msg); err != nil {
				errChan <- fmt.Errorf(err.Error())
				continue
			}

			// look for the command
			handler, found := handlers[msg.Command]

			// if the command isn't found, return an error
			if !found {
				if err := conn.WriteJSON(&mist.Message{Command: msg.Command, Error: "Unknown Command"}); err != nil {
					errChan <- fmt.Errorf(err.Error())
				}
				continue
			}

			// attempt to run the command
			if err := handler(proxy, msg); err != nil {
				if err := conn.WriteJSON(&mist.Message{Command: msg.Command, Error: err.Error()}); err != nil {
					errChan <- fmt.Errorf(err.Error())
				}
				continue
			}
		}
	})

	//
	go http.ListenAndServe(uri, router)
}

// StartWSS starts a mist server listening over a secure websocket
func StartWSS(uri string, errChan chan<- error) {
	errChan <- ErrNotImplemented
}
