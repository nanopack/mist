package clients

import "fmt"

var (
	ErrNotSupported = fmt.Errorf("Unable to perform action: command not supported\n")
)
