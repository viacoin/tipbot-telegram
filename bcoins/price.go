// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package bcoins

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Ticker []struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	PriceUSD string `json:"price_usd"`
	priceBTC string `json:"price_btc"`
}

func GetTicker(name string) (*Ticker, error) {

	url := fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/%s", name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Ticker{}, fmt.Errorf("API coinmarketcap Request error: %s\n", err)
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return &Ticker{}, fmt.Errorf("error api client.DO: %s\n", err)
	}

	defer resp.Body.Close()

	var ticker Ticker
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return &ticker, fmt.Errorf("error decoding insightjson addressInfo: %s", err)
	}

	return &ticker, nil
}

func GetTotalUSD(amount float64, name string) string {
	ticker, _ := GetTicker(name)
	f, err := strconv.ParseFloat((*ticker)[0].PriceUSD, 64)
	if err != nil {
		log.Printf("error GetTotalUSD convert string to float")
	}
	return fmt.Sprintf("%.2f", f*amount)
}

//func Stats() {
//	var c []Coin
//
//	for _, coin := range coins {
//		c = append(c, coin)
//	}
//
//	for _, coin := range c{
//		result, err := database.FindTotalTransferred(coin.Name)
//		if err == nil {
//			fmt.Println(result.Name)
//		}
//	}
//
//	//for _, x := range c {
//	//	fmt.Println(x.Name)
//	//}
//}
