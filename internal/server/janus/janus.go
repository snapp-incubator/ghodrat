package janus

import (
	"context"
	"errors"
	"io"

	"github.com/notedit/janus-go"
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"go.uber.org/zap"
)

type Janus struct {
	Logger *zap.Logger
	Client *client.Client
	Config *Config

	rtpSender *webrtc.RTPSender

	audioTrack        *webrtc.TrackLocalStaticSample
	audioBridgeHandle *janus.Handle

	iceConnectedCtx       context.Context
	iceConnectedCtxCancel context.CancelFunc
}

func (j *Janus) initiate() {
	j.Client.CreatePeerConnection()

	j.iceConnectedCtx, j.iceConnectedCtxCancel = context.WithCancel(context.Background())

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	j.Client.OnICEConnectionStateChange(j.onICEConnectionStateChange)

	// Set a handler for when a new remote track starts, this handler copies inbound RTP packets,
	// replaces the SSRC and sends them back
	j.Client.OnTrack(j.Client.SaveOpusTrack)

	var err error

	// Create a audio track
	j.audioTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion")
	if err != nil {
		j.Logger.Fatal("failed to create audio track", zap.Error(err))
	}

	j.rtpSender = j.Client.AddTrack(j.audioTrack)

	gateway, err := janus.Connect(j.Config.Address)
	if err != nil {
		j.Logger.Fatal("failed to connect to janus", zap.Error(err))
	}

	session, err := gateway.Create()
	if err != nil {
		j.Logger.Fatal("failed to create session", zap.Error(err))
	}

	handle, err := session.Attach("janus.plugin.audiobridge")
	if err != nil {
		j.Logger.Fatal("failed to create handle", zap.Error(err))
	}

	j.audioBridgeHandle = handle

	go j.watchHandle(j.audioBridgeHandle)
}

// readRTCPPackets reads incoming RTCP packets
// Before these packets are returned they are processed by interceptors. For things
// like NACK this needs to be called.
func (j *Janus) readRTCPPackets() {
	const bufferSize = 1500

	rtcpBuf := make([]byte, bufferSize)

	for {
		if _, _, err := j.rtpSender.Read(rtcpBuf); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			j.Logger.Error("failed to read rtcp packets", zap.Error(err))
		}
	}
}

func (j *Janus) onICEConnectionStateChange(connectionState webrtc.ICEConnectionState) {
	j.Logger.Info("connection state has changed", zap.String("state", connectionState.String()))

	if connectionState == webrtc.ICEConnectionStateConnected {
		j.iceConnectedCtxCancel()
	}
}

func (j *Janus) watchHandle(handle *janus.Handle) {
	for {
		msg := <-handle.Events
		switch msg := msg.(type) {
		case *janus.SlowLinkMsg:
			j.Logger.Info("SlowLinkMsg", zap.Int("id", int(handle.ID)))
		case *janus.MediaMsg:
			j.Logger.Info("MediaEvent", zap.String("type", msg.Type), zap.Bool("receiving", msg.Receiving))
		case *janus.WebRTCUpMsg:
			j.Logger.Info("WebRTCUp", zap.Int("id", int(handle.ID)))
		case *janus.HangupMsg:
			j.Logger.Info("HangupEvent", zap.Int("id", int(handle.ID)))
		case *janus.EventMsg:
			j.Logger.Info("EventMsg", zap.Any("data", msg.Plugindata.Data))
		}
	}
}
