package model

type Configuration struct {
	Publisher  PublisherConfig  `mapstructure:",squash"`
	Consumer   ConsumerConfig   `mapstructure:",squash"`
	Database   DatabaseConfig   `mapstructure:",squash"`
	RabbitMq   RabbitMqConfig   `mapstructure:",squash"`
	Api        APIConfig        `mapstructure:",squash"`
	Grpc       GrpcServerConfig `mapstructure:",squash"`
	GrpcClient GrpcClientConfig `mapstructure:",squash"`
	Cache      CacheConfig      `mapstructure:",squash"`
}

type PublisherConfig struct {
	BaseUrl      string `mapstructure:"base_url"`
	CronSchedule string `mapstructure:"cron"`
}

type ConsumerConfig struct {
	NumberOfWorkers int `mapstructure:"workers"`
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

type GrpcClientConfig struct {
	GrpcAddress string `mapstructure:"grpc_address"`
}

type GrpcServerConfig struct {
	Port int `mapstructure:"grpc_port"`
}

type CacheConfig struct {
	Address string `mapstructure:"cache_address"`
}

type RabbitMqConfig struct {
	Username  string `mapstructure:"rabbitmq_username"`
	Password  string `mapstructure:"rabbitmq_password"`
	Host      string `mapstructure:"rabbitmq_host"`
	Port      string `mapstructure:"rabbitmq_port"`
	QueueName string `mapstructure:"rabbitmq_queue_name"`
}
