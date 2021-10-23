package cmd

import (
	"log"
	"os"

	"github.com/snapp-incubator/ghodrat/internal/cmd/janus"
	"github.com/spf13/cobra"
)

const (
	errExecuteCMD = "failed to execute root command"

	short = "WebRTC stress testing tool"
	long  = `ghodrat is a CMD tool used to stress test janus WebRTC media servers`
)

// ExitFailure status code.
const ExitFailure = 1

func Execute() {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{Short: short, Long: long}

	cmd.AddCommand(janus.Command())

	if err := cmd.Execute(); err != nil {
		log.Println(errExecuteCMD, err)
		os.Exit(ExitFailure)
	}
}
