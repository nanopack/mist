//
package commands

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/auth"
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
				if err := auth.Start(viper.GetString("authenticator")); err != nil {
					fmt.Println("Failed to start authenticator!", err)
					os.Exit(1)
				}

				//
				if err := server.Start(viper.GetStringSlice("listeners"), viper.GetString("token")); err != nil {
					fmt.Println("One or more servers failed to start!", err)
					os.Exit(1)
				}

				//
				// if err := replicator.Start(); err != nil {
				// 	os.Exit(1)
				// }

				//
				return
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
	viper.SetDefault("log-level", "INFO")
	viper.SetDefault("listeners", []string{"tcp://127.0.0.1:1445", "http://127.0.0.1:8080", "ws://127.0.0.1:8888"})
	viper.SetDefault("replicator", "")
	viper.SetDefault("token", "")

	// persistent flags; these are the only 2 options that we want overridable from
	// the CLI, all others need to use a config file
	MistCmd.PersistentFlags().String("authenticator", viper.GetString("authenticator"), "desc.")
	MistCmd.PersistentFlags().StringSlice("listeners", viper.GetStringSlice("listeners"), "desc.")
	MistCmd.PersistentFlags().String("log-level", viper.GetString("log-level"), "desc.")
	MistCmd.PersistentFlags().String("replicator", viper.GetString("replicator"), "desc.")
	MistCmd.PersistentFlags().String("token", viper.GetString("token"), "desc.")

	//
	viper.BindPFlag("listeners", MistCmd.PersistentFlags().Lookup("listeners"))
	viper.BindPFlag("authenticator", MistCmd.PersistentFlags().Lookup("authenticator"))
	viper.BindPFlag("log-level", MistCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("replicator", MistCmd.PersistentFlags().Lookup("replicator"))
	viper.BindPFlag("token", MistCmd.PersistentFlags().Lookup("token"))

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
