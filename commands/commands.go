//
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
	config  string //
	daemon  bool   //
	version bool   //

	//
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

			// if --config is passed, attempt to parse the config file
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
		},

		// either run mist as a server, or run it as a CLI depending on what flags
		// are provided
		Run: func(ccmd *cobra.Command, args []string) {

			// if --server is passed start the mist server
			if daemon {

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

	// set config defaults; these are overriden if a --config file is provided
	// (see above)
	viper.SetDefault("log-type", "stdout")
	viper.SetDefault("log-file", "/var/log/mist.log")
	viper.SetDefault("log-level", "INFO")
	viper.SetDefault("listeners", []string{"tcp://127.0.0.1:1445", "http://127.0.0.1:8080", "ws://127.0.0.1:8888"})
	viper.SetDefault("replicator", "")
	viper.SetDefault("token", "")

	// persistent flags; these are the only 2 options that we want overridable from
	// the CLI, all others need to use a config file
	MistCmd.PersistentFlags().String("authenticator", viper.GetString("authenticator"), "Setting this option enables authentication and uses the authenticator provided to store tokens")
	MistCmd.PersistentFlags().StringSlice("listeners", viper.GetStringSlice("listeners"), "A comma delimited list of servers to start")
	MistCmd.PersistentFlags().String("log-type", viper.GetString("log-type"), "The type of logging (stdout, file)")
	MistCmd.PersistentFlags().String("log-file", viper.GetString("log-file"), "If log-type=file, the /path/to/logfile; ignored otherwise")
	MistCmd.PersistentFlags().String("log-level", viper.GetString("log-level"), "Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	MistCmd.PersistentFlags().String("replicator", viper.GetString("replicator"), "not yet implemented")
	MistCmd.PersistentFlags().String("token", viper.GetString("token"), "Auth token used when connecting to a Mist started with an authenticator")

	viper.BindPFlag("authenticator", MistCmd.PersistentFlags().Lookup("authenticator"))
	viper.BindPFlag("listeners", MistCmd.PersistentFlags().Lookup("listeners"))
	viper.BindPFlag("log-level", MistCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("replicator", MistCmd.PersistentFlags().Lookup("replicator"))
	viper.BindPFlag("token", MistCmd.PersistentFlags().Lookup("token"))

	// local flags;
	MistCmd.Flags().StringVar(&config, "config", "", "/path/to/config.yml")
	MistCmd.Flags().BoolVar(&daemon, "server", false, "Run mist as a server")
	MistCmd.Flags().BoolVarP(&version, "version", "v", false, "Display the current version of this CLI")

	// commands
	MistCmd.AddCommand(pingCmd)
	MistCmd.AddCommand(subscribeCmd)
	MistCmd.AddCommand(unsubscribeCmd)
	MistCmd.AddCommand(publishCmd)
	MistCmd.AddCommand(listCmd)

	// hidden/aliased commands
	MistCmd.AddCommand(messageCmd)
	MistCmd.AddCommand(sendCmd)
}
