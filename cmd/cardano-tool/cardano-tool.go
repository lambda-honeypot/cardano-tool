package main

import (
	"github.com/lambda-honeypot/cardano-tool/cardano-cli/pkg/cli"
	"github.com/lambda-honeypot/cardano-tool/cardano-cli/pkg/cmdrunner"
	"os"
)

func main() {
	cmd := cli.NewCardanoTool(cmdrunner.CmdRunner{Command: "cardano-cli"})
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
