package commands

import (
	"fmt"

	"github.com/nanopack/mist/clients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	//
	listCmd = &cobra.Command{
		Use:           "list",
		Short:         "List all subscriptions",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: list,
	}
)

// init
func init() {
	listCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
}

// list
func list(ccmd *cobra.Command, args []string) error {

	// create new mist client
	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf("Failed to connect to '%v' - %v\n", host, err)
		return err
	}

	err = client.List()
	if err != nil {
		fmt.Printf("Failed to list - %v\n", err)
		return err
	}

	msg := <-client.Messages()
	fmt.Println(msg.Data)

	return nil
}
