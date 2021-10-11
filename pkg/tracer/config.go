package tracer

type Config struct {
	Enabled    bool    `koanf:"enabled"`
	Host       string  `koanf:"host"`
	Port       int     `koanf:"port"`
	SampleRate float64 `koanf:"sample-rate"`
}
