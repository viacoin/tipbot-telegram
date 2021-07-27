// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package database

import (
	"fmt"
	"github.com/romanornr/CryptoTwitterTipBot/bcoins"
	"log"
)

type TotalTransferredStat struct {
	Name      string
	Symbol    string
	AmountSat int64
	Amount    float64
}

func addTotalTransferred(transaction bcoins.Transaction) error {
	GetSession()
	stmt, err := db.Prepare("INSERT INTO total_transferred(coin_name, coin_symbol, amount) values(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error prepare statement adding total_transferred to database: %s\n", err)
	}

	_, err = stmt.Exec(transaction.CoinName, transaction.Coinsymbol, transaction.Amount)
	if err != nil {
		return fmt.Errorf("error adding transaction %s to database: %s\n", transaction.Coinsymbol, err)
	}
	return nil
}

func FindTotalTransferred(coinsymbol string) (TotalTransferredStat, error) {
	GetSession()
	rows, err := db.Query("SELECT * FROM total_transferred WHERE coin_symbol = $1", coinsymbol)
	if err != nil {
		return TotalTransferredStat{}, fmt.Errorf("Error db queury SELECT FROM user: %s\n", err)
	}

	stat := TotalTransferredStat{}

	for rows.Next() {
		err = rows.Scan(&stat.Name, &stat.Symbol, &stat.AmountSat)
	}

	rows.Close()
	if stat.Symbol == "" {
		return TotalTransferredStat{}, fmt.Errorf("No stats found for %s\n", coinsymbol)
	}

	stat.Amount = float64(stat.AmountSat) / 100000000 //sat to btc/via/dash/ltc

	return stat, nil
}

func AddOrUpdateTotalTransferred(transaction bcoins.Transaction) {
	_, err := FindTotalTransferred(transaction.Coinsymbol)
	if err != nil {
		fmt.Println("coin not in db.. adding")
		err := addTotalTransferred(transaction) //if coin isn't in db, add it
		if err != nil {
			log.Printf("error adding: %s\n", err)
		}
		return
	}

	_, err = FindTotalTransferred(transaction.Coinsymbol)
	if err != nil {
		log.Printf("error, totalTransferred info added to the database but can't be found after")
	}

	err = updateTotalTransferred(transaction)
	if err != nil {
		log.Printf("error updating transction for %s with amount %d: %s\n", transaction.Coinsymbol, transaction.Amount, err)
	}
}

func updateTotalTransferred(transaction bcoins.Transaction) error {
	GetSession()
	oldStat, err := FindTotalTransferred(transaction.Coinsymbol)
	newAmount := oldStat.AmountSat + transaction.Amount

	stmt, err := db.Prepare("update total_transferred set amount=? where coin_symbol=?")
	if err != nil {
		return fmt.Errorf("error prepare statement adding total_transferred to database: %s\n", err)
	}

	_, err = stmt.Exec(newAmount, transaction.Coinsymbol)
	return err
}

func GetAllCoinStats() []TotalTransferredStat {
	coins := bcoins.GetAllCoins()
	var result []TotalTransferredStat
	for _, coin := range coins {
		dbStats, _ := FindTotalTransferred(coin.Symbol)
		result = append(result, TotalTransferredStat{coin.Name, coin.Symbol, dbStats.AmountSat, dbStats.Amount})
	}
	return result
}
