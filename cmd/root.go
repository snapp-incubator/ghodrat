package main

import (
	"fmt"

	"github.com/snapp-incubator/ghodrat/cmd/janus"
	"github.com/spf13/cobra"
)

const (
	errExecuteCMD = "failed to execute root command"

	short = "WebRTC stress testing tool"
	long  = `ghodrat is a CMD tool used to stress test janus WebRTC media servers`
)

func main() {
	cmd := &cobra.Command{Short: short, Long: long}
	cmd.AddCommand(janus.Command())

	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		panic(map[string]interface{}{"err": err, "msg": errExecuteCMD})
	}
}
