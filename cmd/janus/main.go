package janus

import (
	"os"

	"github.com/snapp-incubator/ghodrat/internal/cmd"
	"github.com/spf13/cobra"
)

const (
	use   = "janus"
	short = "janus WebRTC stress testing tool"
	long  = `janus  is a CMD tool used to stress test janus WebRTC media servers`
)

func Command() *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: short, Long: long, Run: run}

	envFlag := "set config environment, default is dev"
	cmd.Flags().StringP("env", "e", "", envFlag)

	return cmd
}

func run(_ *cobra.Command, _ []string) {
	c := cmd.NewRootCMD()

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
