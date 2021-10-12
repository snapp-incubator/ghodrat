package webrtc

import (
	"fmt"

	"github.com/pion/webrtc/v3"
)

var STUNservers = []string{
	"stun:stun.l.google.com:19302",
	"stun:stun1.l.google.com:19302",
}

// NewPeerConnectionWithOpusMediaEngine returns webrtc.PeerConnection
func NewPeerConnectionWithOpusMediaEngine() (*webrtc.PeerConnection, error) {
	// Prepare the configuration
	config := webrtc.Configuration{
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
		ICEServers: []webrtc.ICEServer{
			{URLs: STUNservers},
		},
	}

	// A MediaEngine defines the codecs supported by a PeerConnection,
	// and the configuration of those codecs
	mediaEngine := &webrtc.MediaEngine{}
	codec := webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeOpus,
			ClockRate:    48000,
			Channels:     2,
			SDPFmtpLine:  "minptime=10;useinbandfec=1",
			RTCPFeedback: nil,
		},
		PayloadType: 111,
	}

	// Add OPUS codec (audio format)
	if err := mediaEngine.RegisterCodec(codec, webrtc.RTPCodecTypeAudio); err != nil {
		return nil, fmt.Errorf("failed to register opus codec: %w", err)
	}

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	return peerConnection, nil
}
