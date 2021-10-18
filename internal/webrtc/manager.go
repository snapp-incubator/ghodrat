package webrtc

import (
	"context"

	"github.com/at-wat/ebml-go/webm"
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/internal/media_server"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

type Manager struct {
	Config      *Config
	Logger      logger.Logger
	MediaServer media_server.MediaServer

	connection *webrtc.PeerConnection

	iceConnectedCtx context.Context
	// iceConnectedCtxCancel context.CancelFunc

	audioTrack  *webrtc.TrackLocalStaticSample
	audioWriter webm.BlockWriteCloser
	// audioBuilder   *samplebuilder.SampleBuilder
	// audioTimestamp time.Duration

	rtpSender *webrtc.RTPSender
}

func (manager Manager) Manage() {
	go func() {
		manager.ReadRTCPPackets()
	}()

	manager.CreateAndSetLocalOffer()

	if err := manager.MediaServer.Call(); err != nil {
		manager.Logger.Fatal("failed to call media-server", logger.Error(err))
	}
}

// ReadRTCPPackets reads incoming RTCP packets
// Before these packets are returned they are processed by interceptors. For things
// like NACK this needs to be called.
func (manager *Manager) ReadRTCPPackets() {
	rtcpBuf := make([]byte, 1500)
	for {
		if _, _, err := manager.rtpSender.Read(rtcpBuf); err != nil {
			manager.Logger.Error("failed to read RTCP packets", logger.Error(err))
		}
	}
}

func (manager *Manager) CreateAndSetLocalOffer() {
	offer, err := manager.connection.CreateOffer(nil)

	if err != nil {
		manager.Logger.Fatal("failed to create local SDP offer", logger.Error(err))
	}

	if err = manager.connection.SetLocalDescription(offer); err != nil {
		manager.Logger.Fatal("failed to set local SDP offer", logger.Error(err))
	}
}

func (manager *Manager) Close() {
	if err := manager.connection.Close(); err != nil {
		manager.Logger.Fatal("failed to close peer connection", logger.Error(err))
	}

	if err := manager.audioWriter.Close(); err != nil {
		manager.Logger.Fatal("failed to close audio writer", logger.Error(err))
	}
}
