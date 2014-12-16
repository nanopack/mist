package mist

import (
// "io"
// "net"
)

//
type (

	// Client...
	Client struct {
		Subscriptions map[chan Message]chan bool
	}
)

// New creates a new Mist Client
func (c *Client) New(opts map[string]string) (Client, error) {

	client := Client{}
	client.Subscriptions = make(map[chan Message]chan bool)

	return client, nil
}

// Subscribe
func (c *Client) Subscribe(tags []string) ([]string, error) {

	// go func{

	// }()

	//
	// res, err := http.Get("http://127.0.0.1:1445/mist?subscribe="+strings.Join(tags, ","))
	// if err != nil {
	//   return tags, err
	// }

	// defer res.Body.Close()

	//
	// var b []byte

	//
	// for {
	//   b = make([]byte, bytes.MinRead)

	//   _, err := res.Body.Read(b)
	//   if err != nil {
	//     if err == io.EOF {
	//       break
	//     } else {
	//       return tags, done, err
	//     }
	//   }

	//   b = bytes.Trim(b, "\x00")
	// }

	return tags, nil
}

// Unsubscribe
func (c *Client) Unsubscribe(tags []string) {
	// unsubscribe from a channel
}

// Subscriptions
// func (c *Client) Subscriptions() {
//   // get a list of all subscriptions
// }
