//
package main

import (
	"fmt"
	"os"

	"github.com/nanopack/mist/commands"
)

//
func main() {

	//
	if err := commands.MistCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
