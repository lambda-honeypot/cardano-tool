package main

import (
	"github.com/lambda-honeypot/cardano-tool/cardano-cli/pkg/cmdrunner"
	"github.com/lambda-honeypot/cardano-tool/cardano-cli/pkg/sendfunds"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	runner := cmdrunner.CmdRunner{Command: "cardano-cli"}
	fs := sendfunds.NewFundSender(runner, "--mainnet", "")

	paymentAddress := "addr1q8q566cvhawynjmw008u5xlzkqaplx33vjhs82ec7f2vzt7m9dtqxjj5kv4u40r5ss7dsy679zcw9xkm07kasdg6u4hs3azrhh"
	stakeAddress := "stake1uymzamdz7228jw96mm56vzqw2famneqh200mnxthqyr2grqza3dkg"
	bal, _ := fs.CreateUTXOFromAddress(paymentAddress)
	log.Infof("Payment balance Before: %d", bal.ADABalance)
	for idx, balance := range bal.TokenBalances {
		log.Infof("Token Before: %s + %d", idx, balance)
	}
	err := fs.GenerateWithdrawTxFile(stakeAddress, paymentAddress)
	if err != nil {
		log.Fatalf("Failed to generate withdraw tx: %v", err)
	}
}
