// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package user

type User struct {
	Id             uint32 `json:"id"`
	Twitter_id     int64  `json:"twitter_id"`
	Telegram_id    int    `json:"telegram_id"`
	Username       string `json:"username"`
	Social         string `json:"social"`
	RegisteredDate int64  `json:"registered_date"`
}
