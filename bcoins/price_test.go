// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package bcoins

import (
	"fmt"
	"testing"
)

func TestGetTicker(t *testing.T) {
	ticker, err := GetTicker("viacoin")
	if err != nil {
		fmt.Println(err)
	}
	if (*ticker)[0].Name != "Viacoin" {
		t.Errorf("error coinmarketcap api. Expected: %s Got: %s", "Viacoin", (*ticker)[0].Name)
	}
}

func TestGetTotalUSD(t *testing.T) {
	result := GetTotalUSD(10, "viacoin")
	if result == "" {
		t.Errorf("error getting usd price")
	}
}

func TestStats(t *testing.T) {
}
