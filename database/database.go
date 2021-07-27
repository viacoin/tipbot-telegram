// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/romanornr/CryptoTwitterTipBot/user"
	"log"
	"time"
)

var db *sql.DB

func GetSession() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", "./database.db")
		checkErr(err)
	}
	return db
}

func Setup() {
	GetSession()
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS user(id INTEGER PRIMARY KEY AUTOINCREMENT, twitter_id INTEGER, telegram_id INTEGER, username VARCHAR(64), social VARCHAR(20), registered_date INTEGER)")
	checkErr(err)
	_, err = statement.Exec()
	checkErr(err)
	setupTotalTransferredTable()
}

func setupTotalTransferredTable() {
	GetSession()
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS total_transferred(coin_name VARCHAR(20), coin_symbol VARCHAR(10), amount INTEGER)")
	checkErr(err)
	_, err = statement.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error SQL: %s", err)
	}
}

func AddUser(user user.User) error {
	GetSession()
	stmt, err := db.Prepare("INSERT INTO user(twitter_id, telegram_id, username, social, registered_date) values(?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error prepare statement adding user to database: %s\n", err)
	}

	user.RegisteredDate = time.Now().Unix()
	_, err = stmt.Exec(user.Twitter_id, user.Telegram_id, user.Username, user.Social, user.RegisteredDate)
	if err != nil {
		return fmt.Errorf("error adding user %s to database: %s\n", user.Username, err)
	}
	fmt.Printf("Added %s to the database\n", user.Username)
	return nil
}

func FindTwitterUser(twitter_id int64) *user.User {
	GetSession()
	rows, err := db.Query("SELECT * FROM user WHERE twitter_id = $1", twitter_id)
	if err != nil {
		log.Printf("Error db queury SELECT FROM user: %s\n", err)
	}

	var id uint32
	//var twitter_id uint64
	var telegram_id int
	var username string
	var social string
	var registered_date int64

	for rows.Next() {
		err = rows.Scan(&id, &twitter_id, &telegram_id, &username, &social, &registered_date)
	}

	user := user.User{}
	user.Id = id
	user.Twitter_id = twitter_id
	user.Username = username
	user.Social = social
	user.RegisteredDate = registered_date

	rows.Close()
	return &user
}

func FindTelegramUserById(telegram_id int) (user.User, error) {
	GetSession()
	rows, err := db.Query("SELECT * FROM user WHERE telegram_id = $? AND social = ?", telegram_id, "telegram")
	if err != nil {
		return user.User{}, fmt.Errorf("Error db queury SELECT FROM user: %s\n", err)
	}

	user := user.User{}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Twitter_id, &user.Telegram_id, &user.Username, &user.Social, &user.RegisteredDate)
	}

	rows.Close()
	if user.Social != "telegram" {
		return user, errors.New("user not found")
	}
	return user, nil
}

func FindTelegramUserByUsername(username string) (user.User, error) {
	GetSession()
	rows, err := db.Query("SELECT * FROM user WHERE username = $1 AND social = $2", username, "telegram")
	if err != nil {
		return user.User{}, fmt.Errorf("Error db queury SELECT FROM user: %s\n", err)
	}

	user := user.User{}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Twitter_id, &user.Telegram_id, &user.Username, &user.Social, &user.RegisteredDate)
	}

	rows.Close()
	if user.Id == 0 {
		return user, errors.New("user not found")
	}
	return user, nil
}

//func FindOrCreateTelegramUserByUsername(user user.User) user.User {
//	_, err := FindTelegramUserByUsername(user.Username)
//	if err != nil {
//		AddUser(user)
//		//AddOrUpdateTelegramUser(user)
//	}
//
//	user, err = FindTelegramUserByUsername(user.Username)
//	if err != nil {
//		log.Printf("error finding telegram username after inserting: %s\n", err)
//	}
//	return user
//}

//count total amount of users in the database
func GetTotalUsers() int {
	GetSession()
	rows, err := db.Query("SELECT COUNT(*) as count FROM user")
	if err != nil {
		log.Printf("Error db queury SELECT FROM user: %s\n", err)
	}
	return checkCount(rows)
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func AddOrUpdateTelegramUser(user user.User) user.User {

	// only try to find telegram id if it's known, if not it's 0 and unknown
	if user.Telegram_id != 0 {
		_, err := FindTelegramUserById(user.Telegram_id)
		if err != nil {
			//fmt.Printf("updating: %s\n with telegram id: %d\n", user.Username, user.Telegram_id)
			err = UpdateTelegramUsername(user)
		}
	}

	// since no telegram id was found, find the user in the database
	// if nothing found, add the user
	dbuser, err := FindTelegramUserByUsername(user.Username)
	if err != nil {
		fmt.Printf("adding telegram user %s into the database\n", user.Username)
		err := AddUser(user) //add if user is not in the database
		if err != nil {
			log.Printf("error adding: %s\n", err)
		}
	}

	//user, err = FindTelegramUserById(user.Telegram_id)
	//if err == nil {
	//	//fmt.Printf("updating: %s\n with telegram id: %d\n", user.Username, user.Telegram_id)
	//	//UpdateTelegramUsername(user)
	//	return user
	//}

	userx, err := FindTelegramUserByUsername(dbuser.Username)
	if err != nil {
		return userx
	}
	//
	//// potentially update the telegram_id
	//if err == nil {
	//	updateTelegramUserId(user)
	//}



	fmt.Printf("username is %s\n", user.Username)
	return user

}

func updateTelegramUserId(user user.User) error {
	GetSession()

	stmt, err := db.Prepare("update user set telegram_id=? where username=? AND social = $2")
	if err != nil {
		return fmt.Errorf("error prepare statement updating user to database: %s\n", err)
	}
	_, err = stmt.Exec(user.Telegram_id, user.Username, user.Social)
	return err
}

func UpdateTelegramUsername(user user.User) error {
	stmt, err := db.Prepare("update user set username=? where telegram_id=? AND social=?")
	if err != nil {
		return fmt.Errorf("error prepare statement updating user to database: %s\n", err)
	}
	res, err := stmt.Exec(user.Username, user.Telegram_id, user.Social)
	if err != nil {
		fmt.Errorf("error exec: %s\n", err)
	}
	fmt.Println(res)
	return err
}
