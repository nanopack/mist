package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	mist "github.com/nanopack/mist/core"

	"github.com/gorilla/mux"
	"github.com/jcelliott/lumber"
)

var (
	// Router ...
	Router = mux.NewRouter()
)

// Message the message we need to receive from http client
type Message struct {
	Tags []string `json:"tags,omitempty"`
	Data string   `json:"data,omitempty"`
}

// init adds http/https as available mist server types
func init() {
	Register("http", StartHTTP)
}

// StartHTTP starts a mist server listening over HTTP
func StartHTTP(uri string, errChan chan<- error) {
	if err := newHTTP(uri); err != nil {
		errChan <- fmt.Errorf("Unable to start mist http listener - %s", err.Error())
	}
}

func newHTTP(address string) error {
	lumber.Info("HTTP server listening at '%s'...\n", address)

	// blocking...
	return http.ListenAndServe(address, routes())
}

// routes registers all api routes with the router
func routes() *mux.Router {
	Router.HandleFunc("/ping", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("pong\n"))
	}).Methods("GET")
	Router.HandleFunc("/publish", handleRequest(publish)).Methods("POST")

	return Router
}

func publish(rw http.ResponseWriter, req *http.Request) {
	body, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		lumber.Error("read body failed '%s'...\n", readErr)
	}
	message := Message{}
	jsonErr := json.Unmarshal(body, &message)
	if jsonErr != nil {
		lumber.Error("publish message error '%s'...\n", jsonErr)
	}
	mist.Publish(message.Tags, message.Data)
}

// handleRequest is a wrapper for the actual route handler, simply to provide some
// debug output
func handleRequest(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		fn(rw, req)
	}
}
