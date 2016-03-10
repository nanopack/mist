package commands

import (
	"github.com/spf13/cobra"
)

var (

	// alias for publish
	messageCmd = &cobra.Command{
		Hidden: true,

		Use:   "message",
		Short: "Publish a message",
		Long:  ``,

		Run: publish,
	}

	//
	publishCmd = &cobra.Command{

		Use:   "publish",
		Short: "Publish a message",
		Long:  ``,

		Run: publish,
	}

	// alias for publish
	sendCmd = &cobra.Command{
		Hidden: true,

		Use:   "send",
		Short: "Publish a message",
		Long:  ``,

		Run: publish,
	}
)

// init
func init() {
}

// publish
func publish(ccmd *cobra.Command, args []string) {
}
