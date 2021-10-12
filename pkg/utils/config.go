package utils

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Source uint8

const (
	File Source = iota
	Env
)

type Configs struct {
	Source    Source
	EnvPrefix string
	FilePath  string
}

const (
	delimeter = "."
	seperator = "_"
)

func (configs Configs) Load(object interface{}) {
	var provider koanf.Provider
	var parser koanf.Parser

	switch configs.Source {
	case Env:
		provider, parser = envConfig(configs.EnvPrefix)
	default:
		provider, parser = fileConfig(configs.FilePath)
	}

	k := koanf.New(delimeter)

	if err := k.Load(provider, parser); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", object); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}
}

func envConfig(prefix string) (koanf.Provider, koanf.Parser) {
	callback := func(source string) string {
		base := strings.ToLower(strings.TrimPrefix(source, prefix))
		return strings.ReplaceAll(base, seperator, delimeter)
	}
	return env.Provider(prefix, delimeter, callback), nil
}

func fileConfig(path string) (koanf.Provider, koanf.Parser) {
	return file.Provider(path), yaml.Parser()
}
