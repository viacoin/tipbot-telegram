// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package insightjson

type AddressInfo struct {
	Address                  string   `json:"addrStr,omitempty"`
	Balance                  float64  `json:"balance"`
	BalanceSat               int64    `json:"balanceSat"`
	TotalReceived            float64  `json:"totalReceived"`
	TotalReceivedSat         int64    `json:"totalReceivedSat"`
	TotalSent                float64  `json:"totalSent"`
	TotalSentSat             int64    `json:"totalSentSat"`
	UnconfirmedBalance       float64  `json:"unconfirmedBalance"`
	UnconfirmedBalanceSat    int64    `json:"unconfirmedBalanceSat"`
	UnconfirmedTxAppearances int64    `json:"unconfirmedTxAppearances"`
	TxAppearances            int64    `json:"txAppearances "`
	TransactionsID           []string `json:"transactions,omitempty"`
}

type UnspentOutputs []struct {
	Address       string  `json:"address"`
	Txid          string  `json:"txid"`
	Vout          int     `json:"vout"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Satoshis      int     `json:"satoshis"`
	Height        int     `json:"height"`
	Confirmations int     `json:"confirmations"`
}
