package janus

import (
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/configs"
	"github.com/snapp-incubator/ghodrat/internal/server/janus"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	use     = `janus`
	short   = `janus WebRTC stress testing tool`
	long    = `janus  is a CMD tool used to stress test janus WebRTC media servers`
	example = `janus --env dev --call-count 5`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Command() *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: short, Long: long, Example: example, Run: run}

	envFlag := "set config environment, default is dev"
	cmd.Flags().StringP("env", "e", "", envFlag)

	return cmd
}

func run(cmd *cobra.Command, _ []string) {
	env := cmd.Flag("env").Value.String()

	configs := configs.Load(env)

	lg := logger.NewZap(configs.Logger)

	server := janus.Janus{
		Config: configs.Janus,
		Logger: lg,
		Client: &client.Client{
			Config: configs.Client,
			Logger: lg,
		},
	}

	server.TearUp()

	closed := make(chan os.Signal, 1)
	signal.Notify(closed, os.Interrupt)
	<-closed

	server.TearDown()
}
