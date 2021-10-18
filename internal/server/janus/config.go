package janus

type Config struct {
	Address    string `koanf:"address"`
	MaxLate    uint16 `koanf:"max-late"`
	SampleRate uint32 `koanf:"sample-rate"`
}
