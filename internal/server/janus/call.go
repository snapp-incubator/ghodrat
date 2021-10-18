package janus

import (
	"fmt"

	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

func (j *Janus) call() error {
	create, err := j.audioBridgeHandle.Request(map[string]interface{}{
		"request": "create",
	})
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	roomID := create.PluginData.Data["room"].(float64)

	j.Logger.Info("room created", logger.Float64("room", roomID))

	join, err := j.audioBridgeHandle.Message(map[string]interface{}{
		"request": "join",
		"room":    roomID,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	j.Logger.Info("joined to room", logger.Float64("id", join.Plugindata.Data["id"].(float64)),
		logger.Any("participants", join.Plugindata.Data["participants"]))

	configure, err := j.audioBridgeHandle.Message(map[string]interface{}{
		"request": "configure",
	}, map[string]interface{}{
		"type": "offer",
		"sdp":  j.Client.GetLocalDescription().SDP,
	})
	if err != nil {
		return fmt.Errorf("failed to send offer: %w", err)
	}

	j.Logger.Info("offer has been sent")

	if configure.Jsep != nil {
		j.Client.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  configure.Jsep["sdp"].(string),
		})
	}

	return nil
}
