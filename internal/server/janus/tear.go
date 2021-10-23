package janus

import "go.uber.org/zap"

func (j *Janus) TearUp(doneChannel chan bool) {
	j.Logger.Info("inititate janus")
	j.initiate()

	j.Logger.Info("read RTCP packets")

	go func() { j.readRTCPPackets() }()

	j.Logger.Info("stream Audio File")

	go func() { j.Client.StreamAudioFile(j.iceConnectedCtx, j.audioTrack.WriteSample, doneChannel) }()

	j.Logger.Info("create and set local offer")
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
