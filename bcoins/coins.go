// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package bcoins

import (
	"fmt"
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/wire"
)

type Coin struct {
	Symbol     string
	Name       string
	Network    *Network
	Insight    *Insight
	FeePerByte int64
	Binance    bool
}

type Insight struct {
	Explorer string
	Api      string
}

type Network struct {
	Name     string
	Symbol   string
	xpubkey  byte
	xprivkey byte
	magic    wire.BitcoinNet
}

var coins = map[string]Coin{
	"via": {Name: "viacoin", Symbol: "via", Network: &Network{"viacoin", "via", 0x47, 0xC7, 0xcbc6680f},
		Insight: &Insight{"https://explorer.viacoin.org", "https://explorer.viacoin.org/api"}, FeePerByte: 110, Binance: true,
	},
	"btc": {Name: "bitcoin", Symbol: "btc", Network: &Network{"bitcoin", "btc", 0x00, 0x80, 0xf9beb4d9},
		Insight: &Insight{"https://insight.bitpay.com", "https://insight.bitpay.com/api"}, FeePerByte: 13, Binance: true,
	},
	"dash": {Name: "dash", Symbol: "dash", Network: &Network{"dash", "dash", 0x4c, 0xcc, 0xd9b4bef9},
		Insight: &Insight{"https://insight.dash.org/insight", "https://insight.dash.org/insight-api"}, FeePerByte: 11, Binance: true,
	},
	"ltc": {Name: "litecoin", Symbol: "ltc", Network: &Network{"litecoin", "ltc", 0x30, 0xb0, 0xfbc0b6db},
		Insight: &Insight{"https://insight.litecore.io", "https://insight.litecore.io/api"}, FeePerByte: 280, Binance: true,
	},
}

// returns all coins in a Coin struct slice
func GetAllCoins() []Coin {
	var coinList []Coin

	for _, coin := range coins {
		coinList = append(coinList, coin)
	}
	return coinList
}

// select a coin by symbol and return Coin struct and error
func SelectCoin(symbol string) (Coin, error) {
	if coins, ok := coins[symbol]; ok {
		return coins, nil
	}
	return Coin{}, fmt.Errorf("altcoin %s not found\n", symbol)
}

func (network Network) GetNetworkParams() *chaincfg.Params {
	networkParams := &chaincfg.MainNetParams
	networkParams.Name = network.Name
	networkParams.Net = network.magic
	networkParams.PubKeyHashAddrID = network.xpubkey
	networkParams.PrivateKeyID = network.xprivkey
	return networkParams
}
