package janus

import (
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/snapp-incubator/ghodrat/internal/webrtc"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
	"github.com/spf13/cobra"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func run(cmd *cobra.Command, args []string) {
	lg := logger.NewZap(nil)

	address, err := cmd.Flags().GetString("address")
	if err != nil {
		lg.Fatal("failed to get address flag", logger.Error(err))
	}

	lg.Info("using Janus websocket", logger.String("address", address))

	audioFile, err := cmd.Flags().GetString("audio-file")
	if err != nil {
		lg.Fatal("failed to get audio-file flag", logger.Error(err))
	}

	lg.Info("checking if audio file exists", logger.String("audio_file", audioFile))

	_, err = os.Stat(audioFile)
	if os.IsNotExist(err) {
		lg.Fatal("audio file does not exist", logger.String("audio_file", audioFile))
	}

	call, err := webrtc.NewCall(address, lg.Named("caller"))
	if err != nil {
		lg.Fatal("failed to create the call", logger.Error(err))
	}

	go func() {
		if err := call.ReadRTCPPackets(); err != nil {
			lg.Error("failed to read RTCP packets", logger.Error(err))
		}
	}()

	go func() {
		if err := call.StreamAudioFile(audioFile); err != nil {
			lg.Fatal("failed to stream audio", logger.Error(err))
		}
	}()

	if err := call.CreateAndSetLocalOffer(); err != nil {
		lg.Fatal("failed to create and set local SDP offer", logger.Error(err))
	}

	if err := call.Call(); err != nil {
		lg.Fatal("failed to call via janus", logger.Error(err))
	}

	closed := make(chan os.Signal, 1)
	signal.Notify(closed, os.Interrupt)
	<-closed

	if err := call.Close(); err != nil {
		lg.Error("failed to close the call", logger.Error(err))
	}
}
