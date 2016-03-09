package commands

import (
	"github.com/spf13/cobra"
)

var (

	//
	subscribeCmd = &cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe tags",
		Long:  ``,

		Run: subscribe,
	}
)

// init
func init() {
}

// subscribe
func subscribe(ccmd *cobra.Command, args []string) {
}
