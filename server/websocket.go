package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/jcelliott/lumber"

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
	lumber.Info("WS server listening at '%s'...\n", uri)

	router := pat.New()
	router.Get("/subscribe/websocket", func(rw http.ResponseWriter, req *http.Request) {

		//
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		//
		conn, err := upgrader.Upgrade(rw, req, nil)
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

				// failing to write is probably because the connection is dead; we dont
				// want mist just looping forever tyring to write to something it will
				// never be able to.
				if err := conn.WriteJSON(msg); err != nil {
					errChan <- fmt.Errorf(err.Error())
					break
				}
			}
		}()

		// if an authenticator was passed, check for a token on connect to see if
		// auth commands are added
		if auth.DefaultAuth != nil && !authenticated {

			//
			var xtoken string
			switch {
			case req.Header.Get("X-AUTH-TOKEN") != "":
				xtoken = req.Header.Get("X-AUTH-TOKEN")
			case req.FormValue("x-auth-token") != "":
				xtoken = req.FormValue("x-auth-token")
			}

			// if the next input matches the token then add auth commands
			if xtoken == authtoken {
				// break // allow connection w/o admin commands
				return // disconnect client
			}

			// add auth commands ("admin" mode)
			for k, v := range auth.GenerateHandlers() {
				handlers[k] = v
			}

			// establish that the socket has already authenticated
			authenticated = true
		}

		// connection loop (blocking); continually read off the connection. Once something
		// is read, check to see if it's a message the client understands to be one of
		// its commands. If so attempt to execute the command.
		for {

			msg := mist.Message{}

			// failing to read is probably because the connection is dead; we dont
			// want mist just looping forever tyring to write to something it will
			// never be able to.
			if err := conn.ReadJSON(&msg); err != nil {
				errChan <- fmt.Errorf(err.Error())
				break
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
