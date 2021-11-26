package config

import (
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*model.Configuration, error) {
	v := viper.New()
	v.SetConfigFile(path + "/.env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read env file: %w", err)
	}
	var configuration model.Configuration
	if err := v.Unmarshal(&configuration); err != nil {
		return nil, fmt.Errorf("Unable to decode into map, %w", err)
	}

	return &configuration, nil
}
