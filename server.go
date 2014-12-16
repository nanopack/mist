package mist

import (
	"bufio"
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

//
type (

	//
	Server struct {
		port string
	}
)

// start
func (s *Server) start(p string, m *Mist) error {
	fmt.Printf("Starting Mist server...\n")

	s.port = p

	//
	go func() {

		//
		l, err := net.Listen("tcp", ":"+s.port)
		if err != nil {
			fmt.Printf("Failed to start Mist Server: %v\n", err)
			os.Exit(1)
		}

		//
		defer l.Close()

		fmt.Printf("Mist listening at %v\n", s.port)

		// Listen for an incoming connection.
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}

			// Handle connections in a new goroutine.
			go handleRequest(conn, m)
		}
	}()

	return nil
}

// handleRequest
func handleRequest(conn net.Conn, m *Mist) {

	fmt.Println("HANDLE REQUEST!")

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

				fmt.Printf("SUB: %+v\nMSG:%+v\n", sub.Sub, msg)

				b, err := json.Marshal(msg)
				if err != nil {
					fmt.Printf("Failed to marshal: %v\n", err)
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
				fmt.Println("EOF!")
				conn.Close()
				m.Unsubscribe(sub)
				done <- true
				// close(sub.Sub)
				break
			} else {
				fmt.Println("Error reading:", err.Error())
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
			fmt.Printf("Unknown command: %+v\n", cmd)
		}
	}

	return
}
