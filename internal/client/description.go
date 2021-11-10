package client

import (
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

func (client *Client) CreateAndSetOffer() {
	offer, err := client.connection.CreateOffer(nil)
	if err != nil {
		client.Logger.Fatal("failed to create SDP offer", zap.Error(err))
	}

	if err = client.connection.SetLocalDescription(offer); err != nil {
		client.Logger.Fatal("failed to set SDP offer", zap.Error(err))
	}
}

func (client *Client) CreateAndSetAnswer() {
	answer, err := client.connection.CreateAnswer(nil)
	if err != nil {
		client.Logger.Fatal("failed to create SDP answer", zap.Error(err))
	}

	if err = client.connection.SetLocalDescription(answer); err != nil {
		client.Logger.Fatal("failed to set SDP answer", zap.Error(err))
	}
}

func (client *Client) GetLocalDescription() *webrtc.SessionDescription {
	return client.connection.LocalDescription()
}

func (client *Client) SetRemoteDescription(sdp webrtc.SessionDescription) {
	client.Logger.Info("remote description", zap.Any("sdp", sdp))
	if err := client.connection.SetRemoteDescription(sdp); err != nil {
		client.Logger.Fatal("failed to set remote SDP answer", zap.Error(err))
	}
}
