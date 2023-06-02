package config

type Config struct {
	ServerEndpoint string `mapstructure:"METRIC_SERVER_ENDPOINT"`
	KafkaEndpoint  string `mapstructure:"METRIC_KAFKA_BROKER_ENDPOINT"`
	RedisEndpoint  string `mapstructure:"METRIC_REDIS_ENDPOINT"`
}
