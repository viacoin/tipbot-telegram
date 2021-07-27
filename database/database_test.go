// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package database

import (
	"fmt"
	"github.com/romanornr/CryptoTwitterTipBot/user"
	"testing"
)

func TestSetup(t *testing.T) {
	Setup()
}

//func TestAddTwitterUser(t *testing.T) {
//	user := user.User{}
//	user.Twitter_id = 233223
//	user.Username = "@RNR_0"
//	user.Social = "twitter"
//
//	AddUser(user)
//}
//
//func TestGetTwitterUser(t *testing.T) {
//	user2 := user.User{}
//	user2.Twitter_id = 233223
//	twitterUser := FindTwitterUser(user2.Twitter_id)
//	if twitterUser.Username != "@RNR_0" {
//		t.Errorf("error getting user from database. Got: %s Exptected: %s\n", twitterUser.Username, "@RNR_0")
//	}
//
//	if twitterUser.Social != "twitter" {
//		t.Errorf("expected twitter as social")
//	}
//}

func TestAddOrUpdateTelegramUser(t *testing.T) {

	user := user.User{}
	user.Username = "@Mike"
	user.Social = "telegram"
	user.Telegram_id = 22
	AddOrUpdateTelegramUser(user)

	u, _ := FindTelegramUserByUsername("@Mike")
	if u.Username != "@Mike" {
		t.Errorf("error creating user: %s\n", u.Username)
	}

	user.Username = "@RomanoRnr"
	AddOrUpdateTelegramUser(user)

	result, _ := FindTelegramUserByUsername("@RomanoRnr")
	if result.Telegram_id != 22 {
		t.Errorf("error updating telegram username for existing telegram user with id.")
	}

	if result.Username != "@RomanoRnr" {
		t.Errorf("error updating telegram username for existing telegram user with id")
	}

	_, err := FindTelegramUserByUsername("@Mike")
	if err == nil {
		t.Errorf("User should not be in the database")
	}
}

func TestFindTelegramUserById(t *testing.T) {
	user := user.User{}
	user.Username = "@RomanoRnr"
	user.Social = "telegram"
	user.Telegram_id = 19

	AddUser(user)

	user, err := FindTelegramUserById(user.Telegram_id)
	if err == nil {
		t.Errorf("user not found by telegram id")
	}
}

// test to update a username
// this can be used when a user changed his telegram username but has the same telegram_id
func TestUpdateTelegramUsername(t *testing.T) {

	user := user.User{}
	user.Username = "@MikeVer"
	user.Social = "telegram"
	user.Telegram_id = 22

	AddUser(user)

	user.Username = "@RomanoRnr"
	err := UpdateTelegramUsername(user)
	if err != nil {
		t.Errorf("error man: %s\n", err)
	}

	user, _  = FindTelegramUserByUsername("@RomanoRnr")
	if user.Telegram_id != 22 {
		t.Errorf("user %s not found\n", user.Username)
	}

	_, err = FindTelegramUserByUsername("@MikeVer")
	if err == nil {
		t.Errorf("user not updated")
	}
}

func TestAddUser(t *testing.T) {
	user1 := user.User{}
	user1.Username = "@RomanoRnr"
	user1.Social = "telegram"
	user1.Telegram_id = 1

	user2 := AddOrUpdateTelegramUser(user1)
	if user2.Id != 0 {
		fmt.Printf("database user id: %d\n", user2.Id)
		t.Errorf("fail")
	}

	user3 := AddOrUpdateTelegramUser(user1)
	if user3.Id != 0 {
		fmt.Printf("database user id: %d\n", user3.Id)
		t.Errorf("fail")
	}

	user4 := user.User{}
	user4.Username = "@MikeVer2"
	user4.Social = "telegram"
	user4.Telegram_id = 28

	AddOrUpdateTelegramUser(user4)
	r, _ := FindTelegramUserByUsername("@MikeVer2")
	if r.Id == 0 {
		t.Errorf("database user id: %d\n", r.Id)
	}

}