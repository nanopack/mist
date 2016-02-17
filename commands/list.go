package commands

import (
	"github.com/spf13/cobra"
)

var (

	//
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all subscriptions",
		Long:  ``,

		Run: list,
	}
)

// init
func init() {
}

// list
func list(ccmd *cobra.Command, args []string) {
}
