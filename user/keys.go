// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package user

import (
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
	"github.com/spf13/viper"
	"github.com/viacoin/viad/chaincfg"
	btcutil "github.com/viacoin/viautil"
	"github.com/viacoin/viautil/hdkeychain"
	"log"
)

// generate private key by using the app.yml file with the master key
func (user User) PrivateKey() *hdkeychain.ExtendedKey {
	master := viper.GetString("master.key")
	masterkey, err := hdkeychain.NewKeyFromString(master)
	if err != nil {
		log.Printf("error could not generate privkey for twitter user: %s\n", err)
	}

	key, _ := masterkey.Child(hdkeychain.HardenedKeyStart + user.Id)
	return key
}

// convert the key to a wif format
func (user User) PrivateKeyWif(network bcoins.Network) *btcutil.WIF {
	privateKey, err := user.PrivateKey().ECPrivKey()
	if err != nil {
		log.Printf("error could not convert private key to ECPrivKey: %s\n", err)
	}

	wif, err := btcutil.NewWIF(privateKey, network.GetNetworkParams(), true)
	if err != nil {
		log.Printf("error could not convert private key to decompressed WIF: %s\n", err)
	}
	return wif
}

// generate a public key by using the private key
// the public key is in a compressed format
func (user User) PublicKey(network bcoins.Network) *btcutil.AddressPubKeyHash {
	account := user.PrivateKeyWif(network)
	addresspubkey, _ := btcutil.NewAddressPubKey(account.PrivKey.PubKey().SerializeCompressed(), &chaincfg.MainNetParams)
	return addresspubkey.AddressPubKeyHash()
}
