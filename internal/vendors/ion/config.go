package ion

type Config struct {
	Address     string   `koanf:"address"`
	StunServers []string `koanf:"stun-servers"`
	Session     string   `koanf:"session"`
}
