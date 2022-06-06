package main

import (
	"fmt"
	"os"

	"github.com/edsonmichaque/buzi/internal/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		if newErr, ok := err.(cmd.CmdError); ok {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(newErr.ExitCode)
		}
	}
}
