package commands

import (
	"github.com/spf13/cobra"
)

var (

	//
	unsubscribeCmd = &cobra.Command{
		Use:   "unsubscribe",
		Short: "Unsubscribe tags",
		Long:  ``,

		Run: unsubscribe,
	}
)

// init
func init() {
}

// unsubscribe
func unsubscribe(ccmd *cobra.Command, args []string) {
}
