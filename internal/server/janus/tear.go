package janus

import (
	"go.uber.org/zap"
)

func (j *Janus) TearUp(doneChannel chan bool) {
	j.initiate()
	go func() { j.readRTCPPackets() }()

	go func() { j.Client.StreamAudioFile(j.iceConnectedCtx, j.audioTrack.WriteSample, doneChannel) }()

	j.Client.CreateAndSetLocalOffer()

	j.Logger.Info("start call")
	if err := j.call(); err != nil {
		j.Logger.Fatal("failed to start a call", zap.Error(err))
	}
}

func (j *Janus) TearDown() {
	j.Client.ClosePeerConnection()
	j.Client.CloseOpusTrack()
}
