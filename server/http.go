package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"

	"github.com/nanopack/mist/auth"
)

var (
	Router = pat.New()
)

// start a mist server listening over HTTP
func startHTTP(uri string, errChan chan<- error)  {
	if err := ListenHTTP(uri); err != nil {
		errChan<- fmt.Errorf("Unable to start mist http listener %v", err)
	}
}

//
func ListenHTTP(address string) error {
	fmt.Printf("HTTP server listening at '%s'...\n", address)

	// blocking...
	return http.ListenAndServe(address, routes())
}

// routes registers all api routes with the router
func routes() *pat.Router {
	// config.Log.Debug("Registering routes...\n")

	//
	Router.Get("/ping", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("pong\n"))
	})
	// Router.Get("/list", handleRequest(list))
	// Router.Get("/subscribe", handleRequest(subscribe))
	// Router.Get("/unsubscribe", handleRequest(unsubscribe))

	// start up the authenticated websocket connection
	Router.Get("/subscribe/websocket", auth.AuthenticateWebsocket())

	return Router
}

// handleRequest is a wrapper for the actual route handler, simply to provide some
// debug output
func handleRequest(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		fn(rw, req)

		// must be after fn if ever going to get rw.status (logging still more meaningful)
		// config.Log.Trace(`%v - [%v] %v %v %v(%s) - "User-Agent: %s", "X-Nanobox-Token: %s"`,
		// 	req.RemoteAddr, req.Proto, req.Method, req.RequestURI,
		// 	rw.Header().Get("status"), req.Header.Get("Content-Length"),
		// 	req.Header.Get("User-Agent"), req.Header.Get("X-Nanobox-Token"))
	}
}
