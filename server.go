package mist

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"strings"
)

// start
func (m *Mist) start() {
	m.log.Info("Starting Mist server...\n")

	//
	go func() {

		//
		l, err := net.Listen("tcp", ":"+m.port)
		if err != nil {
			m.log.Error("%+v\n", err)
		}

		defer l.Close()

		m.log.Info("Mist listening at %+v\n", m.port)

		// Listen for an incoming connection.
		for {
			conn, err := l.Accept()
			if err != nil {
				m.log.Error("%+v\n", err)
			}

			// Handle connections in a new goroutine.
			go m.handleRequest(conn)
		}
	}()
}

// handleRequest
func (m *Mist) handleRequest(conn net.Conn) {

	m.log.Info("Handle request\n")

	var cmd string
	var tags string

	//
	r := bufio.NewReader(conn)

	//
	sub := Subscription{
		Sub: make(chan Message),
	}

	//
	done := make(chan bool)

	// create our 'publish handler'
	go func() {
		for {
			select {

			//
			case msg := <-sub.Sub:

				m.log.Info("SUB: %+v\nMSG:%+v\n", sub.Sub, msg)

				b, err := json.Marshal(msg)
				if err != nil {
					m.log.Error("Failed to marshal: %v\n", err)
				}

				if _, err := conn.Write(b); err != nil {
					break
				}

			//
			case <-done:
				break

			}
		}
	}()

	//
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				m.log.Info("EOF!\n")
				conn.Close()
				m.Unsubscribe(sub)
				done <- true
				// close(sub.Sub)
				break
			} else {
				m.log.Error("Error reading: %+v\n", err.Error())
			}
		}

		split := strings.Split(strings.TrimSpace(l), " ")
		cmd = split[0]

		if len(split) > 1 {
			tags = split[1]
		}

		// if no tags are passed, send a message indicating that the server is up, but
		// a subscription is needed
		if len(tags) <= 0 {
			conn.Write([]byte("Mist Server is running, subscribe to receive updates..."))
		}

		// create a subscription for each tag
		sub.Tags = strings.Split(tags, ",")

		//
		switch cmd {
		case "subscribe":
			m.Subscribe(sub)
		case "unsubscribe":
			m.Unsubscribe(sub)
		case "subscriptions":
			m.List()
		default:
			m.log.Error("Unknown command: %+v\n", cmd)
		}
	}

	return
}
