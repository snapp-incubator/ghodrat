package configs

import (
	"strings"

	"github.com/snapp-incubator/ghodrat/internal"
	"github.com/snapp-incubator/ghodrat/pkg/logger"
	"github.com/snapp-incubator/ghodrat/pkg/tracer"
	"github.com/snapp-incubator/ghodrat/pkg/utils"
)

var (
	envPrefix = strings.ToUpper(internal.Subsystem) + "_"
	filePath  = "./internal/configs/values.yml"
)

type Janus struct {
	Logger *logger.Config `koanf:"logger"`
	Tracer *tracer.Config `koanf:"tracer"`
}

func LoadJanus(environment string) *Janus {
	configs := new(Janus)
	load(environment, configs)
	return configs
}

func load(environment string, configs interface{}) {
	var source utils.Source

	if environment == "prod" {
		source = utils.Env
	} else {
		source = utils.File
	}

	utils.Configs{Source: source, EnvPrefix: envPrefix, FilePath: filePath}.
		Load(configs)
}
