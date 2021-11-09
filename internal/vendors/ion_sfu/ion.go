package ion_sfu

import (
	"math/rand"
	"strconv"

	"github.com/google/uuid"
	sdk "github.com/pion/ion-sdk-go"
	"github.com/pion/webrtc/v3"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"go.uber.org/zap"
)

type Engine struct {
	Logger *zap.Logger
	Config *Config

	engine *sdk.Engine
}

type Client struct {
	Logger     *zap.Logger
	peerClient *client.Client

	serverClient *sdk.Client

	sessionID string
	id        string
}

func NewEngine(cfg *Config, logger *zap.Logger) *Engine {
	sdkCfg := sdk.Config{
		WebRTC: sdk.WebRTCTransportConfig{
			Configuration: webrtc.Configuration{
				ICEServers: []webrtc.ICEServer{
					{URLs: cfg.StunServers},
				},
			},
		},
	}

	engine := sdk.NewEngine(sdkCfg)

	return &Engine{
		Logger: logger,
		Config: cfg,
		engine: engine,
	}
}

func (e *Engine) NewClient(peerClient *client.Client) (*Client, error) {
	cid := ""

	uuid, err := uuid.NewUUID()
	if err != nil {
		e.Logger.Error("failed to generate uuid", zap.Error(err))
		cid = strconv.FormatInt(rand.Int63(), 10)
	} else {
		cid = uuid.String()
	}

	c, err := sdk.NewClient(e.engine, e.Config.Address, "ghodrat_ion_"+cid)
	if err != nil {
		return nil, err
	}

	// Change c.OnTrack for custom packet processing,
	//default approach is described in c.Join()

	c.OnError = func(err error) {
		e.Logger.Error("ion client error", zap.String("cid", cid), zap.Error(err))
	}

	return &Client{
		Logger:       e.Logger,
		serverClient: c,
		peerClient:   peerClient,
		id:           cid,
	}, nil
}
