package client

type Config struct {
	AudioFileAddress string     `koanf:"audio-file-address"`
	Connection       Connection `koanf:"connection"`
}

type Connection struct {
	STUNServer string   `koanf:"stun-server"`
	RTPCodec   RTPCodec `koanf:"rtp-codec"`
}

type RTPCodec struct {
	ClockRate   uint32 `koanf:"clock-rate"`
	Channels    uint16 `koanf:"channels"`
	PayloadType uint8  `koanf:"payload-type"`
}
