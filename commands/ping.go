package commands

import (
	"fmt"
	"os"

	"github.com/nanopack/mist/clients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// ping
func ping(ccmd *cobra.Command, args []string) {

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	//
	client.Ping()

	//
	msg := <-client.Messages()
	fmt.Println(msg.Data)
}
