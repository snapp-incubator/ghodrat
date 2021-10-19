package janus

import (
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

func (j *Janus) TearUp(index int, doneChannel chan bool) {
	j.Logger.Info("Inititate janus", logger.Int("index", index))
	j.initiate()

	j.Logger.Info("read RTCP Packets", logger.Int("index", index))
	go func() { j.readRTCPPackets() }()

	j.Logger.Info("Stream Audio File", logger.Int("index", index))
	go func() { j.Client.StreamAudioFile(j.iceConnectedCtx, j.audioTrack.WriteSample, doneChannel) }()

	j.Logger.Info("Create And Set Local Offer", logger.Int("index", index))
	j.Client.CreateAndSetLocalOffer()

	j.Logger.Info("start call", logger.Int("index", index))
	j.call()
}

func (j *Janus) TearDown() {
	j.Client.ClosePeerConnection()

	if err := j.audioWriter.Close(); err != nil {
		j.Logger.Fatal("failed to close audio writer", logger.Error(err))
	}
}
