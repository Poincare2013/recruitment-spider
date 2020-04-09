package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type tomlConfig struct {
	Mongo       mongoConfig    `toml:"mongo"`
	Rabbit      rabbitConfig   `toml:"rabbit"`
	Concurrency concurrency    `toml:"concurrency"`
}

type concurrency struct {
	Num int
}

type rabbitConfig struct {
	Addr string
}

type mongoConfig struct {
	Addr string
	Db string `toml:"DB"`
}

var myTomlConfig tomlConfig

func GetTomlConfig() *tomlConfig {
	return &myTomlConfig
}

func init() {
	if _, err := toml.DecodeFile("./config/config.toml", &myTomlConfig); err != nil {
		log.Fatal(err)
	}
}
