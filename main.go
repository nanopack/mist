// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//
package main

import (
	"bitbucket.org/nanobox/na-api"
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-discovery"
	"github.com/nanobox-io/golang-mist/authenticate"
	"github.com/nanobox-io/golang-mist/core"
	"github.com/nanobox-io/golang-mist/handlers"
	"github.com/nanobox-io/nanobox-config"
	"os"
	"strings"
	"time"
)

func main() {
	configFile := ""
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		configFile = os.Args[1]
	}

	defaults := map[string]string{
		"tcp_listen_address":  "127.0.0.1:1445",
		"http_listen_address": "127.0.0.1:8080",
		"log_level":           "INFO",
		"multicast_interface": "eth1",
	}

	config.Load(defaults, configFile)

	level := lumber.LvlInt(config.Config["log_level"])

	mist := mist.New()
	api.Name = "MIST"
	api.Logger = lumber.NewConsoleLogger(level)
	api.User = mist

	listen := config.Config["tcp_listen_address"]
	server, err := mist.Listen(listen, nil)

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
	go replicate.Monitor()

	// start up the authenticated websocket connection
	authenticator := authenticate.NewNoopAuthenticator()
	handlers.LoadWebsocketRoute(authenticator)
	api.Start(config.Config["http_listen_address"])
}
