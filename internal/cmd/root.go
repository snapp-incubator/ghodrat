package cmd

import (
	"os"

	"github.com/snapp-incubator/ghodrat/internal/cmd/janus"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
		logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
