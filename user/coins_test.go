// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package user

//
//import (
//	"fmt"
//	"github.com/romanornr/CryptoTwitterTipBot/config"
//	"strings"
//	"testing"
//)
//
//var user1 = User{
//	1,
//	22,
//	int(0),
//	"@RNR_0",
//	"twitter",
//	33223,
//}
//
//func init() {
//	config.GetViperConfig()
//}
//
//func TestUser_PrivateKey(t *testing.T) {
//	key := user1.PrivateKey()
//	expectedKey := "xprv9virnR9sbxn3dDrfuz6fEqYYSUoauF3p2vgLrGjG822Nr439ukQfKtRbwVyL3p3GPsVN4oztBcEvJAntcFYJefrqGx9CrPDnGnEVprCbcZz"
//
//	if key.String() != expectedKey {
//		t.Errorf("wrong private key! expected: %s got: %s\n", expectedKey, key)
//	}
//}
//
//func TestUser_PrivateKeyWif(t *testing.T) {
//	network, _ := SelectNetwork("via")
//	key := user1.PrivateKeyWif(network)
//	if !strings.HasPrefix(key.String(), "W") {
//		t.Errorf("error creating private key in decompressed WIF format. WIF did not start with 7.\nGot: %s\n", key.String())
//	}
//}
//
//func TestUser_PublicKey(t *testing.T) {
//	network, _ := SelectNetwork("via")
//	key := user1.PublicKey(network)
//		fmt.Println(key)
//		expectedKey := "Vw9EhnS6aH2eArtyFEVGy7Y417PWGbPuVF"
//
//		if key.String() != expectedKey {
//			t.Errorf("wrong public key! expected: %s got: %s\n", expectedKey, key)
//		}
//}
//
//func TestUser_AddressInfo(t *testing.T) {
//	network, _ := SelectNetwork("via")
//	addressInfo, _ := user1.AddressInfo(network)
//	fmt.Printf("Balance for %s is: %f\n", user1.Username, addressInfo.Balance)
//}
//
//func TestUser_BuildSignedTx(t *testing.T) {
//		network, _ := SelectNetwork("via")
//		amount := int64(33000000)
//		requiredUtxos := getMinimalRequiredUTXO(amount, user1.GetUnspentOutputs(network))
//		fee := feeEstimator(network.Symbol, len(requiredUtxos))
//		tx, err := user1.BuildSignedTx(network,"Vd1fB3phUKKmV2FXoXLpCyTnugJ3KkpVFE", amount, fee)
//		if err != nil {
//			t.Errorf("%s", err)
//		}
//		if tx.SignedTx == "" {
//			t.Errorf("failed")
//		}
//		//fmt.Println(tx)
//}

//func TestUser_PayTo(t *testing.T) {
//	result, err := user1.ViacoinPayTo("Vd1fB3phUKKmV2FXoXLpCyTnugJ3KkpVFE", 0.05)
//	if err != nil {
//		log.Printf("error: %s\n", err)
//	}
//
//	fmt.Printf("broadcast: %s\n", result)
//}
