package commands

import (
	"fmt"
	"os"

	"github.com/nanopack/mist/clients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var data string //

// init
func init() {
	messageCmd.Flags().StringVar(&data, "data", data, "The string data to message")
	publishCmd.Flags().StringVar(&data, "data", data, "The string data to publish")
	sendCmd.Flags().StringVar(&data, "data", data, "The string data to send")
}

// publish
func publish(ccmd *cobra.Command, args []string) {

	// missing tags
	if tags == nil {
		fmt.Println("Unable to publish - Missing tags")
		os.Exit(1)
	}

	// missing data
	if data == "" {
		fmt.Println("Unable to publish - Missing data")
		os.Exit(1)
	}

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	//
	client.Publish(tags, data)

	//
	msg := <-client.Messages()
	fmt.Println(msg.Data)
}
