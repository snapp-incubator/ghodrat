package ion

import sdk "github.com/pion/ion-sdk-go"

type CallMode int8

const (
	PubSub = iota
	PubOnly
	SubOnly
)

func (c *Client) Call(callMode CallMode) error {
	joinCfg := sdk.NewJoinConfig()

	if callMode == PubOnly {
		joinCfg.SetNoSubscribe()
	}

	if callMode == SubOnly {
		joinCfg.SetNoPublish()
	}

	if err := c.serverClient.Join(c.sessionID, joinCfg); err != nil {
		return err
	}

	return nil
}

func (c *Client) HangUp() {
	c.serverClient.Close()
}
