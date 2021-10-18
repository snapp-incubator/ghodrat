package configs

import (
	"strings"

	"github.com/snapp-incubator/ghodrat/internal"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/server/janus"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
	"github.com/snapp-incubator/ghodrat/pkg/tracer"
	"github.com/snapp-incubator/ghodrat/pkg/utils"
)

var (
	envPrefix = strings.ToUpper(internal.Subsystem) + "_"
	filePath  = "./internal/configs/values.yml"
)

type Configs struct {
	Logger *logger.Config `koanf:"logger"`
	Tracer *tracer.Config `koanf:"tracer"`
	Client *client.Config `koanf:"client"`
	Janus  *janus.Config  `koanf:"janus"`
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
