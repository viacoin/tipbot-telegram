// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
package config

import (
	"github.com/spf13/viper"
	"log"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
)

func GetViperConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("app")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("No configuration file loaded !\n%s", err)
	}
	return err
}
