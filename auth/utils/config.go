package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/auth/config"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*config.Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Invalid config. %v", err)
	}

	config := &config.Config{}
	cerr := viper.Unmarshal(config)

	if cerr != nil {
		return nil, fmt.Errorf("Invalid config Unmarshal failed.")
	}

	validate := validator.New()
	if verr := validate.Struct(config); verr != nil {
		return nil, fmt.Errorf("Missing required fields %v.", verr)
	}

	return config, nil
}
