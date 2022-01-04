package model

type Configuration struct {
	Consumer ConsumerConfig `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Api      APIConfig      `mapstructure:",squash"`
	Grpc     GrpcConfig     `mapstructure:",squash"`
	Cache    CacheConfig    `mapstructure:",squash"`
}

type ConsumerConfig struct {
	BaseUrl         string `mapstructure:"base_url"`
	CronSchedule    string `mapstructure:"cron"`
	NumberOfWorkers int    `mapstructure:"workers"`
}

type DatabaseConfig struct {
	Username string `mapstructure:"database_username"`
	Password string `mapstructure:"database_password"`
	Host     string `mapstructure:"database_host"`
	Port     string `mapstructure:"database_port"`
	Name     string `mapstructure:"database_name"`
}

type APIConfig struct {
	Address string `mapstructure:"api_address"`
}

type GrpcConfig struct {
	Port int `mapstructure:"grpc_port"`
}

type CacheConfig struct {
	Address string `mapstructure:"cache_address"`
}
