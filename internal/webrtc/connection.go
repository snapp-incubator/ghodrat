package webrtc

import (
	"fmt"

	"github.com/pion/webrtc/v3"
)

// NewPeerConnection returns webrtc-peer-connection with opus media-engine
func (manager *Manager) NewPeerConnection() (*webrtc.PeerConnection, error) {
	// A MediaEngine defines the codecs supported by a PeerConnection
	mediaEngine := &webrtc.MediaEngine{}

	// configuration of OPUS codec
	codec := webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeOpus,
			ClockRate:    manager.Config.Connection.RTPCodec.ClockRate,
			Channels:     manager.Config.Connection.RTPCodec.Channels,
			SDPFmtpLine:  "minptime=10;useinbandfec=1",
			RTCPFeedback: nil,
		},
		PayloadType: webrtc.PayloadType(manager.Config.Connection.RTPCodec.PayloadType),
	}

	// Add OPUS codec (audio format)
	if err := mediaEngine.RegisterCodec(codec, webrtc.RTPCodecTypeAudio); err != nil {
		return nil, fmt.Errorf("failed to register opus codec: %w", err)
	}

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	// Prepare the configuration
	config := webrtc.Configuration{
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
		ICEServers: []webrtc.ICEServer{
			{URLs: manager.Config.Connection.STUNServers},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	return peerConnection, nil
}
