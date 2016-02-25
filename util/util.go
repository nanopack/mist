//
package util

import (
	"bufio"
	"io"
	"strings"
)

type (

	//
	reader struct {
		reader *bufio.Reader

		Err   error
		Input struct {
			Cmd string
			Args []string
		}
	}
)

//
func NewReader(r io.Reader) *reader {
	return &reader{
		reader: bufio.NewReader(r),
	}
}

//
func (r *reader) Next() bool {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return false
	}

	// split the line into 3 segments; we want the line broken into three segemnts
	// (the command, tags, and remaining args):
	// "ping" => ["ping"]
	// "subscribe tag,tag" => ["subscribe", "tag,tag"]
	// "publish tag,tag publish message" = ["publish", "tag,tag", "publish message"]
	split := strings.SplitN(strings.TrimSuffix(line, "\n"), " ", 3)

	//
	r.Input.Cmd = split[0]
	r.Input.Args = split[1:]

	//
	return true
}

//
func HaveSameTags(a, b []string) bool {
	for _, vala := range a {
		for _, valb := range b {
			if vala == valb {
				return true
			}
		}
	}
	return false
}
