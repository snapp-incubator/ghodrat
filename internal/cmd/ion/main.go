package ion

import (
	"fmt"

	"github.com/snapp-incubator/ghodrat/internal/config"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	ion_server "github.com/snapp-incubator/ghodrat/internal/vendors/ion-sfu/server"
	"github.com/snapp-incubator/ghodrat/internal/vendors/janus/clients"
	"github.com/spf13/cobra"
	zap2 "go.uber.org/zap"
)

const (
	use     = `ion`
	short   = `ion-sfu WebRTC stress testing tool`
	long    = `ion is a CMD tool used to stress test ion-sfu WebRTC media server`
	example = `ion`
)

func Command() *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: short, Long: long, Example: example, Run: run}

	return cmd
}

func run(cmd *cobra.Command, _ []string) {
	configs := config.New()

	lg := logger.NewZap(configs.Logger)

	//var waitGroup sync.WaitGroup

	//waitGroup.Add(configs.CallCount)

	engine := ion_server.NewEngine(configs.Ion, lg)

	for index := 0; index < configs.CallCount; index++ {
		zap := lg.Named(fmt.Sprintf("goroutine: %d", index+1))

		af, err := clients.NewAudioFactory(configs.Client)
		if err != nil {
			zap.Panic("failed to create audio factory", zap2.Error(err))
		}

		c := &clients.Client{
			Config:       configs.Client,
			Logger:       zap,
			AudioFactory: af,
		}

		peer, err := engine.NewClient(c)
		if err != nil {
			zap.Panic("failed to create and start client")
		}

		zap.Info("starting ion call")

		// Will change for PubSub mode
		if err = peer.Call(ion_server.SubOnly); err != nil {
			zap.Panic("failed to initiate call", zap2.Error(err))
		}
		// it's probably not a good idea.
		defer peer.HangUp()
	}

	select {}
}
