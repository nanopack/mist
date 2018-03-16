package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/clients"
)

var (
	pingCmd = &cobra.Command{
		Use:           "ping",
		Short:         "Ping a running mist server",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: ping,
	}
)

func init() {
	pingCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
}

// ping
func ping(ccmd *cobra.Command, args []string) error {

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf("Failed to connect to '%s' - %s\n", host, err.Error())
		return err
	}

	err = client.Ping()
	if err != nil {
		fmt.Printf("Failed to ping - %s\n", err.Error())
		return err
	}

	msg := <-client.Messages()
	fmt.Println(msg.Data)

	return nil
}
