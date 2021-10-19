package client

type Config struct {
	AudioFileAddress string `koanf:"audio-file-address"`
	Connection       struct {
		STUNServer string `koanf:"stun-server"`
		RTPCodec   struct {
			ClockRate   uint32 `koanf:"clock-rate"`
			Channels    uint16 `koanf:"channels"`
			PayloadType uint8  `koanf:"payload-type"`
		} `koanf:"rtp-codec"`
	} `koanf:"connection"`
}
