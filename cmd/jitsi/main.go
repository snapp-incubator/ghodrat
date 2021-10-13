package janus

import (
	"github.com/spf13/cobra"
)

const (
	use     = `jitsi`
	short   = `jitsi WebRTC stress testing tool`
	long    = `janus  is a CMD tool used to stress test jitsi WebRTC media servers`
	example = `jitsi --env dev `
)

func Command() *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: short, Long: long, Example: example, Run: run}

	// flags
	cmd.Flags().String("env", "", "set config environment, default is dev")

	return cmd
}

func run(cmd *cobra.Command, _ []string) {
	panic("NOT IMPLEMENTED YET")
}
