package model

type Configuration struct {
	Consumer ConsumerConfig `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
}

type ConsumerConfig struct {
	BaseUrl      string `mapstructure:"base_url"`
	CronSchedule string `mapstructure:"cron"`
}

type DatabaseConfig struct {
	Username string `mapstructure:"database_username"`
	Password string `mapstructure:"database_password"`
	Host     string `mapstructure:"database_host"`
	Port     string `mapstructure:"database_port"`
	Name     string `mapstructure:"database_name"`
}
