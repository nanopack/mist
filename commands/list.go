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

	//
	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	//
	client.List()

	//
	msg := <-client.Messages()
	fmt.Println(msg.Data)

}
