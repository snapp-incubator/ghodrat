package janus

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/configs"
	"github.com/snapp-incubator/ghodrat/internal/server/janus"
	"github.com/snapp-incubator/ghodrat/pkg/zap"
	"github.com/spf13/cobra"
)

const (
	use     = `janus`
	short   = `janus WebRTC stress testing tool`
	long    = `janus  is a CMD tool used to stress test janus WebRTC media servers`
	example = `janus --env dev --call-count 5`
)

func Command() *cobra.Command {
	// nolint: exhaustivestruct
	cmd := &cobra.Command{Use: use, Short: short, Long: long, Example: example, Run: run, PreRun: preRun}

	envFlag := "set config environment, default is dev"
	cmd.Flags().StringP("env", "e", "", envFlag)

	return cmd
}

func preRun(cmd *cobra.Command, _ []string) {
	rand.Seed(time.Now().UnixNano())
}

func run(cmd *cobra.Command, _ []string) {
	env := cmd.Flag("env").Value.String()

	configs := configs.Load(env)

	lg := .NewZap(configs.Logger)

	var waitGroup sync.WaitGroup

	waitGroup.Add(configs.CallCount)

	for index := 0; index < configs.CallCount; index++ {
		zap := lg.Named(fmt.Sprintf("goroutine: %d", index+1))

		server := janus.Janus{
			Config: configs.Janus,
			Logger: zap,
			Client: &client.Client{
				Config: configs.Client,
				Logger: zap,
			},
		}

		go func(server janus.Janus) {
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
