package config

type Config struct {
	ServerEndpoint      string `mapstructure:"SERVER_ENDPOINT"`
	UserServiceEndpoint string `mapstructure:"USER_SERVER_ENDPOINT"`
}
