package utils

import (
	"fmt"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/metric/config"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*config.Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Invalid config. %v", err)
	}

	config := &config.Config{}
	err = viper.Unmarshal(config)

	if err != nil {
		return nil, fmt.Errorf("Invalid config Unmarshal failed.")
	}

	return config, nil
}
