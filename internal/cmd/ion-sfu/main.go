package janus

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

const (
	use   = `ion-sfu`
	short = `pion ion-sfu WebRTC media-server stress testing tool`
)

func Command() *cobra.Command {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{Use: use, Short: short, PreRun: preRun, Run: run}

	return cmd
}

func preRun(cmd *cobra.Command, _ []string) {
	rand.Seed(time.Now().UnixNano())
}

func run(cmd *cobra.Command, _ []string) {
	panic("NOT IMPLEMENTED YET")
}
