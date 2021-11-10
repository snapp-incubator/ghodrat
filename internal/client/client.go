package client

import (
	"context"

	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

type Client struct {
	Config *Config
	Logger *zap.Logger

	iceConnectedCtxCancel context.CancelFunc
	connection            *webrtc.PeerConnection
}

type Config struct {
	TrackAddress string    `koanf:"track-address"`
	STUNServer   string    `koanf:"stun-server"`
	RTPCodec     *RTPCodec `koanf:"rtp-codec"`
}

type RTPCodec struct {
	MimeType    string `koanf:"mime-type"`
	ClockRate   uint32 `koanf:"clock-rate"`
	Channels    uint16 `koanf:"channels"`
	PayloadType uint8  `koanf:"payload-type"`
	CodecType   uint8  `koanf:"codec-type"`
}
