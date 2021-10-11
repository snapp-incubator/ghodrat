package cmd

import (
	"log"

	"github.com/snapp-incubator/ghodrat/internal/cmd/janus"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewRootCMD() *cobra.Command {
	root := &cobra.Command{
		Use:   "ghodrat",
		Short: "WebRTC media servers stress testing tool",
		Long:  "Ghodart is a CMD tool used to stress test WebRTC media servers",
		Run: func(cmd *cobra.Command, args []string) {
			logger, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("can't initialize zap logger: %s", err)
			}

			log.Println(cmd.Flags().GetString("address"))

			defer logger.Sync()

			logger.Info("test")
		},
	}

	janus.Register(root)

	return root
}
