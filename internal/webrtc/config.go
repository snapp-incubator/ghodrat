package webrtc

type Config struct {
	AudioFileAddress string
	Connection       struct {
		STUNServers []string
		RTPCodec    struct {
			ClockRate   uint32
			Channels    uint16
			PayloadType uint8
		}
	}
}
