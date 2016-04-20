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
	subscribeCmd = &cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe tags",
		Long:  ``,

		Run: subscribe,
	}
)

// subscribe
func subscribe(ccmd *cobra.Command, args []string) {

	// missing tags
	if tags == nil {
		fmt.Println("Unable to subscribe - Missing tags")
		os.Exit(1)
	}

	client, err := clients.New(host)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	//
	if err := client.Subscribe(tags); err != nil {
		fmt.Printf("Unable to subscribe - %v\n", err.Error())
	}

	// listen for messages on tags
	fmt.Printf("Listening on tags '%v'\n", tags)
	for msg := range client.Messages() {

		// skip handler messages
		if msg.Data != "success" {
			if viper.GetString("log-level") == "DEBUG" {
				fmt.Printf("Message: %#v\n", msg)
			} else {
				fmt.Printf(msg.Data)
			}
		}
	}
}
