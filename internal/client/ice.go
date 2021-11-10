package client

import (
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

func (client *Client) OnIceCandidate(callback func(candidate *webrtc.ICECandidate)) {
	client.connection.OnICECandidate(callback)
}

func (client *Client) AddIceCandidate(c *webrtc.ICECandidateInit) {
	err := client.connection.AddICECandidate(*c)
	if err != nil {
		client.Logger.Fatal("failed to add ice candidate", zap.Error(err))
	}
}

func (client *Client) onICEConnectionStateChange(connectionState webrtc.ICEConnectionState) {
	client.Logger.Info("connection state has changed", zap.String("state", connectionState.String()))
	if connectionState == webrtc.ICEConnectionStateConnected {
		client.iceConnectedCtxCancel()
	}
}
