package config

import (
	_ "embed"
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

//go:embed default-config.toml
var defaultConfigBytes []byte

func InitConfig() error {
	// if provided env variable DB_DSN, use env variables as config source
	if os.Getenv("API_KEY") != "" {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		return nil
	}

	// else use config file
	_, err := os.ReadFile("data/config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll("data", os.ModePerm)
			_ = os.WriteFile("data/config.toml", defaultConfigBytes, os.ModePerm)
			return errors.New("config.toml not found, create default config.toml")
		} else {
			return err
		}
	}

	viper.SetConfigFile("data/config.toml")
	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
