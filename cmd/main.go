package main

import (
	"os"

	"github.com/billglover/character-stats/cmd/cli"
)

func main() {
	err := cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
