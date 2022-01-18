package main

import (
	"github.com/lambda-honeypot/cardano-tool/cardano-cli/pkg/cmdrunner"
	"github.com/lambda-honeypot/cardano-tool/cardano-cli/pkg/sendfunds"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetLevel(log.DebugLevel)
	runner := cmdrunner.CmdRunner{Command: "cardano-cli"}
	fs := sendfunds.NewFundSender(runner, "--mainnet", "")
	//token1 := "1815bee29d9d1eabf78b7f21f29ae55cbad8d06fa470a65ddbf98156.HONEY"
	//token2 := "45ace7db4aec426e119445e867816f31cdebc014b4f642fc1decda41.HONEYChristmas"
	startAddress := ""
	paymentAddressesWithTokens := map[string]sendfunds.PaymentDetails{}

	signingKeyFile := os.Getenv("SIGNING_KEY_FILE")
	if signingKeyFile == "" {
		log.Fatalf("Error SIGNING_KEY_FILE variable is not set")
	}

	bal, _ := fs.CreateUTXOFromAddress(startAddress)
	log.Infof("Balance Before: %d", bal.ADABalance)
	for idx, balance := range bal.TokenBalances {
		log.Infof("Token Before: %s + %d", idx, balance)
	}
	err := fs.PayMultiple(startAddress, signingKeyFile, paymentAddressesWithTokens)
	if err != nil {
		log.Fatalf("Failed to pay multiple wallets: %v", err)
	}
	newbal, _ := fs.CreateUTXOFromAddress(startAddress)
	log.Infof("Balance After: %+v\n", newbal)
}
