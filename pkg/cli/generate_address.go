package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

func init() {
	genAddressCmd.Flags().StringP("output-dir", "o", "", "Output directory to write to")
	genAddressCmd.MarkFlagRequired("output-dir")
	rootCmd.AddCommand(genAddressCmd)
}

var genAddressCmd = &cobra.Command{
	Use:   "gen-addr [mainnet or testnet number]",
	Short: "Generate a new payment address and output the files",
	Long: `Provide either the testnet magic number or mainnet as the first arg.
E.g. cardano-cli gen-addr mainnet OR cardano-cli gen-addr 1097911063`,
	Run:  generateAddressWrapper,
	Args: cobra.MinimumNArgs(1),
}

func generateAddressWrapper(cmd *cobra.Command, args []string) {
	rawDir := cmd.Flags().Lookup("output-dir").Value.String()
	outputDir := normaliseDirectory(rawDir)
	network := networkToUse(args[0])
	generateAddress(outputDir, network)
}

func generateAddress(outputDir string, network string) {
	vKeyFilePath := fmt.Sprintf("%v%v", outputDir, vKeyFile)
	sKeyFilePath := fmt.Sprintf("%v%v", outputDir, sKeyFile)
	err := generatePaymentKeyFiles(vKeyFilePath, sKeyFilePath)
	if err != nil {
		return
	}
	addrFilePath := fmt.Sprintf("%v%v", outputDir, addrFile)
	addrErr := generatePaymentAddrFile(vKeyFilePath, addrFilePath, network)
	if addrErr != nil {
		return
	}
}

//revive:disable:unhandled-error TODO: Probably should handle the errors
func generatePaymentAddrFile(vKeyFilePath string, addrFilePath string, network string) error {
	networks := strings.Split(network, " ")
	var stdout []byte
	var err error
	if len(networks) == 1 {
		stdout, err = commandRunner.ExecuteCommand("address", "build",
			"--payment-verification-key-file", vKeyFilePath, network, "--out-file", addrFilePath)
	} else {
		stdout, err = commandRunner.ExecuteCommand("address", "build",
			"--payment-verification-key-file", vKeyFilePath, networks[0], networks[1], "--out-file", addrFilePath)
	}
	if err != nil {
		fmt.Println(string(stdout))
		fmt.Println(err.Error())
		return err
	}
	addr, err := ioutil.ReadFile(addrFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(fmt.Sprintf("Generated address: %v", string(addr)))
	fmt.Println(fmt.Sprintf("Stored address at: %v", addrFilePath))
	return nil
}

func generatePaymentKeyFiles(vKeyFilePath string, sKeyFilePath string) error {
	stdout, err := commandRunner.ExecuteCommand("address", "key-gen",
		"--verification-key-file", vKeyFilePath, "--signing-key-file", sKeyFilePath)
	if err != nil {
		fmt.Println(string(stdout))
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(fmt.Sprintf("Generated files: %v, %v", vKeyFilePath, sKeyFilePath))
	return nil
}

//revive:enable:unhandled-error
