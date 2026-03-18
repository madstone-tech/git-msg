package main

import (
	"fmt"
	"os"

	"github.com/madstone0-0/git-msg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
