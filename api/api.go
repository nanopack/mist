package api

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/spf13/viper"
)

// Start the api
func Start() error {

	// blocking...
	return http.ListenAndServe(viper.GetString("HTTPAddr"), routes())
}

// routes registers all api routes with the router
func routes() *pat.Router {
	// config.Log.Debug("Registering routes...\n")

	//
	router := pat.New()

	//
	router.Get("/ping", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("pong\n"))
	})

	// blobs
	// router.Get("/list", handleRequest(list))
	// router.Get("/subscribe", handleRequest(subscribe))
	// router.Get("/unsubscribe", handleRequest(unsubscribe))

	return router
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
