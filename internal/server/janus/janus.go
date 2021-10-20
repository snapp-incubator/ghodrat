package janus

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/at-wat/ebml-go/webm"
	"github.com/notedit/janus-go"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/samplebuilder"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"go.uber.org/zap"
)

type Janus struct {
	Logger zap.Logger
	Client *client.Client
	Config *Config

	rtpSender *webrtc.RTPSender

	audioTrack        *webrtc.TrackLocalStaticSample
	audioBridgeHandle *janus.Handle
	audioWriter       webm.BlockWriteCloser
	audioBuilder      *samplebuilder.SampleBuilder
	audioTimestamp    time.Duration

	iceConnectedCtx       context.Context
	iceConnectedCtxCancel context.CancelFunc
}

func (j *Janus) initiate() {
	j.Client.CreatePeerConnection()

	j.audioBuilder = samplebuilder.New(j.Config.MaxLate, &codecs.OpusPacket{}, j.Config.SampleRate)

	path := fmt.Sprintf("/tmp/test-%d.opus", rand.Intn(100))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		j.Logger.Fatal("failed to open audio file", zap.Error(err))
	}

	ws, err := webm.NewSimpleBlockWriter(file, []webm.TrackEntry{
		{
			Name:            "Audio",
			TrackNumber:     1,
			TrackUID:        12345,
			CodecID:         "A_OPUS",
			TrackType:       2,
			DefaultDuration: 20000000,
			Audio: &webm.Audio{
				SamplingFrequency: 48000.0,
				Channels:          2,
			},
		},
	})

	if err != nil {
		j.Logger.Fatal("failed to create block write", zap.Error(err))
	}

	j.audioWriter = ws[0]

	j.iceConnectedCtx, j.iceConnectedCtxCancel = context.WithCancel(context.Background())

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	j.Client.OnICEConnectionStateChange(j.onICEConnectionStateChange)

	// Set a handler for when a new remote track starts, this handler copies inbound RTP packets,
	// replaces the SSRC and sends them back
	j.Client.OnTrack(j.saveOpusTrack)

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
	rtcpBuf := make([]byte, 1500)
	for {
		if _, _, err := j.rtpSender.Read(rtcpBuf); err != nil {
			j.Logger.Error("failed to close audio writer", zap.Error(err))
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
