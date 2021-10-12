package configs

import (
	"strings"

	internal "github.com/snapp-incubator/ghodrat/new_internal"
	"github.com/snapp-incubator/ghodrat/pkg/utils"
)

var (
	envPrefix = strings.ToUpper(internal.Subsystem) + "_"
	filePath  = "./internal/configs/values.yml"
)

type Janus struct{}

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
