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

	// "github.com/nanopack/mist/api"
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
			log = lumber.NewConsoleLogger(lumber.LvlInt(viper.GetString("LogLevel")))
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
			if server != false {
				log.Info("Starting mist server at '%s'...\n", viper.GetString("TCPAddr"))

				//
				mist := mist.New()

				//
				if viper.GetString("MulticastInterface") != "single" {

					// start discovering other mist nodes on the network
					discover, err := discovery.NewDiscovery(viper.GetString("MulticastInterface"), "mist", time.Second*2)
					if err != nil {
						panic(err)
					}
					defer discover.Close()

					// advertise this nodes listen address
					discover.Add("mist", viper.GetString("TCPAddr"))

					// enable replication between mist nodes
					replicate := handlers.EnableReplication(mist, discover)
					fmt.Println(fmt.Sprintf("Starting Mist monitor... \nTCP address: %s\nHTTP address: %s", viper.GetString("TCPAddr"), viper.GetString("HTTPAddr")))
					go replicate.Monitor()
				}

				//
				pgAuth, err := authenticate.NewPostgresqlAuthenticator(viper.GetString("DBUser"), viper.GetString("DBName"), viper.GetString("DBAddr"))
				if err != nil {
					log.Fatal("Unable to start postgresql authenticator %v", err)
					os.Exit(1)
				}

				// start up the authenticated websocket connection
				authenticator := authenticate.NewNoopAuthenticator()
				handlers.LoadWebsocketRoute(authenticator)

				// start a mist server... (Blocking)
				server, err := mist.Listen(viper.GetString("TCPAddr"), handlers.GenerateAdditionalCommands(pgAuth))
				if err != nil {
					log.Fatal("Unable to start mist tcp listener %v", err)
					os.Exit(1)
				}
				defer server.Close()

				// start the API
				// if err := api.Start(); err != nil {
				// 	log.Fatal("Failed to start - %s", err.Error())
				// 	os.Exit(1)
				// }
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {

	//
	viper.SetDefault("TCPAddr", "127.0.0.1:1445")
	viper.SetDefault("HTTPAddr", "127.0.0.1:8080")
	viper.SetDefault("LogLevel", "INFO")
	viper.SetDefault("MulticastInterface", "single")
	viper.SetDefault("DBUser", "postgres")
	viper.SetDefault("DBName", "postgres")
	viper.SetDefault("DBAddr", "127.0.0.1:5432")

	// persistent flags
	MistCmd.PersistentFlags().String("tcp-addr", viper.GetString("TCPAddr"), "desc.")
	MistCmd.PersistentFlags().String("http-addr", viper.GetString("HTTPAddr"), "desc.")
	MistCmd.PersistentFlags().String("log-level", viper.GetString("LogLevel"), "desc.")
	MistCmd.PersistentFlags().String("multicast-interface", viper.GetString("MulticastInterface"), "desc.")
	MistCmd.PersistentFlags().String("db-user", viper.GetString("DBUser"), "desc.")
	MistCmd.PersistentFlags().String("db-name", viper.GetString("DBName"), "desc.")
	MistCmd.PersistentFlags().String("db-addr", viper.GetString("DBAddr"), "desc.")

	// local flags
	MistCmd.Flags().StringVarP(&config, "config", "c", "", "Path to config options")
	MistCmd.Flags().BoolVarP(&server, "server", "s", false, "Run mist as a server")
	MistCmd.Flags().BoolVarP(&version, "version", "v", false, "Display the current version of this CLI")

	//
	viper.BindPFlag("tcp-addr", MistCmd.Flags().Lookup("tcp-addr"))
	viper.BindPFlag("http-addr", MistCmd.Flags().Lookup("http-addr"))
	viper.BindPFlag("log-level", MistCmd.Flags().Lookup("log-level"))
	viper.BindPFlag("multicast-interface", MistCmd.Flags().Lookup("multicast-interface"))
	viper.BindPFlag("db-user", MistCmd.Flags().Lookup("db-user"))
	viper.BindPFlag("db-name", MistCmd.Flags().Lookup("db-name"))
	viper.BindPFlag("db-addr", MistCmd.Flags().Lookup("db-addr"))

	// commands
	MistCmd.AddCommand(sendCmd)
	MistCmd.AddCommand(subscribeCmd)
	MistCmd.AddCommand(unsubscribeCmd)

	// hidden/aliased commands
	MistCmd.AddCommand(messageCmd)
}
