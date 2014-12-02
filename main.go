package mist

import (
  "encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

  "github.com/nanobox-core/utils"
)

const (
	DefaultHost = "localhost"
	DefaultPort = "1445"
	DefaultAddr = DefaultHost + ":" + DefaultPort
	Version     = "0.0.1"
)

//
type (

	//
	Mist struct {
		sync.Mutex //

		Addr          string
		Debugging     bool
		host          string
		port          string
		Subscriptions map[string][]chan string // Subscriptions represent...
	}

  //
  message struct {
    Keys []string
    Data string
  }
)

// New creates a new Mist, setting options, and starting the Mist server
func New(opts map[string]string) (Mist, error) {
	fmt.Println("Initializing 'Mist'...")

	mist := Mist{}

	mist.Subscriptions = make(map[string][]chan string)

	mist.host = utils.SetOption(opts["mist_host"], DefaultHost)
	mist.port = utils.SetOption(opts["mist_port"], DefaultPort)
	mist.Addr = utils.SetOption(opts["mist_addr"], DefaultAddr)

	mist.Debugging = true

	fmt.Printf("Starting Mist server...\n")
	fmt.Printf("Mist listening at %v\n", mist.Addr)

	// start the server
	mist.start()

	return mist, nil
}

// start creates a basic http server to listen for incoming mist subscriptions
func (m *Mist) start() {

	//
	http.HandleFunc("/mist", m.handler)

	//
	go func() {
		if err := http.ListenAndServe(m.Addr, nil); err != nil {
			panic(err)
		}
	}()
}

// handler responds to client requests to subscribe to mist updates. It expects
// subcriptions in the form of '?subscribe=a,b,c'. It then iterates through each
// subscirption, creating a subscription in Mist, and setting up a select loop
// to handle all incoming publish requests.
func (m *Mist) handler(rw http.ResponseWriter, req *http.Request) {

	// pull out query string params and look for a 'subscribe' param
	params := req.URL.Query()
	tags := params["subscribe"][0]

	//
	if m.Debugging {
		fmt.Printf("Mist request with params: '%+v'\n", params)
	}

	// create a wait group to keep the connection open until all messages have been
	// published
	var wg sync.WaitGroup

	// if no tags are passed, send a message indicating that the server is up, but
	// a subscription is needed
	if len(tags) <= 0 {
		rw.Write([]byte("Mist is up... subscribe to receive updates (?subscribe=a,b,c)"))
		rw.(http.Flusher).Flush()
		return
	}

	// iterate over each tag, and create a subscription with Mist
	for _, t := range strings.Split(tags, ",") {

		// increment wait group for each subscription created
		wg.Add(1)

		// create a subscription for each tag
		sub := m.Subscribe([]string{t})

		// create our 'publish handler'
		go func() {
			for {
				select {

				// if a message is publishd on the channel write it to the http writer
				// and flush to the client
				case msg := <-sub:

					//
					if m.Debugging {
						fmt.Printf("Message received: %+v\n", msg)
					}

					if msg == "\n\r" {
						m.Unsubscribe(sub)
            wg.Done()
						return
					}

					// write message and flush
					go func() {
						rw.Write([]byte(msg))
						rw.(http.Flusher).Flush()
					}()
        }
			}
		}()
	}

	// hold the client open until all channels have finished receiving updates
	wg.Wait()
	return

}

// Publish takes a list of tags and iterates through Mist's list of subscriptions,
// looking for matching subscriptions to publish messages too. It ensures that the
// list of recipients is a unique set, so as not to publish the same message more
// than once over a channel
func (m *Mist) Publish(tags []string, data string) {

	// a complete list of recipients (may contain duplicate channels from multiple
	// subscriptions)
	found := make(map[chan string]int)

	// a *unique* list of recipients that will receive broadcasts
	var recipients []chan string

	// iterate through each provided tag looking for subscriptions to publish to
	for _, t := range tags {

		// keep track of how many times a subscription is requested
		used := 0

		// iterate through any matching subscriptions (type []chan string) and add
		// all of that subscriptions channels to the list of recipients
		if sub, ok := m.Subscriptions[t]; ok {
			for _, ch := range sub {

				// ensure that we keep the list of recipients unique, by checking each
				// match against a temporary map of found channels.
				if _, ok := found[ch]; !ok {
					used++

					// update our list of found channels, with a value of how many times
					// that channel has been subscribed to
					found[ch] = used

					// add the channel to our unique list of channels
					recipients = append(recipients, ch)
				}
			}
		}
	}

	// format the data and send it on each unique recipient's channel
  msg := message{Keys: tags, Data: data}

  //
  if m.Debugging {
    fmt.Printf("Publishing: %+v\n", msg)
  }

  b, err := json.Marshal(msg)
  if err != nil {
    panic(err)
  }

  //
	for _, r := range recipients {
		go func() { r <- string(b) }()
	}
}

// Subscribe takes a slice of strings, iterates through each one, and creates a
// new subscription (type []chan string), if it doesn't already exist. It then
// adds a channel under that subscription which will be used to communicate when
// publishing
func (m *Mist) Subscribe(tags []string) chan string {

	//
	if m.Debugging {
		fmt.Printf("Subscribing to: %+v\n", tags)
	}

	// make the channel that will be used when publishing subscriptions. This
	// channel is a 'one-to-many' relationship, in that, it will be used to when
	// publishing messages to one, or any, of the provided subscriptions.
	ch := make(chan string)

	// iterate over each subscription, adding it to our list of subscriptions (if
	// not already found), and then adding the channel into the subscription's list
	// of subscribers.
	for _, t := range tags {

		// if we don't find a subscription, make one (type []chan string), and add
		// it to our list of subscriptions
		if _, ok := m.Subscriptions[t]; !ok {

			//
			if m.Debugging {
				fmt.Printf("Creating subscription: %+v\n", t)
			}

			// new subscription
			var sub []chan string

			// add subscription to list of subscriptions
			m.Lock()
			m.Subscriptions[t] = sub
			m.Unlock()
		}

		// add the channel to each subscription...
		m.Subscriptions[t] = append(m.Subscriptions[t], ch)
	}

	//
	if m.Debugging {
		fmt.Printf("Current subscriptions: %+v\n", m.Subscriptions)
	}

	// and return the channel
	return ch
}

// Unsubscribe iterates through each of Mist's subscriptions looking for subscriptions
// that contain a match for the channel provided. The channel is removed from that
// subscriptions list, and closed. If a subscription is found empty, it is removed
func (m *Mist) Unsubscribe(ch chan string) bool {

	//
	if m.Debugging {
		fmt.Printf("Unsubscribing: '%+v'\n", ch)
	}

	// iterate over Mist's subscriptions looking for subscriptions that have the
	// channel to unsubscribe
	for k, v := range m.Subscriptions {

		// hold all the remaining channels that != the channel to unsubscribe
		var remaining []chan string

		// iterate over each channel in the subscription looking for the channel to
		// unsubscribe
		for _, c := range v {

			// if the channel found isn't the one to unsubscribe, add it to the list of
			// remaining channels
			if c != ch {
				remaining = append(remaining, c)

				// if the channel is found close it
			} else {
				close(ch)

				//
				if m.Debugging {
					fmt.Printf("'%+v' closed\n", ch)
				}
			}

			// set Mist's subscriptions equal to the remaining subscriptions
			m.Subscriptions[k] = remaining

			// if a subscription is empty, remove it
			if len(m.Subscriptions[k]) <= 0 {
				m.Lock()
				delete(m.Subscriptions, k)
				m.Unlock()
			}
		}
	}

	return true
}
