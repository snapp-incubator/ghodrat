package client

import (
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

type Client struct {
	Config *Config
	Logger *zap.Logger

	connection *webrtc.PeerConnection
}
