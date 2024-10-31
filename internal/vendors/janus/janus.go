package janus

import (
	"fmt"

	"github.com/notedit/janus-go"
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"go.uber.org/zap"
)

type Janus struct {
	Logger *zap.Logger
	Client *client.Client
	Config *Config

	audioBridgeHandle *janus.Handle
}

func (j *Janus) initiate() {
	gateway, err := janus.Connect(j.Config.Address)
	if err != nil {
		j.Logger.Fatal("failed to connect to janus", zap.Error(err))
	}

	session, err := gateway.Create()
	if err != nil {
		j.Logger.Fatal("failed to create session", zap.Error(err))
	}

	j.audioBridgeHandle, err = session.Attach("janus.plugin.audiobridge")
	if err != nil {
		j.Logger.Fatal("failed to create handle", zap.Error(err))
	}
}

func (j *Janus) handle() {
	handle := j.audioBridgeHandle

	for {
		msg := <-handle.Events
		switch msg := msg.(type) {
		case *janus.SlowLinkMsg:
			j.Logger.Info("SlowLinkMsg", zap.Int("id", int(handle.ID)))
		case *janus.MediaMsg:
			j.Logger.Info("MediaEvent", zap.String("type", msg.Type), zap.Bool("receiving", msg.Receiving))
		case *janus.WebRTCUpMsg:
			j.Logger.Info("WebRTCUp", zap.Int("id", int(handle.ID)))
		case *janus.HangupMsg:
			j.Logger.Info("HangupEvent", zap.Int("id", int(handle.ID)))
		case *janus.EventMsg:
			j.Logger.Info("EventMsg", zap.Any("data", msg.Plugindata.Data))
		}
	}
}

func (j *Janus) call() error {
	request := map[string]interface{}{"request": "create"}

	create, err := j.audioBridgeHandle.Request(request)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	roomID := create.PluginData.Data["room"].(float64)

	j.Logger.Info("room created", zap.Float64("room", roomID))

	body := map[string]interface{}{"request": "join", "room": roomID}

	join, err := j.audioBridgeHandle.Message(body, nil)
	if err != nil {
		j.Logger.Fatal("failed to join room", zap.Error(err))
	}

	j.Logger.Info("joined to room", zap.Float64(
		"id", join.Plugindata.Data["id"].(float64)),
		zap.Any("participants", join.Plugindata.Data["participants"]),
	)

	body = map[string]interface{}{"request": "configure"}
	jsep := map[string]interface{}{
		"type": "offer",
		"sdp":  j.Client.GetLocalDescription().SDP,
	}

	configure, err := j.audioBridgeHandle.Message(body, jsep)
	if err != nil {
		j.Logger.Fatal("failed to send offer", zap.Error(err))
	}

	j.Logger.Info("local description", zap.Any("sdp", jsep))

	if configure.Jsep == nil {
		j.Logger.Fatal("Jsep should not be nil")
	}

	if configure.Jsep != nil {
		remoteSDP, ok := configure.Jsep["sdp"].(string)
		if !ok {
			j.Logger.Fatal("Jsep should contain SDP")
		}

		description := webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: remoteSDP}
		j.Client.SetRemoteDescription(description)
	}

	return nil
}
