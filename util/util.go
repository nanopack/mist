//
package util

import (
  "bufio"
  "net"
	"strings"
)

type (

  //
	reader struct {
		reader *bufio.Reader

		Err    error
		Input    []string
	}
)

//
func NewReader(conn net.Conn) (r *reader) {
	r = &reader{
		reader: bufio.NewReader(conn),
	}

	return
}

//
func (r *reader) Next() bool {
	line, err := r.reader.ReadString('\n')

	if err != nil {
		r.Err = err
		return false
	}

  line = strings.TrimSuffix(line, "\n")

  //
	r.Input = strings.SplitN(line, " ", 3)

  //
	return true
}
