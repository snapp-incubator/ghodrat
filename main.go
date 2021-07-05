package main

import (
	"context"
	"github.com/notedit/janus-go"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggreader"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal("failed to create peer connection", zap.Error(err))
	}

	defer peerConnection.Close()

	iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Println("connection state has changed", zap.String("state", connectionState.String()))
		if connectionState == webrtc.ICEConnectionStateConnected {
			iceConnectedCtxCancel()
		}
	})

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion")
	if err != nil {
		log.Fatal("failed to create audio track", zap.Error(err))
	}

	rtpSender, err := peerConnection.AddTrack(audioTrack)
	if err != nil {
		log.Fatal("failed to add audio track", zap.Error(err))
	}

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, err := rtpSender.Read(rtcpBuf); err != nil {
				return
			}
		}
	}()

	go func() {
		file, err := os.Open("./audio.ogg")
		if err != nil {
			log.Fatal("failed to open audio file", zap.Error(err))
		}

		ogg, _, err := oggreader.NewWith(file)
		if err != nil {
			log.Fatal("failed to read ogg audio", zap.Error(err))
		}

		<-iceConnectedCtx.Done()

		var lastGranule uint64
		for {
			pageData, pageHeader, err := ogg.ParseNextPage()
			if err == io.EOF {
				log.Println("all audio pages parsed and sent")
				os.Exit(0)
			}

			if err != nil {
				log.Fatal("failed to parse ogg", zap.Error(err))
			}
			sampleCount := float64(pageHeader.GranulePosition - lastGranule)
			lastGranule = pageHeader.GranulePosition
			sampleDuration := time.Duration((sampleCount/48000)*1000) * time.Millisecond

			if err = audioTrack.WriteSample(media.Sample{Data: pageData, Duration: sampleDuration}); err != nil {
				log.Fatal("failed to write media sample", zap.Error(err))
			}

			time.Sleep(sampleDuration)
		}
	}()

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Fatal("failed to create offer", zap.Error(err))
	}

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		log.Fatal("failed to set offer", zap.Error(err))
	}

	gateway, err := janus.Connect("ws://172.16.76.243:8188")
	if err != nil {
		log.Fatal("failed to connect to janus", zap.Error(err))
	}

	session, err := gateway.Create()
	if err != nil {
		log.Fatal("failed to create session", zap.Error(err))
	}

	handle, err := session.Attach("janus.plugin.audiobridge")
	if err != nil {
		log.Fatal("failed to create handle", zap.Error(err))
	}

	_, err = handle.Request(map[string]interface{}{
		"request": "create",
		"room":    12345,
	})
	if err != nil {
		log.Println("failed to create room", zap.Error(err))
	}

	//roomID := create.PluginData.Data["room"].(float64)
	roomID := 12345

	log.Println("room created", zap.Float64("room", float64(roomID)))

	join, err := handle.Message(map[string]interface{}{
		"request": "join",
		"room":    roomID,
	}, nil)
	if err != nil {
		log.Fatal("failed to join room", zap.Error(err))
	}

	log.Println("joined to room", zap.Float64("id", join.Plugindata.Data["id"].(float64)),
		zap.Any("participants", join.Plugindata.Data["participants"]))

	configure, err := handle.Message(map[string]interface{}{
		"request": "configure",
	}, map[string]interface{}{
		"type": "offer",
		"sdp":  peerConnection.LocalDescription().SDP,
	})
	if err != nil {
		log.Fatal("failed to send offer", zap.Error(err))
	}

	log.Println("offer has been sent")

	if configure.Jsep != nil {
		if err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  configure.Jsep["sdp"].(string),
		}); err != nil {
			log.Fatal("failed to set SDP answer", zap.Error(err))
		}
	}

	select {}
}
