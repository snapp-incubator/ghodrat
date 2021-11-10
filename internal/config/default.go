package config

import (
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/tracer"
	"github.com/snapp-incubator/ghodrat/internal/vendors/ion_sfu"
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
			STUNServer: "stun:stun.l.google.com:19302",

			// OPUS audio (ogg)
			TrackAddress: "./static/audio.ogg",
			RTPCodec: &client.RTPCodec{
				MimeType:    "audio/opus",
				ClockRate:   48000,
				Channels:    2,
				PayloadType: 111,
				CodecType:   1,
			},

			// VP8 video (ivf)
			// TrackAddress: "./static/video.ivf",
			// RTPCodec: &client.RTPCodec{
			// 	MimeType:    "video/vp8",
			// 	ClockRate:   90000,
			// 	Channels:    2,
			// 	PayloadType: 96,
			// 	CodecType:   2,
			// },
		},
		Janus: &janus.Config{
			Address: "ws://janus-dispatching-testing.apps.private.okd4.teh-1.snappcloud.io",
		},
		Ion_sfu: &ion_sfu.Config{
			Address: "ws://localhost:7000/ws",
		},
	}
}
