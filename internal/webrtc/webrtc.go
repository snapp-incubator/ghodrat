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
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
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

	peerConnection, err := NewPeerConnectionWithOpusMediaEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}
	c.peerConnection = peerConnection

	audioBuilder := samplebuilder.New(10, &codecs.OpusPacket{}, 48000)
	c.audioBuilder = audioBuilder

	w, err := os.OpenFile(fmt.Sprintf("/tmp/test-%d.opus", rand.Intn(100)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}

	ws, err := webm.NewSimpleBlockWriter(w, []webm.TrackEntry{
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

	rtpSender, err := peerConnection.AddTrack(audioTrack)
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

// ReadRTCPPackets reads incoming RTCP packets
// Before these packets are returned they are processed by interceptors. For things
// like NACK this needs to be called.
func (c *Call) ReadRTCPPackets() error {
	rtcpBuf := make([]byte, 1500)
	for {
		if _, _, err := c.rtpSender.Read(rtcpBuf); err != nil {
			return err
		}
	}
}

func (c *Call) StreamAudioFile(audioFileAddress string) error {
	file, err := os.Open(audioFileAddress)
	if err != nil {
		return fmt.Errorf("failed to open audio file: %w", err)
	}

	// Open on oggfile in non-checksum mode.
	ogg, _, err := oggreader.NewWith(file)
	if err != nil {
		return fmt.Errorf("failed to read ogg audio: %w", err)
	}

	// Wait for connection established
	<-c.iceConnectedCtx.Done()

	// Keep track of last granule, the difference is the amount of samples in the buffer
	var lastGranule uint64
	for {
		pageData, pageHeader, err := ogg.ParseNextPage()
		if err == io.EOF {
			c.logger.Info("all audio pages parsed and sent")
			os.Exit(0)
		}
		if err != nil {
			return fmt.Errorf("failed to parse ogg: %w", err)
		}

		// The amount of samples is the difference between the last and current timestamp
		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

		if err = c.audioTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); err != nil {
			return fmt.Errorf("failed to write media sample: %w", err)
		}

		time.Sleep(sampleDuration)
	}
}

func (c *Call) CreateAndSetLocalOffer() error {
	offer, err := c.peerConnection.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("failed to create offer: %w", err)
	}

	if err = c.peerConnection.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("failed to set offer: %w", err)
	}

	return nil
}

func (c *Call) watchHandle(handle *janus.Handle) {
	for {
		msg := <-handle.Events
		switch msg := msg.(type) {
		case *janus.SlowLinkMsg:
			c.logger.Info("SlowLinkMsg", logger.Int("id", int(handle.ID)))
		case *janus.MediaMsg:
			c.logger.Info("MediaEvent", logger.String("type", msg.Type), logger.Bool("receiving", msg.Receiving))
		case *janus.WebRTCUpMsg:
			c.logger.Info("WebRTCUp", logger.Int("id", int(handle.ID)))
		case *janus.HangupMsg:
			c.logger.Info("HangupEvent", logger.Int("id", int(handle.ID)))
		case *janus.EventMsg:
			c.logger.Info("EventMsg", logger.Any("data", msg.Plugindata.Data))
		}
	}
}

func (c *Call) Call() error {

	create, err := c.audioBridgeHandle.Request(map[string]interface{}{
		"request": "create",
	})
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	roomID := create.PluginData.Data["room"].(float64)

	c.logger.Info("room created", logger.Float64("room", roomID))

	join, err := c.audioBridgeHandle.Message(map[string]interface{}{
		"request": "join",
		"room":    roomID,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	c.logger.Info("joined to room", logger.Float64("id", join.Plugindata.Data["id"].(float64)),
		logger.Any("participants", join.Plugindata.Data["participants"]))

	configure, err := c.audioBridgeHandle.Message(map[string]interface{}{
		"request": "configure",
	}, map[string]interface{}{
		"type": "offer",
		"sdp":  c.peerConnection.LocalDescription().SDP,
	})
	if err != nil {
		return fmt.Errorf("failed to send offer: %w", err)
	}

	c.logger.Info("offer has been sent")

	if configure.Jsep != nil {
		if err := c.peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  configure.Jsep["sdp"].(string),
		}); err != nil {
			return fmt.Errorf("failed to set SDP answer: %w", err)
		}
	}

	return nil
}

func (c *Call) Close() error {
	if err := c.peerConnection.Close(); err != nil {
		return fmt.Errorf("failed to close peer connection: %w", err)
	}

	if err := c.audioWriter.Close(); err != nil {
		return fmt.Errorf("failed to close audio writer: %w", err)
	}

	return nil
}
