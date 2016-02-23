package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/core"
)

var (

	//
	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Ping a running mist server",
		Long:  ``,

		Run: ping,
	}
)

// init
func init() {
}

// ping
func ping(ccmd *cobra.Command, args []string) {
}
