package config

import (
	"github.com/spf13/viper"
)

var Cfg Config

type Config struct {
	Differ Differ `mapstructure:"differ"`
}

type Differ struct {
	Database01 Database `mapstructure:"database01"`
	Database02 Database `mapstructure:"database02"`
}

type Database struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Protocol string `mapstructure:"protocol"`
}

func LoadFile(filename string) (err error){
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.Unmarshal(&Cfg)
	return nil
}