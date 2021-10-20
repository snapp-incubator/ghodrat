package configs

import (
	"github.com/snapp-incubator/ghodrat/internal"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/server/janus"
	"github.com/snapp-incubator/ghodrat/internal/tracer"
	"github.com/snapp-incubator/ghodrat/pkg/utils"
)

var (
	envPrefix = internal.Subsystem + "_"
	filePath  = "./internal/configs/values.yml"
)

type Configs struct {
	Logger    *logger.Config `koanf:"logger"`
	Tracer    *tracer.Config `koanf:"tracer"`
	CallCount int            `koanf:"call-count"`
	Client    *client.Config `koanf:"client"`
	Janus     *janus.Config  `koanf:"janus"`
}

func Load(environment string) *Configs {
	configs := new(Configs)

	var source utils.Source

	if environment == "prod" {
		source = utils.Env
	} else {
		source = utils.File
	}

	utils.Configs{Source: source, EnvPrefix: envPrefix, FilePath: filePath}.
		Load(configs)

	return configs
}
