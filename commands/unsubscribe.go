package commands

import (
	"fmt"

	"github.com/nanopack/mist/clients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	//
	unsubscribeCmd = &cobra.Command{
		Use:           "unsubscribe",
		Short:         "Unsubscribe tags",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: unsubscribe,
	}
)

func init() {
	unsubscribeCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
	unsubscribeCmd.Flags().StringSliceVar(&tags, "tags", tags, "Tags to unsubscribe from")
}

// unsubscribe. seems like a useless command, since it creates a new client, then unsubscribes itself
func unsubscribe(ccmd *cobra.Command, args []string) error {

	// missing tags
	if tags == nil {
		fmt.Println("Unable to unsubscribe - Missing tags")
		return fmt.Errorf("")
	}

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf("Failed to connect to '%v' - %v\n", host, err)
		return err
	}

	err = client.Unsubscribe(tags)
	if err != nil {
		fmt.Printf("Failed to unsubscribe - %v\n", err)
		return err
	}

	fmt.Println("success")

	return nil
}
