//
package commands

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/server"
)

var (
	log lumber.Logger

	//
	config  string //
	daemon  bool   //
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
			if daemon {

				//
				if viper.GetString("multicast-interface") != "single" {
					server.EnableDiscovery()
					server.EnableReplication()
				}

				//
				server.Start()
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {

	// local flags;
	MistCmd.Flags().StringVarP(&config, "config", "", "", "Path to config options")
	MistCmd.Flags().BoolVarP(&daemon, "server", "", false, "Run mist as a server")
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
