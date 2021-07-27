// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package user

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
	"github.com/romanornr/CryptoTwitterTipBot/insightjson"
	"github.com/viacoin/viad/txscript"
	"github.com/viacoin/viad/wire"
	btcutil "github.com/viacoin/viautil"
	"io/ioutil"
	"log"
	"net/http"
)

//type errorBroadcast string
const ErrNotEnoughBalance = "16: bad-txns-vout-negative. Code:-26"
const ErrNotEnoughFee = "66: insufficient priority. Code:-26"
const ErrTransactionTooSmall = "64: dust. Code:-26"
const ErrTxDecodeFailed = "Something seems wrong: TX decode failed. Code:-22"

//the cost in satoshi to create 1 output
func getCostForOneOutput(coinsymbol string) (satoshi int64) {
	coin, err := bcoins.SelectCoin(coinsymbol)
	if err != nil {
		log.Printf("error getting cost for one output for %s: %s\n", coinsymbol, err)
	}
	return coin.FeePerByte
}

//estimate fee in satoshi by knowing the amount of inputs required for the tx
//use the fee per byte for a coin and multiple by the estimated size
//return the required fee in satoshi
func feeEstimator(coinsymbol string, inputs int) (satoshi int64) {
	coin, err := bcoins.SelectCoin(coinsymbol)
	if err != nil {
		log.Printf("error getting cost for one output for %s: %s\n", coinsymbol, err)
	}

	feePerByte, defaultOutputs := coin.FeePerByte, 2 // 1 to receiver and 1 change back to source address
	estimatedSize := inputs*180 + 1*34 + 10 + defaultOutputs
	return feePerByte * int64(estimatedSize)
}

//prepare for allowing the user to create a payment to a destination address.
//this function will build a signed transaction & broadcast it
func (user User) PayTo(network bcoins.Network, destination string, amount float64) (insightjson.Txid, bcoins.Transaction, error) {
	amountSat, err := btcutil.NewAmount(amount)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("this amount does not seem to be in float64")
	}

	requiredUtxos := getMinimalRequiredUTXO(int64(amountSat), user.GetUnspentOutputs(network))
	fee := feeEstimator(network.Symbol, len(requiredUtxos))

	tx, err := user.BuildSignedTx(network, destination, int64(amountSat), fee)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, err
	}

	result, tx, err := BroadcastTransaction(network, tx)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, err
	}

	return result, tx, nil
}

// withdraw entire balance & use the correct fee. This function is similair to the PayTo() function
func (user User) Withdraw(network bcoins.Network, destination string) (insightjson.Txid, bcoins.Transaction, error) {
	fee := feeEstimator(network.Symbol, len(user.GetUnspentOutputs(network)))

	addressInfo, _ := user.AddressInfo(network)
	tx, err := user.BuildSignedTx(network, destination, addressInfo.BalanceSat-fee, fee)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, err
	}
	result, tx, err := BroadcastTransaction(network, tx)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, err
	}

	return result, tx, nil
}

//allow to recover when signing panics.
//this can happen with a wrong address for example
func recoverBuildSignedTx() {
	if r := recover(); r != nil {
		fmt.Sprintf("recovered from BuildSignedTx: %s\n", r)
	}
}

