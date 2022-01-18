package cli

import (
	"fmt"
	"strings"
)

const (
	vKeyFile = "pay.vkey"
	sKeyFile = "pay.skey"
	addrFile = "pay.addr"
)

func networkToUse(arg string) string {
	if arg == "mainnet" {
		return "--mainnet"
	}
	return fmt.Sprintf("--testnet-magic %s", arg)
}

func normaliseDirectory(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}
	return fmt.Sprintf("%v/", path)
}
