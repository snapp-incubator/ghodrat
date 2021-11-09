package janus

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	zap2 "go.uber.org/zap"

	"github.com/snapp-incubator/ghodrat/internal/config"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/vendors/janus/clients"
	janus_server "github.com/snapp-incubator/ghodrat/internal/vendors/janus/server"
	"github.com/spf13/cobra"
)

const (
	use   = `janus`
	short = `janus WebRTC media-server stress testing tool`
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

		af, err := clients.NewAudioFactory(configs.Client)
		if err != nil {
			zap.Panic("failed to create audio factory", zap2.Error(err))
		}

		server := janus_server.Janus{
			Config: configs.Janus,
			Logger: zap,
			Client: &clients.Client{
				Config:       configs.Client,
				Logger:       zap,
				AudioFactory: af,
			},
		}

		go func(server janus_server.Janus) {
			doneChannel := make(chan bool)
			server.TearUp(doneChannel)
			<-doneChannel
			server.TearDown()
			waitGroup.Done()
		}(server)
	}

	waitGroup.Wait()

	lg.Info("all calls has been finished successfully")
}
