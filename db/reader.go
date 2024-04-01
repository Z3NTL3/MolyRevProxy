/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package db

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"
	"github.com/spf13/viper"
)

func ReadConfig() {
	var err error
	defer func(err_ *error) {
		fmt.Println(errors.New(*err_).ErrorStack())
		os.Exit(-1)
	}(&err)

	cwd, err := os.Getwd()

	viper.AddConfigPath(cwd)
	viper.SetConfigName("config")
	viper.SetConfigFile("yaml")

	err = viper.ReadInConfig()
}
