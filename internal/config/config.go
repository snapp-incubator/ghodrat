package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/snapp-incubator/ghodrat/internal/client"
	"github.com/snapp-incubator/ghodrat/internal/logger"
	"github.com/snapp-incubator/ghodrat/internal/tracer"
	ion_sfu "github.com/snapp-incubator/ghodrat/internal/vendors/ion-sfu"
	"github.com/snapp-incubator/ghodrat/internal/vendors/janus"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "ghodrat_"
)

type Config struct {
	Logger    *logger.Config  `koanf:"logger"`
	Tracer    *tracer.Config  `koanf:"tracer"`
	CallCount int             `koanf:"call-count"`
	Client    *client.Config  `koanf:"client"`
	Janus     *janus.Config   `koanf:"janus"`
	IonSfu    *ion_sfu.Config `koanf:"ion-sfu"`
}

// New reads configuration with viper.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "_", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	log.Printf("following configuration is loaded:\n%+v", instance)

	return instance
}
