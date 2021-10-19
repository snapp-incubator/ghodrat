package janus

import (
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

func (j *Janus) TearUp(doneChannel chan bool) {
	j.Logger.Info("Inititate janus")
	j.initiate()

	j.Logger.Info("read RTCP Packets")
	go func() { j.readRTCPPackets() }()

	j.Logger.Info("Stream Audio File")
	go func() { j.Client.StreamAudioFile(j.iceConnectedCtx, j.audioTrack.WriteSample, doneChannel) }()

	j.Logger.Info("Create And Set Local Offer")
	j.Client.CreateAndSetLocalOffer()

	j.Logger.Info("start call")
	j.call()
}

func (j *Janus) TearDown() {
	j.Client.ClosePeerConnection()

	if err := j.audioWriter.Close(); err != nil {
		j.Logger.Fatal("failed to close audio writer", logger.Error(err))
	}
}
