package main

import (
	"fmt"
	"os"

	"github.com/micnncim/bitbrew/cli"
)

func main() {
	if err := cli.New().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
