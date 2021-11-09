package client

import "github.com/pion/webrtc/v3"

func (client *Client) OnICEConnectionStateChange(callback func(webrtc.ICEConnectionState)) {
	client.connection.OnICEConnectionStateChange(callback)
}
