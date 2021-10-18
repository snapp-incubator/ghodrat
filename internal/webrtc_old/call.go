package webrtc

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/at-wat/ebml-go/webm"
	"github.com/notedit/janus-go"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/samplebuilder"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

type Call struct {
	logger logger.Logger

	peerConnection *webrtc.PeerConnection

	audioBridgeHandle *janus.Handle

	audioTrack *webrtc.TrackLocalStaticSample
	rtpSender  *webrtc.RTPSender

	audioWriter    webm.BlockWriteCloser
	audioBuilder   *samplebuilder.SampleBuilder
	audioTimestamp time.Duration

	iceConnectedCtx       context.Context
	iceConnectedCtxCancel context.CancelFunc
}

func NewCall(janusAddress string, logger logger.Logger) (*Call, error) {
	c := &Call{logger: logger}

	var err error
	c.peerConnection, err = NewPeerConnectionWithOpusMediaEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	c.audioBuilder = samplebuilder.New(10, &codecs.OpusPacket{}, 48000)

	path := fmt.Sprintf("/tmp/test-%d.opus", rand.Intn(100))
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
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
		return nil, fmt.Errorf("failed to create block write: %w", err)
	}

	c.audioWriter = ws[0]

	iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())
	c.iceConnectedCtx = iceConnectedCtx
	c.iceConnectedCtxCancel = iceConnectedCtxCancel

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	c.peerConnection.OnICEConnectionStateChange(c.onICEConnectionStateChange)

	// Set a handler for when a new remote track starts, this handler copies inbound RTP packets,
	// replaces the SSRC and sends them back
	c.peerConnection.OnTrack(c.saveOpusTrack)

	// Create a audio track
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion")
	if err != nil {
		return nil, fmt.Errorf("failed to create audio track: %w", err)
	}

	rtpSender, err := c.peerConnection.AddTrack(audioTrack)
	if err != nil {
		return nil, fmt.Errorf("failed to create RTP sender: %w", err)
	}
	c.rtpSender = rtpSender

	gateway, err := janus.Connect(janusAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to janus: %w", err)
	}

	session, err := gateway.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	handle, err := session.Attach("janus.plugin.audiobridge")
	if err != nil {
		return nil, fmt.Errorf("failed to create handle: %w", err)
	}

	c.audioBridgeHandle = handle

	go c.watchHandle(c.audioBridgeHandle)

	return c, nil
}

func (c *Call) onICEConnectionStateChange(connectionState webrtc.ICEConnectionState) {
	c.logger.Info("connection state has changed", logger.String("state", connectionState.String()))
	if connectionState == webrtc.ICEConnectionStateConnected {
		c.iceConnectedCtxCancel()
	}
}

func (c *Call) saveOpusTrack(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	c.logger.Info(
		"track has started",
		logger.Int("payload_type", int(track.PayloadType())),
		logger.String("mime_type", track.Codec().RTPCodecCapability.MimeType),
	)

	for {
		// Read RTP packets being sent to Pion
		rtp, _, err := track.ReadRTP()
		if err != nil {
			if err == io.EOF {
				return
			}
			c.logger.Fatal("failed to read RTP", logger.Error(err))
		}
		switch track.Kind() {
		case webrtc.RTPCodecTypeAudio:
			c.audioBuilder.Push(rtp)
			for {
				sample := c.audioBuilder.Pop()
				if sample == nil {
					break
				}
				if c.audioWriter != nil {
					c.audioTimestamp += sample.Duration
					if _, err := c.audioWriter.Write(true, int64(c.audioTimestamp/time.Millisecond), sample.Data); err != nil {
						c.logger.Fatal("failed to write audio", logger.Error(err))
					}
				}
			}
		default:
			c.logger.Fatal("only audio type is supported")
		}
	}
}

func (c *Call) watchHandle(handle *janus.Handle) {
	for {
		msg := <-handle.Events
		switch msg := msg.(type) {
		case *janus.SlowLinkMsg:
			c.logger.Info("SlowLinkMsg",
				logger.Int("id", int(handle.ID)),
			)
		case *janus.MediaMsg:
			c.logger.Info("MediaEvent",
				logger.String("type", msg.Type),
				logger.Bool("receiving", msg.Receiving),
			)
		case *janus.WebRTCUpMsg:
			c.logger.Info("WebRTCUp",
				logger.Int("id", int(handle.ID)),
			)
		case *janus.HangupMsg:
			c.logger.Info("HangupEvent",
				logger.Int("id", int(handle.ID)),
			)
		case *janus.EventMsg:
			c.logger.Info("EventMsg",
				logger.Any("data", msg.Plugindata.Data),
			)
		}
	}
}
