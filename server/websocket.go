package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/server/handlers"
	"github.com/nanopack/mist/util"
)

//
var wsCommands map[string]mist.WSHandler

//
func init() {

	// add WS handlers
	wsCommands = handlers.GenerateWSCommands()
}

//
func ListenWS(mixins map[string]mist.WSHandler) http.HandlerFunc {

	fmt.Println("LISTEN WS!")

	//
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("HERE????!?!?!?!")

		//
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		fmt.Println("UPGRADER??", upgrader)

		//
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("ERR???", err)
			return
		}

		fmt.Println("UPGRDAE!", conn)

		// add any additional handlers
		for k, v := range mixins {
			wsCommands[k] = v
		}

		fmt.Println("COMMANDS?", wsCommands)

		// we don't want this to be buffered
		client, err := mist.NewProxy(0)
		if err != nil {
			fmt.Println("BIONK!", err)
		}

		write := make(chan string)
		done := make(chan bool)
		defer func() {
			client.Close()
			close(done)
		}()

		// the gorilla websocket package must have all writes come from the
		// same process.
		go func() {

			fmt.Println("HERE????")

			for {
				select {
				case event := <-client.Messages():
					if msg, err := json.Marshal(event); err == nil {
						conn.WriteMessage(websocket.TextMessage, msg)
					}
				case msg := <-write:
					conn.WriteMessage(websocket.TextMessage, []byte(msg))
				case <-done:
					close(write)
					return
				}
			}
		}()

		for {
			msgType, frame, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("BZONK!", err)
			}

			//
			if msgType != websocket.TextMessage {
				write <- "{\"success\":false,\"error\":\"I don't understand binary messages\"}"
				continue
			}

			input := struct {
				Cmd string `json:"command"`
			}{}

			//
			if err := json.Unmarshal(frame, &input); err != nil {
				write <- "{\"success\":false,\"error\":\"Invalid json\"}"
				continue
			}

			//
			handler, found := wsCommands[input.Cmd]

			//
			if !found {
				write <- "{\"success\":false,\"error\":\"unknown command\"}"
				continue
			}

			// something needs to happen with this error..
			if err := handler.Handle(client, frame, write); err != nil {
				fmt.Println("BROZONK!", err)
				return
			}
		}
	}
}

//
func AuthenticateWebsocket(a auth.Authenticator) http.HandlerFunc {

	fmt.Println("AUTH WEBSOCKET???", a)

	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("HERE!!?!?!", r.FormValue("x-auth-token"), r.Header.Get("x-auth-token"))

		//
		var token string
		switch {
		case r.Header.Get("x-auth-token") != "":
			token = r.Header.Get("x-auth-token")
		case r.FormValue("x-auth-token") != "":
			token = r.FormValue("x-auth-token")
		default:
			token = "unauthorized"
		}

		fmt.Println("TOKEN??", token)

		// if they have no tags registered for the token, then they are not authorized
		// to connect to mist
		if tags, err := a.GetTagsForToken(token); err != nil || len(tags) == 0 {
			fmt.Println("BRONK??", err)
			w.WriteHeader(401)
			return
		}

		fmt.Println("HERE!")

		// overwrite the subscribe command so that we can add authentication to it.
		mixins := map[string]mist.WSHandler{
			"subscribe": {0, handleWSAuthSubscribe(token, a)},
		}

		//
		wsUpgrade := ListenWS(mixins)
		wsUpgrade(w, r)
	}
}

//
func handleWSAuthSubscribe(token string, a auth.Authenticator) func(client mist.Client, frame []byte, write chan<- string) error {
	return func(client mist.Client, frame []byte, write chan<- string) error {

		authTags, err := a.GetTagsForToken(token)
		if err != nil || len(authTags) == 0 {
			write <- "{\"success\":false,\"command\":\"subscribe\"}"
			return nil
		}

		//
		tags := struct {
			Tags []string `json:"tags"`
		}{}

		// error would already be caught by unmarshalling the command
		if err := json.Unmarshal(frame, &tags); err != nil {
			fmt.Println("BUNGK!")
		}

		if !util.HaveSameTags(authTags, tags.Tags) {
			write <- "{\"success\":false,\"command\":\"subscribe\"}"
			return nil
		}
		client.Subscribe(tags.Tags)
		write <- "{\"success\":true,\"command\":\"subscribe\"}"
		return nil
	}
}
