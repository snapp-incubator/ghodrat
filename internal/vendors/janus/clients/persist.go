package clients

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/at-wat/ebml-go/webm"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3/pkg/media/samplebuilder"

	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

type AudioFactory struct {
	audioWriter    webm.BlockWriteCloser
	audioBuilder   *samplebuilder.SampleBuilder
	audioTimestamp time.Duration
}

func NewAudioFactory(cfg *Config) (*AudioFactory, error) {
	factory := new(AudioFactory)
	factory.audioBuilder = samplebuilder.New(cfg.AudioMaxLate, &codecs.OpusPacket{}, cfg.AudioSampleRate)

	file, err := os.CreateTemp(os.TempDir(), "ghodrat-*.opus")
	if err != nil {
		return nil, err
	}

	ws, err := webm.NewSimpleBlockWriter(file, []webm.TrackEntry{
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
		return nil, err
	}

	factory.audioWriter = ws[0]

	return factory, nil
}

func (client *Client) SaveOpusTrack(track *webrtc.TrackRemote) {
	client.Logger.Info(
		"track has started",
		zap.Int("payload_type", int(track.PayloadType())),
		zap.String("mime_type", track.Codec().RTPCodecCapability.MimeType),
	)

	for {
		// Read RTP packets being sent to Pion
		rtp, _, err := track.ReadRTP()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			client.Logger.Fatal("failed to read RTP", zap.Error(err))
		}

		switch track.Kind() {
		case webrtc.RTPCodecTypeAudio:
			client.AudioFactory.audioBuilder.Push(rtp)

			for {
				sample := client.AudioFactory.audioBuilder.Pop()
				if sample == nil {
					break
				}

				if client.AudioFactory.audioWriter != nil {
					client.AudioFactory.audioTimestamp += sample.Duration
					if _, err := client.AudioFactory.audioWriter.Write(
						true,
						int64(client.AudioFactory.audioTimestamp/time.Millisecond),
						sample.Data); err != nil {
						client.Logger.Fatal("failed to write audio", zap.Error(err))
					}
				}
			}
		default:
			client.Logger.Fatal("only audio type is supported")
		}
	}
}

func (client *Client) CloseOpusTrack() {
	if err := client.AudioFactory.audioWriter.Close(); err != nil {
		client.Logger.Fatal("failed to close audio writer", zap.Error(err))
	}
}
