package sendfunds

import (
	"encoding/json"
	"fmt"
	"github.com/lambda-honeypot/cardano-tool/cardano-cli/pkg/cli"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// FundSender some struct
type FundSender struct {
	network string
	magic   string
	runner  cli.Runner
}

// NewFundSender init method for fundsender
func NewFundSender(runner cli.Runner, network, magic string) FundSender {
	return FundSender{network, magic, runner}
}

// RewardQuery struct to unmarshall json
type RewardQuery struct {
	RewardAccountBalance int
}

// TipQuery struct to unmarshall json
type TipQuery struct {
	Slot int
}

func (fs *FundSender) CreateUTXOFromAddress(tokenAddress string) (*FullUTXO, error) {
	queryReturn, err := fs.runner.ExecuteCommand("query", "utxo", "--address", tokenAddress, fs.network, fs.magic)
	if err != nil {
		return nil, fmt.Errorf("stdin: %s stderr: %v", queryReturn, err)
	}
	log.Debugf("UTXO query result\n%s", string(queryReturn))
	fullUTXO, err := parseFullUTXO(string(queryReturn), tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse utxo query: %s", string(queryReturn))
	}

	return fullUTXO, nil
}

func (fs *FundSender) GenerateWithdrawTxFile(rewardAddress string, paymentAddress string) error {
	dir, err := ioutil.TempDir("", "gen_withdraw")
	if err != nil {
		return err
	}
	tmpFile := dir + "/tx.tmp"
	rawFile := "tx.raw"
	paramsFile := dir + "/" + strings.ReplaceAll(fmt.Sprintf("%s-params.json", fs.network), "--", "")

	defer os.RemoveAll(dir)
	slot, _ := fs.getCurrentSlot()
	err = fs.createParamsFile(paramsFile)
	if err != nil {
		return fmt.Errorf("failed to create params file: %v", err)
	}
	rewardBalance, err := fs.getRewardBalance(rewardAddress)
	if err != nil {
		return fmt.Errorf("failed to get reward balance for %s, with: %v", rewardAddress, err)
	}
	withdrawalString := fmt.Sprintf("%s+%d", rewardAddress, rewardBalance)
	utxoDetails, err := fs.CreateUTXOFromAddress(paymentAddress)
	if err != nil {
		return fmt.Errorf("failed to create utxofrom address: %s with error: %v", rewardAddress, err)
	}
	if err != nil {
		return fmt.Errorf("failed to generate token details with: %v", err)
	}
	paymentDetails := map[string]PaymentDetails{}
	if err != nil {
		return fmt.Errorf("failed to generate token details with: %v", err)
	}
	err = fs.createRawTxFile(utxoDetails, paymentAddress, tmpFile, paymentDetails, slot, 0, 0, []TokenDetails{}, withdrawalString)
	if err != nil {
		return fmt.Errorf("failed to create tmp tx file for fee calc: %v", err)
	}
	minFee, err := fs.calculateMinimumFee(utxoDetails, paymentDetails, tmpFile, paramsFile, "2")
	if err != nil {
		return fmt.Errorf("failed to calculate min fee: %v", err)
	}
	log.Infof("calculate min fee: %d", minFee)
	txOutAdaAmount := utxoDetails.ADABalance - minFee + rewardBalance
	err = fs.createRawTxFile(utxoDetails, paymentAddress, rawFile, paymentDetails, slot, txOutAdaAmount, minFee, []TokenDetails{}, withdrawalString)
	if err != nil {
		return fmt.Errorf("failed to create raw tx file for payment: %v", err)
	}
	return nil
}

func (fs *FundSender) PayMultiple(sourceAddress, signingKeyFile string, paymentDetails map[string]PaymentDetails) error {
	dir, err := ioutil.TempDir("", "pay_multi")
	if err != nil {
		return err
	}
	tmpFile := dir + "/tx.tmp"
	rawFile := dir + "/tx.raw"
	txSignedFile := dir + "/tx.signed"
	paramsFile := dir + "/" + strings.ReplaceAll(fmt.Sprintf("%s-params.json", fs.network), "--", "")

	defer os.RemoveAll(dir)
	slot, _ := fs.getCurrentSlot()
	err = fs.createParamsFile(paramsFile)
	if err != nil {
		return fmt.Errorf("failed to create params file: %v", err)
	}
	utxoDetails, err := fs.CreateUTXOFromAddress(sourceAddress)
	if err != nil {
		return fmt.Errorf("failed to create utxofrom address: %s with error: %v", sourceAddress, err)
	}
	paymentTokenDetails, err := generateTokenDetailsAndVerify(utxoDetails, paymentDetails)
	if err != nil {
		return fmt.Errorf("failed to generate token details with: %v", err)
	}
	err = fs.createRawTxFile(utxoDetails, sourceAddress, tmpFile, paymentDetails, slot, 0, 0, []TokenDetails{}, "")
	if err != nil {
		return fmt.Errorf("failed to create tmp tx file for fee calc: %v", err)
	}
	minFee, err := fs.calculateMinimumFee(utxoDetails, paymentDetails, tmpFile, paramsFile, "1")
	if err != nil {
		return fmt.Errorf("failed to calculate min fee: %v", err)
	}
	log.Infof("calculate min fee: %d", minFee)
	totalADAinLovelace := 0
	for _, paymentDetail := range paymentDetails {
		totalADAinLovelace += paymentDetail.AdaAmount
	}
	txOutAdaAmount := utxoDetails.ADABalance - totalADAinLovelace - minFee
	err = fs.createRawTxFile(utxoDetails, sourceAddress, rawFile, paymentDetails, slot, txOutAdaAmount, minFee, paymentTokenDetails, "")
	if err != nil {
		return fmt.Errorf("failed to create raw tx file for payment: %v", err)
	}
	err = fs.signTransactionFile(rawFile, signingKeyFile, txSignedFile)
	if err != nil {
		return fmt.Errorf("failed to sign tx file for send: %v", err)
	}
	err = fs.sendTransaction(txSignedFile)
	if err != nil {
		return fmt.Errorf("failed to send signed tx: %v", err)
	}
	return nil
}

func (fs *FundSender) calculateMinimumFee(utxo *FullUTXO, paymentAddresses map[string]PaymentDetails, tempFile, paramsFile, witnessCount string) (int, error) {
	transactionOutCount := len(paymentAddresses) + 1
	minFeeArgs := fs.calculateMinFeeArgs(paramsFile, tempFile, witnessCount, utxo.TXCount, transactionOutCount)
	log.Debugf("%s", minFeeArgs)
	minFeeReturn, err := fs.runner.ExecuteCommand(minFeeArgs...)
	log.Debugf("MIN FEE: %s", minFeeReturn)
	if err != nil {
		return 0, fmt.Errorf("stdin: %s stderr: %v", minFeeReturn, err)
	}
	minFeeSplit := strings.Fields(string(minFeeReturn))
	minFee, err := strconv.Atoi(minFeeSplit[0])
	if err != nil {
		return 0, err
	}
	return minFee, nil
}

func (fs *FundSender) createParamsFile(paramsFile string) error {
	queryProtocolArgs := fs.queryProtocolParamsArgs(paramsFile)
	log.Debugf("%s", queryProtocolArgs)
	queryProtocolReturn, err := fs.runner.ExecuteCommand(queryProtocolArgs...)
	log.Debugf("%s", string(queryProtocolReturn))
	if err != nil {
		return fmt.Errorf("stdin: %s stderr: %v", queryProtocolReturn, err)
	}
	return nil
}

func (fs *FundSender) createRawTxFile(utxo *FullUTXO, sourceAddress, outFile string, paymentDetails map[string]PaymentDetails, currentSlot, txOutAdaAmount, minFee int, txOutTokenAmounts []TokenDetails, withdrawalString string) error {
	rawTxArgs := buildRawTransactionArgs(utxo, sourceAddress, outFile, currentSlot, txOutAdaAmount, minFee, paymentDetails, txOutTokenAmounts, withdrawalString)
	log.Debugf("%s", rawTxArgs)
	buildRawReturn, err := fs.runner.ExecuteCommand(rawTxArgs...)
	log.Debugf("%s", string(buildRawReturn))
	if err != nil {
		return fmt.Errorf("stdin: %s stderr: %v", buildRawReturn, err)
	}
	return nil
}

func generateTokenDetailsAndVerify(utxo *FullUTXO, paymentDetails map[string]PaymentDetails) ([]TokenDetails, error) {
	sendTotals := make(map[string]int)
	var returnTokens []TokenDetails
	for _, paymentDetail := range paymentDetails {
		for _, tokenDetail := range paymentDetail.PaymentTokens {
			sendTotals[tokenDetail.TokenID] += tokenDetail.TokenAmount
		}
	}

	sortedKeys := make([]string, 0, len(sendTotals))
	for k := range sendTotals {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	for _, tokenID := range sortedKeys {
		sendTokenAmount := sendTotals[tokenID]
		if utxo.TokenBalances[tokenID] < sendTokenAmount {
			return nil, fmt.Errorf("total send token amount for token: %s is %d - this is greater than source wallet balance of %d", tokenID, sendTokenAmount, utxo.TokenBalances[tokenID])
		}
		adjustedAmount := utxo.TokenBalances[tokenID] - sendTokenAmount
		returnTokens = append(returnTokens, TokenDetails{TokenID: tokenID, TokenAmount: adjustedAmount})
	}
	return returnTokens, nil
}

func (fs *FundSender) getCurrentSlot() (int, error) {
	var tipQuery TipQuery
	jsQuery, err := fs.runner.ExecuteCommand("query", "tip", fs.network, fs.magic)
	log.Debugf("Query tip: %s", string(jsQuery))
	if err != nil {
		return 0, fmt.Errorf("stdin: %s stderr: %v", jsQuery, err)
	}
	err = json.Unmarshal(jsQuery, &tipQuery)
	if err != nil {
		return 0, err
	}
	return tipQuery.Slot, nil
}

func (fs *FundSender) getRewardBalance(stakeAddress string) (int, error) {
	var rewardQueries []RewardQuery
	jsQuery, err := fs.runner.ExecuteCommand("query", "stake-address-info", fs.network, fs.magic, "--address", stakeAddress)
	log.Debugf("Query Reward: %s", string(jsQuery))
	if err != nil {
		return 0, fmt.Errorf("stdin: %s stderr: %v", jsQuery, err)
	}
	err = json.Unmarshal(jsQuery, &rewardQueries)
	if err != nil {
		return 0, err
	}
	return rewardQueries[0].RewardAccountBalance, nil
}

func (fs *FundSender) sendTransaction(txSignedFile string) error {
	txSubmitReturn, err := fs.runner.ExecuteCommand("transaction", "submit", "--tx-file", txSignedFile, fs.network, fs.magic)
	log.Infof("%s", string(txSubmitReturn))
	if err != nil {
		return fmt.Errorf("stdin: %s stderr: %v", txSubmitReturn, err)
	}
	return nil
}

func (fs *FundSender) signTransactionFile(txRawFile, signingKeyFile, txSignedFile string) error {
	txSignReturn, err := fs.runner.ExecuteCommand("transaction", "sign", "--tx-body-file", txRawFile, "--signing-key-file", signingKeyFile, fs.network, fs.magic, "--out-file", txSignedFile)
	log.Infof("%s", string(txSignReturn))
	if err != nil {
		return fmt.Errorf("stdin: %s stderr: %v", txSignReturn, err)
	}
	return nil
}
