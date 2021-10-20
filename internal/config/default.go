package config

import (
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/server/janus"
	"github.com/snapp-incubator/ghodrat/internal/tracer"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Logger: &logger.Config{
			Development: true,
			Encoding:    "console",
			Level:       "warn",
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
			Connection: client.Connection{
				STUNServer: "stun:stun.l.google.com:19302",
				RTPCodec: client.RTPCodec{
					ClockRate:   48000,
					Channels:    2,
					PayloadType: 111,
				},
			},
		},

		Janus: &janus.Config{
			Address:    "ws://localhost:8080",
			MaxLate:    10,
			SampleRate: 48000,
		},
	}
}
