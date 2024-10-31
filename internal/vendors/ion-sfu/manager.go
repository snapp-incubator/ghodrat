package ion_sfu

import "context"

func (ionSfu *IonSfu) StartCall(doneChannel chan bool) {
	iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())

	ionSfu.generateSID()
	ionSfu.dial()

	ionSfu.Client.CreatePeerConnection(iceConnectedCtxCancel)

	go ionSfu.readMessage()

	ionSfu.Client.ReadTrack(doneChannel, iceConnectedCtx)

	ionSfu.Client.CreateAndSetOffer()
	ionSfu.Client.OnIceCandidate(ionSfu.onIceCandidate)
	ionSfu.offer()
}

func (ionSfu *IonSfu) HangUp() {
	ionSfu.Client.ClosePeerConnection()
}
