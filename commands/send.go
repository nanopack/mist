package commands

import (
	"github.com/spf13/cobra"
)

var (

	// alias for send
	messageCmd = &cobra.Command{
		Hidden: true,

		Use:   "message",
		Short: "Sends a message",
		Long:  ``,

		Run: send,
	}

	//
	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Sends a message",
		Long:  ``,

		Run: send,
	}
)

// init
func init() {
}

// send
func send(ccmd *cobra.Command, args []string) {
}
