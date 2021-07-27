// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package user

import (
	"encoding/json"
	"fmt"
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
	"github.com/romanornr/CryptoTwitterTipBot/insightjson"
	"net/http"
	"time"
)

func (user User) AddressInfo(network bcoins.Network) (*insightjson.AddressInfo, error) {

	insightExplorer, _ := GetInsightExplorer(network.Symbol)
	publickey := user.PublicKey(network)

	url := fmt.Sprintf("%s/addr/%s", insightExplorer.Api, publickey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &insightjson.AddressInfo{}, fmt.Errorf("API Request error: %s\n", err)
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return &insightjson.AddressInfo{}, fmt.Errorf("error api client.DO: %s\n", err)
	}

	defer resp.Body.Close()

	var addressInfo insightjson.AddressInfo
	if err := json.NewDecoder(resp.Body).Decode(&addressInfo); err != nil {
		return &insightjson.AddressInfo{}, fmt.Errorf("error decoding insightjson addressInfo: %s", err)
	}

	return &addressInfo, nil
}
