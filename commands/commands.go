// Package commands ...
package commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/server"
)

var (

	//
	config  string             //
	daemon  bool               //
	host    = "127.0.0.1:1445" //
	tags    []string           //
	version bool               //

	// MistCmd ...
	MistCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,

		// parse the config if one is provided, or use the defaults. Create a new logger
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// convert the log level
			logLvl := lumber.LvlInt(viper.GetString("log-level"))

			// configure the logger
			// lumber.Prefix("[hoader]")
			switch viper.GetString("log-type") {
			case "stdout":
				lumber.Level(logLvl)
			case "file":
				// logger := lumber.NewFileLogger(viper.GetString("log-file"), logLvl, lumber.ROTATE, 5000, 1, 100)
				// lumber.SetLogger(logger)
			}
		},

		// either run mist as a server, or run it as a CLI depending on what flags
		// are provided
		Run: func(ccmd *cobra.Command, args []string) {

			// if --server is passed start the mist server
			if daemon {

				// if --config is passed, attempt to parse a config file
				if config != "" {

					// get the filepath
					abs, err := filepath.Abs(config)
					if err != nil {
						lumber.Error("Error reading filepath: ", err.Error())
					}

					// get the config name
					base := filepath.Base(abs)

					// get the path
					path := filepath.Dir(abs)

					//
					viper.SetConfigName(strings.Split(base, ".")[0])
					viper.AddConfigPath(path)

					// Find and read the config file; Handle errors reading the config file
					if err := viper.ReadInConfig(); err != nil {
						lumber.Fatal("Failed to read config file: ", err.Error())
						os.Exit(1)
					}
				}

				//
				if err := auth.Start(viper.GetString("authenticator")); err != nil {
					lumber.Fatal("Failed to start authenticator: ", err)
					os.Exit(1)
				}

				//
				if err := server.Start(viper.GetStringSlice("listeners"), viper.GetString("token")); err != nil {
					lumber.Fatal("One or more servers failed to start: ", err)
					os.Exit(1)
				}

				//
				// if err := replicator.Start(); err != nil {
				// 	os.Exit(1)
				// }

				// just "hang" out"; this needs to be updated to be/do something real
				done := make(chan bool)
				<-done

				//
				return
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {

	// persistent config flags
	MistCmd.PersistentFlags().String("authenticator", "", "Setting this option enables authentication and uses the authenticator provided to store tokens")
	viper.BindPFlag("authenticator", MistCmd.PersistentFlags().Lookup("authenticator"))

	MistCmd.PersistentFlags().StringSlice("listeners", []string{"tcp://127.0.0.1:1445", "http://127.0.0.1:8080", "ws://127.0.0.1:8888"}, "A comma delimited list of servers to start")
	viper.BindPFlag("listeners", MistCmd.PersistentFlags().Lookup("listeners"))

	MistCmd.PersistentFlags().String("log-type", "stdout", "The type of logging (stdout, file)")
	viper.BindPFlag("log-type", MistCmd.PersistentFlags().Lookup("log-type"))

	MistCmd.PersistentFlags().String("log-file", "/var/log/mist.log", "If log-type=file, the /path/to/logfile; ignored otherwise")
	viper.BindPFlag("log-file", MistCmd.PersistentFlags().Lookup("log-file"))

	MistCmd.PersistentFlags().String("log-level", "INFO", "Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	viper.BindPFlag("log-level", MistCmd.PersistentFlags().Lookup("log-level"))

	MistCmd.PersistentFlags().String("replicator", "", "not yet implemented")
	viper.BindPFlag("replicator", MistCmd.PersistentFlags().Lookup("replicator"))

	MistCmd.PersistentFlags().String("token", "", "Auth token used when connecting to a Mist started with an authenticator")
	viper.BindPFlag("token", MistCmd.PersistentFlags().Lookup("token"))

	MistCmd.PersistentFlags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
	MistCmd.PersistentFlags().StringSliceVar(&tags, "tags", tags, "Tags used when subscribing, publishing, or unsubscribing")

	// local flags;
	MistCmd.Flags().StringVar(&config, "config", config, "/path/to/config.yml")
	viper.BindPFlag("config", MistCmd.Flags().Lookup("config"))

	MistCmd.Flags().BoolVar(&daemon, "server", daemon, "Run mist as a server")
	MistCmd.Flags().BoolVar(&version, "version", version, "Display the current version of this CLI")

	// commands
	MistCmd.AddCommand(pingCmd)
	MistCmd.AddCommand(subscribeCmd)
	MistCmd.AddCommand(publishCmd)
	MistCmd.AddCommand(unsubscribeCmd)
	MistCmd.AddCommand(listCmd)

	// hidden/aliased commands
	MistCmd.AddCommand(messageCmd)
	MistCmd.AddCommand(sendCmd)
}
