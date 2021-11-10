package ion_sfu

import "context"

func (ion_sfu *Ion_sfu) StartCall(doneChannel chan bool) {
	iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())

	ion_sfu.generateSID()
	ion_sfu.dial()

	ion_sfu.Client.CreatePeerConnection(iceConnectedCtxCancel)

	go ion_sfu.readMessage()

	ion_sfu.Client.ReadTrack(doneChannel, iceConnectedCtx)

	ion_sfu.Client.CreateAndSetOffer()
	ion_sfu.Client.OnIceCandidate(ion_sfu.onIceCandidate)
	ion_sfu.offer()
}

func (j *Ion_sfu) HangUp() {
	j.Client.ClosePeerConnection()
}
