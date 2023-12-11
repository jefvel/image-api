package config

import (
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	DbAddress string `mapstructure:"POSTGRES_ADDRESS"`
	DbPort    int    `mapstructure:"POSTGRES_PORT"`
	DbUser    string `mapstructure:"POSTGRES_USER"`
	DbPass    string `mapstructure:"POSTGRES_PASSWORD"`
	DbName    string `mapstructure:"POSTGRES_DB"`
}

var config Config

func DecoderErrorUnset(c *mapstructure.DecoderConfig) {
	c.ErrorUnset = true
}

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Could not read config file", err)
	}

	if err := viper.Unmarshal(&config, DecoderErrorUnset); err != nil {
		log.Fatalln("Could not find config values", err)
	}
}

func GetConfig() Config {
	return config
}
