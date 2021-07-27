// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package database

import (
	"fmt"
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
	"testing"
)

func TestFindTotalTransferred(t *testing.T) {
	tx := bcoins.Transaction{}
	tx.CoinName = "viacoin"
	tx.Coinsymbol = "via"
	tx.Amount = 100
	amount, err := FindTotalTransferred(tx.Coinsymbol)
	if err == nil {
		//t.Errorf("Expected nothing found\n")
	}
	fmt.Println(amount, err)
}

func TestAddOrUpdateTotalTransferred(t *testing.T) {
	tx := bcoins.Transaction{}
	tx.Coinsymbol = "via"
	tx.CoinName = "viacoin"
	tx.Amount = 100000000
	AddOrUpdateTotalTransferred(tx)
	AddOrUpdateTotalTransferred(tx)
	AddOrUpdateTotalTransferred(tx)
	amount, err := FindTotalTransferred(tx.Coinsymbol)
	if err != nil {
		t.Errorf("error finding transaction after adding: %s\n", err)
	}
	fmt.Println(amount)
}

func TestGetAllCoinStats(t *testing.T) {
	stats := GetAllCoinStats()
	fmt.Println(stats)
}
