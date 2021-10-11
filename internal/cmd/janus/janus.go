package janus

import "github.com/spf13/cobra"

func Register(root *cobra.Command) {
	janusCMD := &cobra.Command{
		Use:     "janus",
		Short:   "Stress test Janus media server",
		Example: "ghodrat janus --address ws://127.0.0.1:8188/ --audio-address ./audio.ogg --call-count 5",
		Run:     run,
	}

	janusCMD.Flags().String("address", "ws://127.0.0.1:8188/", "Janus media server websocket address")
	janusCMD.Flags().String("audio-file", "./static/audio.ogg", "audio file used to stream to Janus")
	janusCMD.Flags().Uint("call-count", 1, "number of concurrent calls")

	root.AddCommand(janusCMD)
}
