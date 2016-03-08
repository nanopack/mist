package server

import (
	"fmt"
	"net/url"
	"time"

	"github.com/nanopack/mist/core"
)

//
var (

	//
	ErrNotImplemented = fmt.Errorf("Error: Not Implemented\n")

	// this is a map of the supported listeners/servers that can be started by mist
	listeners = map[string]func(uri string, errChan chan<- error){}
	handlers  = map[string]mist.Handler{}

	//
	authorized = true
)

// Start attempts to individually start mist servers from a list of provided
// listeners; the listeners provided is a comma delimited list of uri strings
// (scheme:[//[user:pass@]host[:port]][/]path[?query][#fragment])
func Start(uris []string) error {

	// this chan is given to each individual server start as a way for them to
	// communcate back their startup status
	errChan := make(chan error, len(uris))

	// iterate over each of the provided listener uris attempting to start them
	// individually; if one isn't supported it skipped
	for _, uri := range uris {

		// parse the uri string into a url object
		url, err := url.Parse(uri)
		if err != nil {
			return err
		}

		// check to see if the scheme is supported; if not, indicate as such and
		// continue
		server, ok := listeners[url.Scheme]
		if !ok {
			fmt.Printf("Unsupported scheme '%v'", url.Scheme)
			continue
		}

		// attempt to start the server
		go server(url.Host, errChan)
	}

	// handle errors that happen during startup by reading off errChan and returning
	// on any error received. If no errors are received after 1 second per server
	// assume successful starts.
	select {
	case err := <-errChan:
		return err
	case <-time.After(time.Second * time.Duration(len(uris))):
		// no errors
	}

	// handle errors that happen after initial start; if any errors are received they
	// are logged and the servers try to just keep running
	// go func() {
	// 	for err := range errChan {
	// 		fmt.Println("ERR!", err)
	// 		// write to a log
	// 	}
	// }()

	// REMOVE THIS WHEN DONE: we'll just hold the connection open for now to see output
	for err := range errChan {
		fmt.Println("ERR!", err)
		// write to a log
	}

	return nil
}
