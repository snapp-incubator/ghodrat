package janus

import (
	"fmt"

	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

func (j *Janus) call() error {
	create, err := j.audioBridgeHandle.Request(map[string]interface{}{
		"request": "create",
	})
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	roomID := create.PluginData.Data["room"].(float64)

	j.Logger.Info("room created", zap.Float64("room", roomID))

	body := map[string]interface{}{"request": "join", "room": roomID}
	join, err := j.audioBridgeHandle.Message(body, nil)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	j.Logger.Info("joined to room", zap.Float64("id", join.Plugindata.Data["id"].(float64)),
		zap.Any("participants", join.Plugindata.Data["participants"]))

	body = map[string]interface{}{"request": "configure"}
	jsep := map[string]interface{}{"type": "offer",
		"sdp": j.Client.GetLocalDescription().SDP}
	configure, err := j.audioBridgeHandle.Message(body, jsep)
	if err != nil {
		return fmt.Errorf("failed to send offer: %w", err)
	}

	j.Logger.Info("offer has been sent", zap.Any("sdp", jsep))

	if configure.Jsep != nil {
		j.Client.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  configure.Jsep["sdp"].(string),
		})
	}

	return nil
}
