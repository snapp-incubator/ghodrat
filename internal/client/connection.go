package client

import (
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

// InitiatePeerConnection returns webrtc-peer-connection with opus media-engine
func (manager *Client) InitiatePeerConnection() {
	var err error

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
	if err = mediaEngine.RegisterCodec(codec, webrtc.RTPCodecTypeAudio); err != nil {
		manager.Logger.Fatal("failed to register opus codec", logger.Error(err))
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
	manager.connection, err = api.NewPeerConnection(config)
	if err != nil {
		manager.Logger.Fatal("failed to close peer connection", logger.Error(err))
	}
}

func (manager *Client) OnICEConnectionStateChange(callback func(webrtc.ICEConnectionState)) {
	manager.connection.OnICEConnectionStateChange(callback)
}

// NewPeerConnection returns webrtc-peer-connection with opus media-engine
func (manager *Client) ClosePeerConnection() {
	if err := manager.connection.Close(); err != nil {
		manager.Logger.Fatal("failed to close peer connection", logger.Error(err))
	}
}
