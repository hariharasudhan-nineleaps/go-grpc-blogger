package config

type Config struct {
	ServerEndpoint          string `mapstructure:"SERVER_ENDPOINT"`
	UserServiceEndpoint     string `mapstructure:"USER_SERVER_ENDPOINT"`
	KafkaEndpoint           string `mapstructure:"KAFKA_BROKER_ENDPOINT"`
	ActivityServiceEndpoint string `mapstructure:"ACTIVITY_SERVER_ENDPOINT"`
}
