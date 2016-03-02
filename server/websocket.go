package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/util"
)

//
var wsCommands map[string]Handler

//
func init() {

	// add WS handlers
	wsCommands = GenerateHandlers()
}

// start a mist server listening over HTTP
// func startWS(uri string, errChan chan<- error)  {
// 	if err := ListenWS(uri); err != nil {
// 		errChan<- fmt.Errorf("Unable to start mist http listener %v", err)
// 	}
// }

//
func ListenWS(mixins map[string]Handler) http.HandlerFunc {

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
		proxy := mist.NewProxy(0)

		write := make(chan string)
		done := make(chan bool)
		defer func() {
			proxy.Close()
			close(done)
		}()

		// the gorilla websocket package must have all writes come from the
		// same process.
		go func() {

			fmt.Println("HERE????")

			for {
				select {
				case event := <-proxy.Messages():
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
				Args []string `json:"args"`
			}{}

			//
			if err := json.Unmarshal(frame, &input); err != nil {
				write <- "{\"success\":false,\"error\":\"Invalid json\"}"
				continue
			}

			//
			handler, found := wsCommands[input.Cmd]

			//
			var response string
			switch {

			// no command found
			case !found:
				write <- fmt.Sprintf("{\"success\":false,\"error\":\"Unknown command '%s'\"}", input.Cmd)
				// continue

			//
		case handler.NumArgs != len(input.Args):
				write <- fmt.Sprintf("{\"success\":false,\"error\":\"Wrong number of args for '%s'. Expected %v got %v\"}", input.Cmd, handler.NumArgs, len(input.Args))

			// execute command
			default:
				fmt.Println("WSS EXECUTE! ", input.Cmd)
				response = handler.Handle(proxy, input.Args)
			}

			// write the response from the command back to the connection
			fmt.Println("WSS WRITING RESPONSE! ", response)
			write <- response
		}
	}
}

//
// func handleWSPing(client mist.Proxy, frame []byte, write chan<- string) error {
// 	write <- "{\"success\":true,\"command\":\"ping\"}"
// 	return nil
// }
//
// //
// func handleWSSubscribe(client mist.Proxy, frame []byte, write chan<- string) error {
// 	tags := struct {
// 		Tags []string `json:"tags"`
// 	}{}
//
// 	// error would already be caught by unmarshalling the command
// 	if err := json.Unmarshal(frame, &tags); err != nil {
// 		fmt.Println("BUNK!", err)
// 	}
//
// 	//
// 	client.Subscribe(tags.Tags)
//
// 	write <- "{\"success\":true,\"command\":\"subscribe\"}"
//
// 	return nil
// }
//
// //
// func handleWSUnubscribe(client mist.Proxy, frame []byte, write chan<- string) error {
// 	tags := struct {
// 		Tags []string `json:"tags"`
// 	}{}
//
// 	// error would already be caught by unmarshalling the command
// 	if err := json.Unmarshal(frame, &tags); err != nil {
// 		fmt.Println("BUNK!", err)
// 	}
//
// 	//
// 	client.Unsubscribe(tags.Tags)
//
// 	write <- "{\"success\":true,\"command\":\"unsubscribe\"}"
//
// 	return nil
// }
//
// //
// func handleWSList(client mist.Proxy, frame []byte, write chan<- string) (err error) {
//
// 	//
// 	list := struct {
// 		Subscriptions [][]string `json:"subscriptions"`
// 		Command       string     `json:"command"`
// 		Success       bool       `json:"success"`
// 	}{}
//
// 	if list.Subscriptions, err = client.List(); err != nil {
// 		return err
// 	}
//
// 	list.Command = "list"
// 	list.Success = true
//
// 	bytes, err := json.Marshal(list)
// 	if err != nil {
// 		return err
// 	}
//
// 	//
// 	write <- string(bytes)
//
// 	return
// }

//
func handleWSAuthSubscribe(token string, a auth.Authenticator) func(proxy mist.Proxy, frame []byte, write chan<- string) error {
	return func(proxy mist.Proxy, frame []byte, write chan<- string) error {

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
		proxy.Subscribe(tags.Tags)
		write <- "{\"success\":true,\"command\":\"subscribe\"}"
		return nil
	}
}
