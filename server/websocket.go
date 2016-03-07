package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
)

// init
func init() {

	// add websockets as an available server type
	listeners["ws"] = StartWS
	listeners["wss"] = StartWSS
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
			// CheckOrigin: func(r *http.Request) bool {
			// 	return true
			// },
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

		// read and publish mist messages to connected clients (non-blocking)
		go func() {

			// convert mist messages to text messages
			for msg := range proxy.Pipe {
				b, err := json.Marshal(msg)
				if err != nil {
					// log this error and continue?
				}

				//
				conn.WriteMessage(websocket.TextMessage, b)
			}
		}()

		//
		write := make(chan string)
		defer close(write)

		// read and publish websocket messages to connected clients (non-blocking)
		go func() {
			for msg := range write {
				conn.WriteMessage(websocket.TextMessage, []byte(msg))
			}
		}()

		// add basic WS handlers for this socket
		handlers = GenerateHandlers()

		// check for authentication
		switch {

		// authentication wanted...
		case auth.DefaultAuth != nil:
			//
			var token string
			switch {
			case r.Header.Get("x-auth-token") != "":
				token = r.Header.Get("x-auth-token")
			case r.FormValue("x-auth-token") != "":
				token = r.FormValue("x-auth-token")
			}

			// if the websocket is connected with the required token, add auth command
			// handlers
			if token == auth.Token {
				for k, v := range auth.GenerateHandlers() {
					handlers[k] = v
				}
			}

		// no authentication wanted; authorize the proxy
		default:
			// proxy.Authorized = true
		}

		// add a reader that reads off the connection (blocking)
		for {

			//
			msgType, frame, err := conn.ReadMessage()
			if err != nil {
				write <- fmt.Sprintf("{\"success\":false,\"error\":\"%v\"}", err.Error())
				// maybe log this also?
			}

			//
			if msgType != websocket.TextMessage {
				write <- "{\"success\":false,\"error\":\"I don't understand binary messages\"}"
				continue
			}

			fmt.Printf("FRAME! %#v\n", string(frame))

			//
			input := struct {
				Cmd  string   `json:"command"`
				Args []string `json:"tags"`
			}{}

			//
			if err := json.Unmarshal(frame, &input); err != nil {
				write <- fmt.Sprintf("{\"success\":false,\"error\":\"%v\"}", err.Error())
				continue
			}

			fmt.Printf("INPUT! %#v\n", input)

			//
			handler, found := handlers[input.Cmd]

			//
			// var err error
			switch {

			// command not found
			case !found:
				err = fmt.Errorf("Unknown command")

			// wrong number of arguments
			case handler.NumArgs != len(input.Args):
				err = fmt.Errorf("Wrong number of args. Expected %v got %v\"", handler.NumArgs, len(input.Args))

			// execute command
			default:
				fmt.Println("WSS EXECUTE! ", input.Cmd)
				err = handler.Handle(proxy, input.Args)

			}

			// if something failed along the way, respond accordingly...
			if err != nil {
				fmt.Println("FAIL!")
				write <- fmt.Sprintf("{\"success\":false, \"command\":\"%v\", \"error\":\"%v\"}", input.Cmd, err.Error())

				// break
				continue
			}

			// ...otherwise write a successful response
			fmt.Println("WSS WRITING RESPONSE! ")
			write <- fmt.Sprintf("{\"success\":true, \"command\":\"%v\", \"tags\":\"%v\"}", input.Cmd, input.Args)
		}
	})

	//
	go http.ListenAndServe(uri, router)
}

// StartWSS starts a mist server listening over a secure websocket
func StartWSS(uri string, errChan chan<- error) {
	errChan <- ErrNotImplemented
}
