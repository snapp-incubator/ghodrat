package ion

import (
	"fmt"

	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/config"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/vendors/ion_sfu"
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

	engine := ion_sfu.NewEngine(nil, lg)

	for index := 0; index < configs.CallCount; index++ {
		zap := lg.Named(fmt.Sprintf("goroutine: %d", index+1))

		c := &client.Client{
			Config: configs.Client,
			Logger: zap,
		}

		peer, err := engine.NewClient(c)
		if err != nil {
			zap.Panic("failed to create and start client")
		}

		zap.Info("starting ion call")

		// Will change for PubSub mode
		if err = peer.Call(ion_sfu.SubOnly); err != nil {
			zap.Panic("failed to initiate call", zap2.Error(err))
		}
		// it's probably not a good idea.
		defer peer.HangUp()
	}

	select {}
}
