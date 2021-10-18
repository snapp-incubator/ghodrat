package webrtc

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

// ReadRTCPPackets reads incoming RTCP packets
// Before these packets are returned they are processed by interceptors. For things
// like NACK this needs to be called.
func (c *Call) ReadRTCPPackets() error {
	rtcpBuf := make([]byte, 1500)
	for {
		if _, _, err := c.rtpSender.Read(rtcpBuf); err != nil {
			return err
		}
	}
}

func (c *Call) StreamAudioFile(audioFileAddress string) error {
	file, err := os.Open(audioFileAddress)
	if err != nil {
		return fmt.Errorf("failed to open audio file: %w", err)
	}

	// Open on oggfile in non-checksum mode.
	ogg, _, err := oggreader.NewWith(file)
	if err != nil {
		return fmt.Errorf("failed to read ogg audio: %w", err)
	}

	// Wait for connection established
	<-c.iceConnectedCtx.Done()

	// Keep track of last granule, the difference is the amount of samples in the buffer
	var lastGranule uint64
	for {
		pageData, pageHeader, err := ogg.ParseNextPage()
		if err == io.EOF {
			c.logger.Info("all audio pages parsed and sent")
			os.Exit(0)
		}
		if err != nil {
			return fmt.Errorf("failed to parse ogg: %w", err)
		}

		// The amount of samples is the difference between the last and current timestamp
		sampleCount := float64(pageHeader.GranulePosition - lastGranule)
		lastGranule = pageHeader.GranulePosition
		sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

		if err = c.audioTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); err != nil {
			return fmt.Errorf("failed to write media sample: %w", err)
		}

		time.Sleep(sampleDuration)
	}
}

func (c *Call) CreateAndSetLocalOffer() error {
	offer, err := c.peerConnection.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("failed to create offer: %w", err)
	}

	if err = c.peerConnection.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("failed to set offer: %w", err)
	}

	return nil
}

func (c *Call) Call() error {

	create, err := c.audioBridgeHandle.Request(map[string]interface{}{
		"request": "create",
	})
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	roomID := create.PluginData.Data["room"].(float64)

	c.logger.Info("room created", logger.Float64("room", roomID))

	join, err := c.audioBridgeHandle.Message(map[string]interface{}{
		"request": "join",
		"room":    roomID,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	c.logger.Info("joined to room", logger.Float64("id", join.Plugindata.Data["id"].(float64)),
		logger.Any("participants", join.Plugindata.Data["participants"]))

	configure, err := c.audioBridgeHandle.Message(map[string]interface{}{
		"request": "configure",
	}, map[string]interface{}{
		"type": "offer",
		"sdp":  c.peerConnection.LocalDescription().SDP,
	})
	if err != nil {
		return fmt.Errorf("failed to send offer: %w", err)
	}

	c.logger.Info("offer has been sent")

	if configure.Jsep != nil {
		if err := c.peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  configure.Jsep["sdp"].(string),
		}); err != nil {
			return fmt.Errorf("failed to set SDP answer: %w", err)
		}
	}

	return nil
}

func (c *Call) Close() error {
	if err := c.peerConnection.Close(); err != nil {
		return fmt.Errorf("failed to close peer connection: %w", err)
	}

	if err := c.audioWriter.Close(); err != nil {
		return fmt.Errorf("failed to close audio writer: %w", err)
	}

	return nil
}
