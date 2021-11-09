package client

import (
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

type Client struct {
	Config *Config
	Logger *zap.Logger

	AudioFactory *AudioFactory

	connection *webrtc.PeerConnection
}

type Config struct {
	AudioFileAddress string `koanf:"audio-file-address"`
	AudioMaxLate     uint16 `koanf:"audio-max-late"`
	AudioSampleRate  uint32 `koanf:"sample-rate"`

	STUNServer string `koanf:"stun-server"`
	RTPCodec   struct {
		MimeType    string `koanf:"mime-type"`
		ClockRate   uint32 `koanf:"clock-rate"`
		Channels    uint16 `koanf:"channels"`
		PayloadType uint8  `koanf:"payload-type"`
		CodecType   uint8  `koanf:"codec-type"`
	} `koanf:"rtp-codec"`
}
