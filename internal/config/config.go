package config

import (
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*model.Configuration, error) {
	v := viper.New()
	v.SetConfigFile(path + "/.env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Failed to read env file")
	}
	var configuration model.Configuration
	err := v.Unmarshal(&configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode into map")
	}

	return &configuration, nil
}
