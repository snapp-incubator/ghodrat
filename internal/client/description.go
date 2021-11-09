package client

import (
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

func (client *Client) CreateAndSetLocalOffer() {
	offer, err := client.connection.CreateOffer(nil)
	if err != nil {
		client.Logger.Fatal("failed to create local SDP offer", zap.Error(err))
	}

	if err = client.connection.SetLocalDescription(offer); err != nil {
		client.Logger.Fatal("failed to set local SDP offer", zap.Error(err))
	}
}

func (client *Client) GetLocalDescription() *webrtc.SessionDescription {
	return client.connection.LocalDescription()
}

func (client *Client) SetRemoteDescription(sdp webrtc.SessionDescription) {
	if err := client.connection.SetRemoteDescription(sdp); err != nil {
		client.Logger.Fatal("failed to set remote SDP answer", zap.Error(err))
	}
}
