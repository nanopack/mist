//
package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-discovery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/api"
	"github.com/nanopack/mist/authenticate"
	"github.com/nanopack/mist/core"
	"github.com/nanopack/mist/handlers"
)

var (
	log lumber.Logger

	//
	config  string //
	server  bool   //
	version bool   //

	//
	MistCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,

		// parse the config if one is provided, or use the defaults. Set the backend
		// driver to be used
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// create a new logger
			log = lumber.NewConsoleLogger(lumber.LvlInt(viper.GetString("log-level")))
			log.Prefix("[mist]")

			// if --config is passed, attempt to parse the config file
			if config != "" {

				//
				viper.SetConfigName("config")
				viper.AddConfigPath(config)

				// Find and read the config file; Handle errors reading the config file
				if err := viper.ReadInConfig(); err != nil {
					panic(fmt.Errorf("Fatal error config file: %s \n", err))
				}
			}
		},

		// either run mist as a server, or run it as a CLI depending on what flags
		// are provided
		Run: func(ccmd *cobra.Command, args []string) {

			// if --server is passed start the mist server; Assuming an http server for
			// the time being. At some point this may be configurable
			if server {

				//
				mist := mist.New()

				//
				if viper.GetString("multicast-interface") != "single" {

					// start discovering other mist nodes on the network
					discover, err := discovery.NewDiscovery(viper.GetString("multicast-interface"), "mist", time.Second*2)
					if err != nil {
						panic(err)
					}
					defer discover.Close()

					// advertise this nodes listen address
					discover.Add("mist", viper.GetString("tcp-addr"))

					// enable replication between mist nodes
					replicate := handlers.EnableReplication(mist, discover)
					fmt.Println(fmt.Sprintf("Starting Mist monitor... \nTCP address: %s\nHTTP address: %s", viper.GetString("tcp-addr"), viper.GetString("http-addr")))
					go replicate.Monitor()
				}

				//
				pgAuth, err := authenticate.NewPostgresqlAuthenticator(viper.GetString("db-user"), viper.GetString("db-name"), viper.GetString("db-addr"))
				if err != nil {
					log.Fatal("Unable to start postgresql authenticator ", err)
					os.Exit(1)
				}

				// start a mist server listening over TCP; this is a non-blocking server
				// because we also want to start a web server and will leave the blocking
				// up to it.
				log.Info("Starting mist server (TCP) at '%s'...\n", viper.GetString("tcp-addr"))
				server, err := mist.Listen(viper.GetString("tcp-addr"), handlers.GenerateAdditionalCommands(pgAuth))
				if err != nil {
					log.Fatal("Unable to start mist tcp listener ", err)
					os.Exit(1)
				}
				defer server.Close()

				// start a mist server listening over HTTP (blocking)
				log.Info("Starting mist server (HTTP) at '%s'...\n", viper.GetString("http-addr"))
				if err := api.Start(); err != nil {
					log.Fatal("Failed to start - ", err.Error())
					os.Exit(1)
				}
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {

	// local flags;
	MistCmd.Flags().StringVarP(&config, "config", "c", "", "Path to config options")
	MistCmd.Flags().BoolVarP(&server, "server", "s", false, "Run mist as a server")
	MistCmd.Flags().BoolVarP(&version, "version", "v", false, "Display the current version of this CLI")

	// set config defaults; these are overriden if a --config file is provided
	// (see above)
	viper.SetDefault("tcp-addr", "127.0.0.1:1445")
	viper.SetDefault("http-addr", "127.0.0.1:8080")
	viper.SetDefault("log-level", "INFO")
	viper.SetDefault("multicast-interface", "single")
	viper.SetDefault("db-user", "postgres")
	viper.SetDefault("db-name", "postgres")
	viper.SetDefault("db-addr", "127.0.0.1:5432")

	// persistent flags; these are the only 2 options that we want overridable from
	// the CLI, all others need to use a config file
	MistCmd.PersistentFlags().String("tcp-addr", viper.GetString("tcp-addr"), "desc.")
	viper.BindPFlag("tcp-addr", MistCmd.PersistentFlags().Lookup("tcp-addr"))

	MistCmd.PersistentFlags().String("log-level", viper.GetString("log-level"), "desc.")
	viper.BindPFlag("log-level", MistCmd.PersistentFlags().Lookup("log-level"))

	// commands
	MistCmd.AddCommand(listCmd)
	MistCmd.AddCommand(pingCmd)
	MistCmd.AddCommand(publishCmd)
	MistCmd.AddCommand(subscribeCmd)
	MistCmd.AddCommand(unsubscribeCmd)

	// hidden/aliased commands
	MistCmd.AddCommand(messageCmd)
	MistCmd.AddCommand(sendCmd)
}
