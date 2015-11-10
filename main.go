// -*- mode: go; tab-width: 2; indent-tabs-mode: 1; st-rulers: [70] -*-
// vim: ts=4 sw=4 ft=lua noet
//--------------------------------------------------------------------
// @author Daniel Barney <daniel@nanobox.io>
// Copyright (C) Pagoda Box, Inc - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly
// prohibited. Proprietary and confidential
//
// @doc
//
// @end
// Created :   12 August 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
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
