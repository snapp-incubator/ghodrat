package client

import (
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
)

func (manager *Client) CreateAndSetLocalOffer() {
	offer, err := manager.connection.CreateOffer(nil)
	if err != nil {
		manager.Logger.Fatal("failed to create local SDP offer", logger.Error(err))
	}

	if err = manager.connection.SetLocalDescription(offer); err != nil {
		manager.Logger.Fatal("failed to set local SDP offer", logger.Error(err))
	}
}

func (manager *Client) GetLocalDescription() *webrtc.SessionDescription {
	return manager.connection.LocalDescription()
}

func (manager *Client) SetRemoteDescription(sdp webrtc.SessionDescription) {
	if err := manager.connection.SetRemoteDescription(sdp); err != nil {
		manager.Logger.Fatal("failed to set remote SDP answer", logger.Error(err))
	}
}
