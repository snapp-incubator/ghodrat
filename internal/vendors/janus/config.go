package janus

type Config struct {
	Address    string `koanf:"address"`
	MaxLate    int    `koanf:"max-late"`
	SampleRate int    `koanf:"sample-rate"`
}