//build a signed transaction by using the Network, destination, amount and fee
//the required UTXO will be spent and the change will be send back
func (user User) BuildSignedTx(network bcoins.Network, destination string, amount int64, fee int64) (bcoins.Transaction, error) {

	defer recoverBuildSignedTx()

	var transaction bcoins.Transaction

	unspentOutputs := user.GetUnspentOutputs(network)

	wif := user.PrivateKeyWif(network)
	addresspubkey := user.PublicKey(network)

	tx := wire.NewMsgTx(wire.TxVersion)
	sourceUTXOs := getMinimalRequiredUTXO(amount+fee, unspentOutputs)
	availableAmountToSpend := int64(0) // amount in UTXO available

	for idx := range sourceUTXOs {
		availableAmountToSpend += sourceUTXOs[idx].Amount
		sourceUTXO := wire.NewOutPoint(sourceUTXOs[idx].Hash, sourceUTXOs[idx].TxIndex)
		sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
		tx.AddTxIn(sourceTxIn)
	}

	// create tx outs
	destinationAddress, err := btcutil.DecodeAddress(destination, network.GetNetworkParams())
	if err != nil {
		return bcoins.Transaction{}, fmt.Errorf("error decoding destination addresss %s\n", destinationAddress.String())
	}

	destinationScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return bcoins.Transaction{}, fmt.Errorf("destination script failed: %s\n", destinationScript)
	}

	// tx to sent to the destination
	destinationOutput := wire.NewTxOut(amount, destinationScript)
	tx.AddTxOut(destinationOutput)

	// calculate change
	change := availableAmountToSpend - amount
	change -= fee
	if change < 0 {
		err = fmt.Errorf("Amount to big too send. Change is: %d\n", change)
	}

	// change address, sent left over UTXO back to user his own address
	// if there's change left, sent it back to the source wallet
	// or send back when the change is equal or higher than the cost of making one extra output
	if change != 0 || change >= getCostForOneOutput(network.Symbol) {
		ChangeAddr := addresspubkey
		changeSendToScript, err := txscript.PayToAddrScript(ChangeAddr)
		if err != nil {
			return bcoins.Transaction{}, fmt.Errorf("error creating changeSendToScript: %s\n", err)
		}

		// tx out to sent back to user his own address
		changeOutput := wire.NewTxOut(change, changeSendToScript)
		tx.AddTxOut(changeOutput)
	}

	sourceAddress := addresspubkey

	sourcePKScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return bcoins.Transaction{}, fmt.Errorf("error creating sourcePkScript: %s\n", err)
	}

	for i := range sourceUTXOs {
		sigScript, err := txscript.SignatureScript(tx, i, sourcePKScript, txscript.SigHashAll, wif.PrivKey, true)
		if err != nil {
			return bcoins.Transaction{}, fmt.Errorf("error creating Signature script: %s\n", err)
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	var signedTx bytes.Buffer
	err = tx.Serialize(&signedTx)
	if err != nil {
		return bcoins.Transaction{}, fmt.Errorf("error serializing signed tx: %s\n", err)
	}
	transaction.TxId = tx.TxHash().String()
	transaction.Amount = amount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()
	transaction.Coinsymbol = network.Symbol
	transaction.CoinName = network.Name

	return transaction, err
}

//broadcast a signed transaction to the blockexplorer.
//the transaction will either be denied or rejected by the network
func BroadcastTransaction(network bcoins.Network, tx bcoins.Transaction) (insightjson.Txid, bcoins.Transaction, error) {
	jsonData := insightjson.InsightRawTx{Rawtx: tx.SignedTx}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("error broadcasting viacoin tx because jsonMarshal failed: %s\n", err)
	}

	insightExplorer, _ := GetInsightExplorer(network.Symbol)
	insightExplorerBroacastApi := fmt.Sprintf("%s/tx/send", insightExplorer.Api)

	response, err := http.Post(insightExplorerBroacastApi, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("Error broadcasting Viacoin transaction with blockexplorer: %s\n", err)
	}

	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("error reading response from viacoin blockexplorer broadcast: %s\n", err)
	}

	if response.StatusCode != 200 { // some error handling if broadcasting fails
		rejectReason := string(result)
		switch rejectReason {
		case ErrNotEnoughBalance:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("not enough balance to cover the transaction including the required fees")
		case ErrNotEnoughFee:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("fee needs to be higher")
		case ErrTransactionTooSmall:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("transaction too small (dust transaction)\nTx does not meet the minimal amount")
		case ErrTxDecodeFailed:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("transaction decode failed !\n Maybe a wrong address?")
		default:
			return insightjson.Txid{}, bcoins.Transaction{}, fmt.Errorf("%s\n", string(result))
		}
	}

	var txid = insightjson.Txid{}
	err = json.Unmarshal([]byte(result), &txid)
	if err != nil {
		return txid, bcoins.Transaction{}, fmt.Errorf("something went wrong with receiving your txid")
	}

	return txid, tx, nil
}
