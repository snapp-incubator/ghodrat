package janus

import (
	"io"
	"time"

	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

func (j *Janus) saveOpusTrack(track *webrtc.TrackRemote) {
	j.Logger.Info(
		"track has started",
		zap.Int("payload_type", int(track.PayloadType())),
		zap.String("mime_type", track.Codec().RTPCodecCapability.MimeType),
	)

	for {
		// Read RTP packets being sent to Pion
		rtp, _, err := track.ReadRTP()
		if err != nil {
			if err == io.EOF {
				return
			}
			j.Logger.Fatal("failed to read RTP", zap.Error(err))
		}
		switch track.Kind() {
		case webrtc.RTPCodecTypeAudio:
			j.audioBuilder.Push(rtp)
			for {
				sample := j.audioBuilder.Pop()
				if sample == nil {
					break
				}
				if j.audioWriter != nil {
					j.audioTimestamp += sample.Duration
					if _, err := j.audioWriter.Write(true, int64(j.audioTimestamp/time.Millisecond), sample.Data); err != nil {
						j.Logger.Fatal("failed to write audio", zap.Error(err))
					}
				}
			}
		default:
			j.Logger.Fatal("only audio type is supported")
		}
	}
}
