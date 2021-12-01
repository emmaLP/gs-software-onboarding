package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	envVars := env{
		"BASE_URL":          "localhost:8000",
		"DATABASE_USERNAME": "test_username",
		"DATABASE_PASSWORD": "test_password",
		"DATABASE_HOST":     "localhost",
		"DATABASE_PORT":     "30000",
		"DATABASE_NAME":     "hackernews",
	}
	tests := []struct {
		name        string
		expected    *model.Configuration
		filePath    string
		envVars     env
		expectedErr string
	}{
		{
			name:     "Successfully load app.env file",
			filePath: "../..",
			envVars:  envVars,
			expected: &model.Configuration{
				Consumer: model.ConsumerConfig{
					BaseUrl:         "localhost:8000",
					NumberOfWorkers: 5,
					CronSchedule:    "*/15 * * * *",
				},
				Database: model.DatabaseConfig{
					Username: "test_username",
					Password: "test_password",
					Host:     "localhost",
					Port:     "30000",
					Name:     "hackernews",
				},
			}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.envVars {
				err := os.Setenv(k, fmt.Sprintf("%v", v))
				assert.NoError(t, err)
			}

			cfg, err := LoadConfig(test.filePath)
			if test.expectedErr != "" {
				assert.EqualErrorf(t, err, test.expectedErr, "Error should be: %v, got: %v", test.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, cfg)
			}

			os.Clearenv()
		})
	}
}

type env map[string]interface{}
