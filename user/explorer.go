// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package user

import (
	"fmt"
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
)

func GetInsightExplorer(symbol string) (bcoins.Insight, error) {
	coin, err := bcoins.SelectCoin(symbol)
	if err != nil {
		return bcoins.Insight{}, fmt.Errorf("this altcoin does not have an insight explorer")
	}
	return *coin.Insight, nil
}
