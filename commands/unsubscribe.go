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
	unsubscribeCmd = &cobra.Command{
		Use:   "unsubscribe",
		Short: "Unsubscribe tags",
		Long:  ``,

		Run: unsubscribe,
	}
)

// unsubscribe
func unsubscribe(ccmd *cobra.Command, args []string) {

	// missing tags
	if tags == nil {
		fmt.Println("Unable to unsubscribe - Missing tags")
		os.Exit(1)
	}

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	//
	client.Unsubscribe(tags)

	//
	msg := <-client.Messages()
	fmt.Println(msg.Data)
}
