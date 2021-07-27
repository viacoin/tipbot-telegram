// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package main

import (
	"github.com/romanornr/CryptoTwitterTipBot/bots"
	"github.com/spf13/viper"
)

func main() {
	telegram := bots.NewTelegramBot(viper.GetString("telegram.token"))
	//telegram2 := bots.NewTelegramBot(viper.GetString("telegram.token2"))

	//go telegram2.Update()
	telegram.Update()
	//telegram2.Update()
}
