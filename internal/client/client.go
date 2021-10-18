package client

import (
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

type Client struct {
	Config *Config
	Logger logger.Logger

	connection *webrtc.PeerConnection
}
