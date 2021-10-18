package webrtc

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

func (manager *Manager) StreamAudioFile() error {
	audioFileAddress := manager.Config.AudioFileAddress

	_, err := os.Stat(audioFileAddress)
	if os.IsNotExist(err) {
		manager.Logger.Fatal("audio file does not exist", logger.String("path", audioFileAddress))
	}

	file, err := os.Open(audioFileAddress)
	if err != nil {
		manager.Logger.Fatal("failed to open audio file", logger.Error(err))
	}

	// Open on oggfile in non-checksum mode.
	ogg, _, err := oggreader.NewWith(file)
	if err != nil {
		manager.Logger.Fatal("failed to read ogg audio", logger.Error(err))
	}

	// Wait for connection established
	<-manager.iceConnectedCtx.Done()

	// Keep track of last granule, the difference is the amount of samples in the buffer
	var lastGranule uint64
	for {
		pageData, pageHeader, err := ogg.ParseNextPage()
		if err != nil {
			if err == io.EOF {
				manager.Logger.Info("all audio pages parsed and sent")
				os.Exit(0)
			}
			manager.Logger.Fatal("failed to parse ogg", logger.Error(err))
		}

		// The amount of samples is the difference between the last and current timestamp
		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

		if err = manager.audioTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); err != nil {
			return fmt.Errorf("failed to write media sample: %w", err)
		}

		time.Sleep(sampleDuration)
	}
}
