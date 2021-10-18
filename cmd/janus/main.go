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
	example = `janus --env dev --url ws://127.0.0.1:8188/ --audio-file ./static/audio.ogg --call-count 5`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Command() *cobra.Command {
	cmd := &cobra.Command{Use: use, Short: short, Long: long, Example: example, Run: run}

	cmd.Flags().String("env", "", "set config environment, default is dev")
	cmd.Flags().String("url", "ws://127.0.0.1:8188/", "Janus media server websocket url")
	cmd.Flags().String("audio-file", "./static/audio.ogg", "audio file used to stream to Janus")
	cmd.Flags().Uint("call-count", 1, "number of concurrent calls")

	return cmd
}

func run(cmd *cobra.Command, _ []string) {
	env := cmd.Flag("env").Value.String()

	configs := configs.LoadJanus(env)

	lg := logger.NewZap(configs.Logger)

	url, err := cmd.Flags().GetString("url")
	if err != nil {
		lg.Fatal("failed to get url flag", logger.Error(err))
	}
	lg.Info("using Janus websocket", logger.String("url", url))

	audioFile, err := cmd.Flags().GetString("audio-file")
	if err != nil {
		lg.Fatal("failed to get audio-file flag", logger.Error(err))
	}
	lg.Info("checking if audio file exists", logger.String("audio_file", audioFile))

	_, err = os.Stat(audioFile)
	if os.IsNotExist(err) {
		lg.Fatal("audio file does not exist", logger.String("audio_file", audioFile))
	}

	client := &client.Client{
		Config: nil,
		Logger: nil,
	}

	server := janus.Janus{
		Config: nil,
		Logger: nil,
		Client: client,
	}

	server.Initiate()

	go func() { server.ReadRTCPPackets() }()

	go func() { client.StreamAudioFile() }()

	client.CreateAndSetLocalOffer()

	server.Call()

	closed := make(chan os.Signal, 1)
	signal.Notify(closed, os.Interrupt)
	<-closed

	server.Close()
}
