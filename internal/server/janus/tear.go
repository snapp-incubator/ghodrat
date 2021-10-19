package janus

import (
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

func (j *Janus) TearUp() {
	j.initiate()

	go func() { j.readRTCPPackets() }()

	go func() { j.Client.StreamAudioFile(j.iceConnectedCtx, j.audioTrack.WriteSample) }()

	j.Client.CreateAndSetLocalOffer()

	j.call()
}

func (j *Janus) TearDown() {
	j.Client.ClosePeerConnection()

	if err := j.audioWriter.Close(); err != nil {
		j.Logger.Fatal("failed to close audio writer", logger.Error(err))
	}
}
