package config

type Config struct {
	ServerEndpoint string `mapstructure:"AUTH_SERVER_ENDPOINT" validate:"required"`
	JWTSecret      string `mapstructure:"AUTH_JWT_SECRET" validate:"required"`
}
