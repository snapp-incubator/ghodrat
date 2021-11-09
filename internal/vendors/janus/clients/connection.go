package clients

import (
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

// InitiatePeerConnection returns webrtc-peer-connection with opus media-engine.
func (client *Client) CreatePeerConnection() {
	var err error

	// A MediaEngine defines the codecs supported by a PeerConnection
	mediaEngine := &webrtc.MediaEngine{}

	// configuration of OPUS codec
	codec := webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeOpus,
			ClockRate:    client.Config.Connection.RTPCodec.ClockRate,
			Channels:     client.Config.Connection.RTPCodec.Channels,
			SDPFmtpLine:  "minptime=10;useinbandfec=1",
			RTCPFeedback: nil,
		},
		PayloadType: webrtc.PayloadType(client.Config.Connection.RTPCodec.PayloadType),
	}

	// Add OPUS codec (audio format)
	if err = mediaEngine.RegisterCodec(codec, webrtc.RTPCodecTypeAudio); err != nil {
		client.Logger.Fatal("failed to register opus codec", zap.Error(err))
	}

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	// Prepare the configuration
	// nolint: exhaustivestruct
	config := webrtc.Configuration{
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{client.Config.Connection.STUNServer}},
		},
	}

	// Create a new RTCPeerConnection
	client.connection, err = api.NewPeerConnection(config)
	if err != nil {
		client.Logger.Fatal("failed to create peer connection", zap.Error(err))
	}
}

func (client *Client) OnICEConnectionStateChange(callback func(webrtc.ICEConnectionState)) {
	client.connection.OnICEConnectionStateChange(callback)
}

// NewPeerConnection returns webrtc-peer-connection with opus media-engine.
func (client *Client) ClosePeerConnection() {
	if err := client.connection.Close(); err != nil {
		client.Logger.Fatal("failed to close peer connection", zap.Error(err))
	}
}
