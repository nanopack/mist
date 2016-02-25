package server

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/nanobox-io/golang-discovery"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/server/handlers"
)

//
func Start() {

	//
	pgAuth, err := auth.NewPostgresql(viper.GetString("db-user"), viper.GetString("db-name"), viper.GetString("db-addr"))
	if err != nil {
		fmt.Println("Unable to start postgresql authenticator ", err)
		os.Exit(1)
	}

	//
	Auth = pgAuth

	// start a mist server listening over TCP; this is a non-blocking server
	// because we also want to start a web server and will leave the blocking
	// up to it.
	server, err := ListenTCP(viper.GetString("tcp-addr"), handlers.GenerateAuthCommands(Auth))
	if err != nil {
		fmt.Println("Unable to start mist tcp listener ", err)
		os.Exit(1)
	}
	fmt.Println("SERVER!", server)
	// defer server.Close()

	// start a mist server listening over HTTP (blocking)
	if err := ListenHTTP(viper.GetString("http-addr")); err != nil {
		fmt.Println("Unable to start mist http listener ", err)
		os.Exit(1)
	}

	fmt.Println("STARTED!")
}

//
func ConfigureAsMultinode(mist *mist.Mist) {

	// start discovering other mist nodes on the network
	discover, err := discovery.NewDiscovery(viper.GetString("multicast-interface"), "mist", time.Second*2)
	if err != nil {
		panic(err)
	}
	defer discover.Close()

	// advertise this nodes listen address
	discover.Add("mist", viper.GetString("tcp-addr"))

	// enable replication between mist nodes
	// replicate := handlers.EnableReplication(mist, discover)
	// fmt.Println(fmt.Sprintf("Starting Mist monitor... \nTCP address: %s\nHTTP address: %s", viper.GetString("tcp-addr"), viper.GetString("http-addr")))

	// go replicate.Monitor()
}
