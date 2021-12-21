package config

import (
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*model.Configuration, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("Failed to read config: %w", err)
		}
	}
	if err := bindEnvs(v, model.Configuration{}); err != nil {
		return nil, fmt.Errorf("Failed to bind environment variables: %w", err)
	}
	var configuration model.Configuration
	setDefaults(v)
	if err := v.Unmarshal(&configuration); err != nil {
		return nil, fmt.Errorf("Unable to decode into map, %w", err)
	}

	return &configuration, nil
}

// bindEnv is a workaround for a known issue in viper
// The issue means that env variables cannot be read unless there is a blank config file or every config value set to have a default
// Issue ref: https://github.com/spf13/viper/issues/761
func bindEnvs(v *viper.Viper, config model.Configuration) error {
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(config, &envKeysMap); err != nil {
		return fmt.Errorf("Failed to determine keys %w", err)
	}
	for k := range *envKeysMap {
		if bindErr := v.BindEnv(k); bindErr != nil {
			return fmt.Errorf("Failed to bind env variables. %w", bindErr)
		}
	}
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("cron", "*/15 * * * *")
	v.SetDefault("workers", 5)

	v.SetDefault("api_address", ":8080")
}
