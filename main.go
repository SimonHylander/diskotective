package main

import (
	"github.com/simonhylander/diskotective/cmd"
	"github.com/simonhylander/diskotective/cmd/seed"
	"os"
	//"github.com/simonhylander/diskotective/cmd"
)

func main() {
	args := os.Args

	if len(args) == 2 && args[1] == "seed" {
		seed.Execute()
		return
	}

	cmd.Execute()
}
