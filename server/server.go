package server

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanopack/mist/auth"
)

var (
	// ErrNotImplemented not implemented error
	ErrNotImplemented = fmt.Errorf("error: not implemented")

	// this is a map of the supported servers that can be started by mist
	servers    = map[string]handleFunc{}
	serversTex sync.RWMutex

	authtoken string // used when determining if auth command handlers should be added
)

type (
	handleFunc func(uri string, errChan chan<- error)
)

// Register registers a new mist server
func Register(name string, auth handleFunc) {
	serversTex.Lock()
	servers[name] = auth
	serversTex.Unlock()
}

// Start attempts to individually start mist servers from a list of provided
// listeners; the listeners provided is a comma delimited list of uri strings
// (scheme:[//[user:pass@]host[:port]][/]path[?query][#fragment])
func Start(uris []string, token string) error {
	// check to see if a token is provided; an authenticator cannot work without
	// a token and so it should error here informing that.
	if auth.IsConfigured() && token == "" {
		return fmt.Errorf("an authenticator has been specified but no token provided")
	}

	// set the authtoken
	authtoken = token

	// this chan is given to each individual server start as a way for them to
	// communcate back their startup status
	errChan := make(chan error, len(uris))

	// iterate over each of the provided listener uris attempting to start them
	// individually; if one isn't supported it gets skipped
	for i := range uris {

		// parse the uri string into a url object
		url, err := url.Parse(uris[i])
		if err != nil {
			return err
		}

		// check to see if the scheme is supported; if not, indicate as such and
		// continue
		serversTex.RLock()
		server, ok := servers[url.Scheme]
		serversTex.RUnlock()
		if !ok {
			lumber.Error("Unsupported scheme '%s'", url.Scheme)
			continue
		}

		// attempt to start the server
		lumber.Info("Starting '%s' server...", url.Scheme)
		go server(url.Host, errChan)
	}

	// handle errors that happen during startup by reading off errChan and returning
	// on any error received. If no errors are received after 1 second per server
	// assume successful starts.
	select {
	case err := <-errChan:
		lumber.Error("Failed to start - %s", err.Error())
		return err
	case <-time.After(time.Second * time.Duration(len(uris))):
		// no errors
	}

	// handle errors that happen after initial start; if any errors are received they
	// are logged and the servers just try to keep running
	for err := range errChan {
		// log these errors and continue
		lumber.Error("Server error - %s", err.Error())
	}

	return nil
}
