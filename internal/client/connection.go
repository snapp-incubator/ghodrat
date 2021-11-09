package client

import (
	"context"

	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

// InitiatePeerConnection returns webrtc-peer-connection with opus media-engine.
func (client *Client) CreatePeerConnection(iceConnectedCtxCancel context.CancelFunc) {
	var err error

	// A MediaEngine defines the codecs supported by a PeerConnection
	mediaEngine := &webrtc.MediaEngine{}

	// configuration of OPUS codec
	codec := webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     client.Config.RTPCodec.MimeType,
			ClockRate:    client.Config.RTPCodec.ClockRate,
			Channels:     client.Config.RTPCodec.Channels,
			SDPFmtpLine:  "minptime=10;useinbandfec=1",
			RTCPFeedback: nil,
		},
		PayloadType: webrtc.PayloadType(client.Config.RTPCodec.PayloadType),
	}

	rtpCodecType := webrtc.RTPCodecType(client.Config.RTPCodec.CodecType)
	if err = mediaEngine.RegisterCodec(codec, rtpCodecType); err != nil {
		client.Logger.Fatal("failed to register opus codec", zap.Error(err))
	}

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	// Prepare the configuration
	// nolint: exhaustivestruct
	config := webrtc.Configuration{
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{client.Config.STUNServer}},
		},
	}

	// Create a new RTCPeerConnection
	client.connection, err = api.NewPeerConnection(config)
	if err != nil {
		client.Logger.Fatal("failed to create peer connection", zap.Error(err))
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	client.connection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		client.Logger.Info("connection state has changed", zap.String("state", connectionState.String()))
		if connectionState == webrtc.ICEConnectionStateConnected {
			iceConnectedCtxCancel()
		}
	})
}

// NewPeerConnection returns webrtc-peer-connection with opus media-engine.
func (client *Client) ClosePeerConnection() {
	if err := client.connection.Close(); err != nil {
		client.Logger.Fatal("failed to close peer connection", zap.Error(err))
	}
}
