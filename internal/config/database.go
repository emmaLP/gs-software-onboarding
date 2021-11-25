package config

type DatabaseConfig struct {
	Username string `mapstructure:"database_username"`
	Password string `mapstructure:"database_password"`
	Host     string `mapstructure:"database_host"`
	Port     string `mapstructure:"database_port"`
	Name     string `mapstructure:"database_name"`
}
