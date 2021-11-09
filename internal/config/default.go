package config

import (
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/tracer"
	"github.com/snapp-incubator/ghodrat/internal/vendors/janus/clients"
	janus_server "github.com/snapp-incubator/ghodrat/internal/vendors/janus/server"
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

		Client: &clients.Config{
			AudioFileAddress: "./static/audio.ogg",
			AudioMaxLate:     10,
			AudioSampleRate:  48000,
			Connection: clients.Connection{
				STUNServer: "stun:stun.l.google.com:19302",
				RTPCodec: clients.RTPCodec{
					ClockRate:   48000,
					Channels:    2,
					PayloadType: 111,
				},
			},
		},

		Janus: &janus_server.Config{
			Address:    "ws://localhost:8080",
			MaxLate:    10,
			SampleRate: 48000,
		},
	}
}
