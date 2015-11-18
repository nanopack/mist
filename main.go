// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package main

import (
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-discovery"
	"github.com/nanobox-io/nanobox-api"
	"github.com/nanobox-io/nanobox-config"
	"github.com/nanopack/mist/authenticate"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/handlers"
	"os"
	"time"
	"flag"
	"fmt"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.Parse()

	defaults := map[string]string{
		"tcp_listen_address":  "127.0.0.1:1445",
		"http_listen_address": "127.0.0.1:8080",
		"log_level":           "INFO",
		"multicast_interface": "eth1",
		"pg_user":             "postgres",
		"pg_database":         "postgres",
		"pg_address":          "127.0.0.1:5432",
	}

	config.Load(defaults, configFile)

	level := lumber.LvlInt(config.Config["log_level"])

	mist := mist.New()
	api.Name = "MIST"
	api.Logger = lumber.NewConsoleLogger(level)
	api.User = mist

	user := config.Config["pg_user"]
	database := config.Config["pg_database"]
	address := config.Config["pg_address"]

	pgAuth, err := authenticate.NewPostgresqlAuthenticator(user, database, address)
	if err != nil {
		api.Logger.Fatal("unable to start postgresql authenticator %v", err)
		os.Exit(1)
	}

	authCommands := handlers.GenerateAdditionalCommands(pgAuth)

	listen := config.Config["tcp_listen_address"]
	server, err := mist.Listen(listen, authCommands)

	if err != nil {
		api.Logger.Fatal("unable to start mist tcp listener %v", err)
		os.Exit(1)
	}
	defer server.Close()

	// start discovering other mist nodes on the network
	discover, err := discovery.NewDiscovery(config.Config["multicast_interface"], "mist", time.Second*2)
	if err != nil {
		panic(err)
	}
	defer discover.Close()

	// advertise this nodes listen address
	discover.Add("mist", listen)

	// enable replication between mist nodes
	replicate := handlers.EnableReplication(mist, discover)
	fmt.Println(fmt.Sprintf("Starting Mist monitor... \nTCP address: %s\nHTTP address: %s", config.Config["tcp_listen_address"],
		config.Config["http_listen_address"]))
	go replicate.Monitor()

	// start up the authenticated websocket connection
	authenticator := authenticate.NewNoopAuthenticator()
	handlers.LoadWebsocketRoute(authenticator)
	api.Start(config.Config["http_listen_address"])
}
