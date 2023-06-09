package config

type Config struct {
	ServerEndpoint      string `mapstructure:"USER_SERVER_ENDPOINT"`
	AuthServiceEndpoint string `mapstructure:"AUTH_SERVICE_ENDPOINT"`
}
