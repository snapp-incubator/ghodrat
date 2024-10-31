package janus

import (
	"context"

	"go.uber.org/zap"
)

func (j *Janus) StartCall(doneChannel chan bool) {
	iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())

	j.Client.CreatePeerConnection(iceConnectedCtxCancel)

	j.initiate()

	go j.handle()

	j.Client.ReadTrack(doneChannel, iceConnectedCtx)

	j.Client.CreateAndSetOffer()

	j.Logger.Info("start call")

	if err := j.call(); err != nil {
		j.Logger.Fatal("failed to start a call", zap.Error(err))
	}
}

func (j *Janus) HangUp() {
	j.Client.ClosePeerConnection()
}
