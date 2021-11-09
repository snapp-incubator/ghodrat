package config

import (
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/tracer"
	"github.com/snapp-incubator/ghodrat/internal/vendors/janus"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Logger: &logger.Config{
			Development: true,
			Encoding:    "console",
			Level:       "info",
		},
		Tracer: &tracer.Config{
			Enabled:    false,
			Host:       "localhost",
			Port:       6831,
			SampleRate: 0.1,
		},
		CallCount: 1,

		Client: &client.Config{
			AudioFileAddress: "./static/audio.ogg",
			AudioMaxLate:     10,
			AudioSampleRate:  48000,
			STUNServer:       "stun:stun.l.google.com:19302",
			RTPCodec: &client.RTPCodec{
				MimeType:    "",
				ClockRate:   2,
				Channels:    2,
				PayloadType: 1,
				CodecType:   1,
			},
		},

		Janus: &janus.Config{
			Address:    "ws://localhost:8080",
			MaxLate:    10,
			SampleRate: 48000,
		},
	}
}
