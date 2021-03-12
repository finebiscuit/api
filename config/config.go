package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port uint

	DBType   string
	DBSource string
}

func New(v *viper.Viper) *Config {
	return &Config{
		Port:     v.GetUint("port"),
		DBType:   v.GetString("database-type"),
		DBSource: v.GetString("database-source"),
	}
}
