package client

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/ivfreader"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
)

var (
	audioTrackCodecCapability = webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}
	audioTrackCodecId         = "audio"

	videoTrackCodecCapability = webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}
	videoTrackCodecId         = "video"
)

func (client *Client) ReadTrack(doneChannel chan bool, connectedCtx context.Context) {
	_, err := os.Stat(client.Config.TrackAddress)
	if os.IsNotExist(err) {
		panic("Track Not Exists")
	}

	var trackCodecCapability webrtc.RTPCodecCapability
	var trackCodecId string

	mimeType := client.Config.RTPCodec.MimeType
	isAudioTrack := strings.Split(mimeType, "/")[0] == "audio"

	if isAudioTrack {
		trackCodecCapability = audioTrackCodecCapability
		trackCodecId = audioTrackCodecId
	} else {
		trackCodecCapability = videoTrackCodecCapability
		trackCodecId = videoTrackCodecId
	}

	track, trackErr := webrtc.NewTrackLocalStaticSample(trackCodecCapability, trackCodecId, "ghodrat")
	if trackErr != nil {
		panic(trackErr)
	}

	rtpSender, trackErr := client.connection.AddTrack(track)
	if trackErr != nil {
		panic(trackErr)
	}

	go readRTCP(rtpSender)

	if isAudioTrack {
		go audioTrack(client.Config.TrackAddress, doneChannel, track, connectedCtx)
	} else {
		go videoTrack(client.Config.TrackAddress, doneChannel, track, connectedCtx)
	}
}

// Read incoming RTCP packets
// Before these packets are returned they are processed by interceptors. For things
// like NACK this needs to be called.
func readRTCP(rtpSender *webrtc.RTPSender) {
	rtcpBuf := make([]byte, 1500)
	for {
		if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
			return
		}
	}
}

const oggPageDuration = time.Millisecond * 20

func audioTrack(address string, doneChannel chan bool, track *webrtc.TrackLocalStaticSample, connectedCtx context.Context) {
	// Open a OGG file and start reading using our OGGReader
	file, oggErr := os.Open(address)
	if oggErr != nil {
		panic(oggErr)
	}

	// Open on oggfile in non-checksum mode.
	ogg, _, oggErr := oggreader.NewWith(file)
	if oggErr != nil {
		panic(oggErr)
	}

	// Wait for connection established
	<-connectedCtx.Done()

	// Keep track of last granule, the difference is the amount of samples in the buffer
	var lastGranule uint64

	// It is important to use a time.Ticker instead of time.Sleep because
	// * avoids accumulating skew, just calling time.Sleep didn't compensate for the time spent parsing the data
	// * works around latency issues with Sleep (see https://github.com/golang/go/issues/44343)
	ticker := time.NewTicker(oggPageDuration)
	for ; true; <-ticker.C {
		pageData, pageHeader, oggErr := ogg.ParseNextPage()
		if oggErr == io.EOF {
			doneChannel <- true
			return
		}

		if oggErr != nil {
			panic(oggErr)
		}

		// The amount of samples is the difference between the last and current timestamp
		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

		sample := media.Sample{Data: pageData, Duration: sampleDuration}
		if oggErr = track.WriteSample(sample); oggErr != nil {
			panic(oggErr)
		}
	}
}

func videoTrack(address string, doneChannel chan bool, track *webrtc.TrackLocalStaticSample, connectedCtx context.Context) {
	// Open a IVF file and start reading using our IVFReader
	file, ivfErr := os.Open(address)
	if ivfErr != nil {
		panic(ivfErr)
	}

	ivf, header, ivfErr := ivfreader.NewWith(file)
	if ivfErr != nil {
		panic(ivfErr)
	}

	// Wait for connection established
	<-connectedCtx.Done()

	// Send our video file frame at a time. Pace our sending so we send it at the same speed it should be played back as.
	// This isn't required since the video is timestamped, but we will such much higher loss if we send all at once.
	//
	// It is important to use a time.Ticker instead of time.Sleep because
	// * avoids accumulating skew, just calling time.Sleep didn't compensate for the time spent parsing the data
	// * works around latency issues with Sleep (see https://github.com/golang/go/issues/44343)
	ticker := time.NewTicker(time.Millisecond * time.Duration((float32(header.TimebaseNumerator)/float32(header.TimebaseDenominator))*1000))
	for ; true; <-ticker.C {
		frame, _, ivfErr := ivf.ParseNextFrame()
		if ivfErr == io.EOF {
			doneChannel <- true
			return
		}

		if ivfErr != nil {
			panic(ivfErr)
		}

		sample := media.Sample{Data: frame, Duration: time.Second}
		if ivfErr = track.WriteSample(sample); ivfErr != nil {
			panic(ivfErr)
		}
	}
}
