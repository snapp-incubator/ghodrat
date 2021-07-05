package janus

import (
	z "github.com/moeen/ghodrat/internal/logger/zap"
	"github.com/moeen/ghodrat/internal/webrtc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func run(cmd *cobra.Command, args []string) {
	logger, err := z.InitZapLogger()
	if err != nil {
		log.Fatal(err)
	}

	address, err := cmd.Flags().GetString("address")
	if err != nil {
		logger.Fatal("failed to get address flag", zap.Error(err))
	}

	logger.Info("using Janus websocket", zap.String("address", address))

	audioFile, err := cmd.Flags().GetString("audio-file")
	if err != nil {
		logger.Fatal("failed to get audio-file flag", zap.Error(err))
	}

	logger.Info("checking if audio file exists", zap.String("audio_file", audioFile))

	_, err = os.Stat(audioFile)
	if os.IsNotExist(err) {
		logger.Fatal("audio file does not exist", zap.String("audio_file", audioFile))
	}

	call, err := webrtc.NewCall(logger.Named("caller"))
	if err != nil {
		logger.Fatal("failed to create the call", zap.Error(err))
	}

	go func() {
		if err := call.ReadRTCPPackets(); err != nil {
			logger.Error("failed to read RTCP packets", zap.Error(err))
		}
	}()

	go func() {
		if err := call.StreamAudioFile(audioFile); err != nil {
			logger.Fatal("failed to stream audio", zap.Error(err))
		}
	}()

	if err := call.CreateAndSetLocalOffer(); err != nil {
		logger.Fatal("failed to create and set local SDP offer", zap.Error(err))
	}

	if err := call.Call(address); err != nil {
		logger.Fatal("failed to call via janus", zap.Error(err))
	}

	closed := make(chan os.Signal, 1)
	signal.Notify(closed, os.Interrupt)
	<-closed

	if err := call.Close(); err != nil {
		logger.Error("failed to close the call", zap.Error(err))
	}

	return
}
