package server

import (
	"fmt"
	"net/url"
	"time"
)

//
var (

	//
	listeners = map[string]func(uri string, errChan chan<- error) {
		"tcp": startTCP,
		"http": startHTTP,
		// "https": startHTTPS,
		// "ws": startWS,
		// "wss": startWSS,
	}
)

//
func Start(schemas []string) error {

	fmt.Println("SERVER START!")

	// NOTE look at viper string slice as an alternative to this
	errChan := make(chan error, len(schemas))
	for _, schema := range schemas {

		thing, err := url.Parse(schema)
		if err != nil {
			return err
		}

		server, ok := listeners[thing.Scheme]

		//
		if !ok {
			fmt.Errorf("Unsuported schema %v", thing.Scheme)
			continue
		}

		// start the server
		go server(thing.Host, errChan)
	}

	// handle errors that happen during startup
	select {
	case err := <-errChan:
		return err
	case <-time.After(time.Second*5):
		fmt.Println("NO ERRORS!")
		// no errors
	}

	// handle errors that happen after initial start
	// go func() {
	// 	for err := range errChan {
	// 		fmt.Println("ERR!", err)
	// 		// write to a log
	// 	}
	// }()

	// we'll just hold the connection open for now to see output
	for err := range errChan {
		fmt.Println("ERR!", err)
		// write to a log
	}

	fmt.Println("DONE!!!!!!!!!")

	return nil
}

// EnableDiscovery starts discovering other mist nodes on the network
// func EnableDiscovery() error {
//
// 	discover, err := discovery.NewDiscovery(viper.GetString("multicast-interface"), "mist", time.Second*2)
// 	if err != nil {
// 		return err
// 	}
// 	defer discover.Close()
//
// 	// advertise this nodes listen address
// 	discover.Add("mist", viper.GetString("tcp-addr"))
// }

// EnableReplication enables replication between mist nodes
// func EnableReplication() {
//
// 	mist := mist.New()
//
// 	replicate := handlers.EnableReplication(mist, discover)
// 	fmt.Println(fmt.Sprintf("Starting Mist monitor... \nTCP address: %s\nHTTP address: %s", viper.GetString("tcp-addr"), viper.GetString("http-addr")))
//
// 	go replicate.Monitor()
// }
