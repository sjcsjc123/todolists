package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

var Conf *Config

func init() {
	Conf = &Config{
		viper.New(),
	}
	Conf.SetConfigName("application")
	Conf.SetConfigType("yaml")
	Conf.AddConfigPath(".")
	err := Conf.ReadInConfig()
	if err != nil {

	}
	Conf.WatchConfig()
}
