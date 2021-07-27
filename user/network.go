// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package user

import (
	"fmt"
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
)

func SelectNetwork(symbol string) (bcoins.Network, error) {
	coin, err := bcoins.SelectCoin(symbol)
	if err != nil {
		return bcoins.Network{}, fmt.Errorf("Network for %s not found\n", symbol)
	}
	return *coin.Network, nil
}
