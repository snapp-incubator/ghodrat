package main

import (
	"github.com/moeen/ghodrat/internal/cmd"
	"os"
)

func main() {
	c := cmd.NewRootCMD()

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
