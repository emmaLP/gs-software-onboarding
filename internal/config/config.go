package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Configuration struct {
	Consumer ConsumerConfig `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
}

func LoadConfig(path string) (*Configuration, error) {
	v := viper.New()
	v.SetConfigFile(path + "/.env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Failed to read env file")
	}
	var configuration Configuration
	err := v.Unmarshal(&configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode into map")
	}

	return &configuration, nil
}
