package mist

import (
	"fmt"
	"net"
	"strings"

	"github.com/nanopack/mist/util"
)

var (

	// a map of the available commands that the server will respond to; format is
	// "name":{#args, handler}
	serverCommands = map[string]Handler{
		"ping":               {0, handlePing},
		"list":               {0, handleList},
		"subscribe":          {1, handleSubscribe},
		"unsubscribe":        {1, handleUnubscribe},
		"publish":            {2, handlePublish},
		"enable-replication": {0, handleEnableReplication},
	}
)

type (
	//
	Handler struct {
		ArgCount int
		Handle   func(Client, []string) string
	}
)

// newConnection takes an incoming connection from a mist client (or other client)
// and sets up a new subscription for that connection, and a 'publish Handler'
// that is used to publish messages to the data channel of the subscription
func newConnection(client Client, conn net.Conn, commands map[string]Handler) {

	// make a done channel
	done := make(chan bool)

	// clean up everything that we have setup
	defer func() {
		conn.Close()
		client.Close()
		close(done)
	}()

	// create a 'publish handler' for this connection
	go func() {
		for {

			// when a message is recieved on the subscriptions channel write the message
			// to the connection
			select {
			case msg := <-client.Messages():

				if _, err := conn.Write([]byte(fmt.Sprintf("publish %v %v\n", strings.Join(msg.Tags, ","), msg.Data))); err != nil {
					break
				}

			// return if we are done
			case <-done:
				return
			}
		}
	}()

	//
	r := util.NewReader(conn)
	for r.Next() {

		// what should we do with errors?
		if r.Err != nil {
			// r.Err
		}

		cmd := r.Input[0]
		args := r.Input[1:]

		//
		handler, found := commands[cmd]

		//
		var response string
		switch {

		// no command found
		case !found:
			response = fmt.Sprintf("Error: Unknown Command '%s'", cmd)

		// incorrect number of arguments for command
		case handler.ArgCount != len(args):
			response = fmt.Sprintf("Error: Wrong number of arguments for '%v'. Expected %v got %v.", cmd, handler.ArgCount, len(args))

		// execute command
		default:
			response = handler.Handle(client, args)
		}

		// only send if a response is given
		if response != "" {
			if _, err := conn.Write([]byte(response + "\n")); err != nil {
				break
			}
		}
	}
}

//
func handlePing(client Client, args []string) string {
	return "pong"
}

//
func handleSubscribe(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Subscribe(tags)
	return fmt.Sprintf("subscribed '%v'", tags)
}

//
func handleUnubscribe(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	client.Unsubscribe(tags)
	return fmt.Sprintf("unsubscribed '%v'", tags)
}

//
func handlePublish(client Client, args []string) string {
	tags := strings.Split(args[0], ",")
	sub := args[1]
	client.Publish(tags, sub)
	return fmt.Sprintf("published '%v' to '%v'", tags, sub)
}

//
func handleList(client Client, args []string) string {
	list, err := client.List()
	if err != nil {
		return err.Error()
	}
	tmp := make([]string, len(list))

	for idx, subscription := range list {
		tmp[idx] = strings.Join(subscription, ",")
	}

	response := strings.Join(tmp, " ")
	return fmt.Sprintf("list %v", response)
}

//
func handleEnableReplication(client Client, args []string) string {
	client.(Replicatable).EnableReplication()
	return "replication enabled"
}
