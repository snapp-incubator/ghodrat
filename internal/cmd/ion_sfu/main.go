package ion_sfu

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/config"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	ion_sfu "github.com/snapp-incubator/ghodrat/internal/vendors/ion-sfu"
	"github.com/spf13/cobra"
)

const (
	use   = `ion-sfu`
	short = `ion-sfu WebRTC media-server stress testing tool`
)

func Command() *cobra.Command {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{Use: use, Short: short, Run: run, PreRun: preRun}

	return cmd
}

func preRun(cmd *cobra.Command, _ []string) {
	rand.Seed(time.Now().UnixNano())
}

func run(_ *cobra.Command, _ []string) {
	configs := config.New()

	lg := logger.NewZap(configs.Logger)

	var waitGroup sync.WaitGroup
	waitGroup.Add(configs.CallCount)

	for index := 0; index < configs.CallCount; index++ {
		zap := lg.Named(fmt.Sprintf("goroutine: %d", index+1))

		server := ion_sfu.Ion_sfu{
			Config: configs.Ion_sfu,
			Logger: zap,
			Client: &client.Client{
				Config: configs.Client,
				Logger: zap,
			},
		}

		go func(server ion_sfu.Ion_sfu) {
			doneChannel := make(chan bool)
			server.StartCall(doneChannel)
			<-doneChannel
			server.HangUp()
			waitGroup.Done()
		}(server)
	}

	waitGroup.Wait()
	lg.Info("all calls has been finished successfully")
}
