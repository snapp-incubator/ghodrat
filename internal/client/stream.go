package client

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

func (client *Client) StreamAudioFile(
	connectedCtx context.Context,
	trackWriter func(sample media.Sample) error,
	doneChannel chan bool) {
	audioFileAddress := client.Config.AudioFileAddress

	_, err := os.Stat(audioFileAddress)
	if os.IsNotExist(err) {
		client.Logger.Fatal("audio file does not exist", logger.String("path", audioFileAddress))
	}

	file, err := os.Open(audioFileAddress)
	if err != nil {
		client.Logger.Fatal("failed to open audio file", logger.Error(err))
	}

	// Open on oggfile in non-checksum mode.
	ogg, _, err := oggreader.NewWith(file)
	if err != nil {
		client.Logger.Fatal("failed to read ogg audio", logger.Error(err))
	}

	// Wait for connection established
	<-connectedCtx.Done()

	// Keep track of last granule, the difference is the amount of samples in the buffer
	var lastGranule uint64

	for {
		pageData, pageHeader, err := ogg.ParseNextPage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				client.Logger.Info("all audio pages parsed and sent")
				os.Exit(0)
			}

			client.Logger.Fatal("failed to parse ogg", logger.Error(err))
		}

		// The amount of samples is the difference between the last and current timestamp
		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		// nolint: gomnd
		sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond
		// nolint: exhaustivestruct
		sample := media.Sample{Data: pageData, Duration: sampleDuration}

		if err = trackWriter(sample); err != nil {
			client.Logger.Fatal("failed to write media sample", logger.Error(err))
		}

		time.Sleep(sampleDuration)
	}
}

// OnTrack sets an event handler which is called when remote track arrives from a remote peer.
func (client *Client) OnTrack(callback func(*webrtc.TrackRemote)) {
	client.connection.OnTrack(
		func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
			callback(track)
		},
	)
}

// AddTrack adds a Track to the PeerConnection.
func (client *Client) AddTrack(track *webrtc.TrackLocalStaticSample) *webrtc.RTPSender {
	rtpSender, err := client.connection.AddTrack(track)
	if err != nil {
		client.Logger.Fatal("failed to create RTP sender", logger.Error(err))
	}

	return rtpSender
}
