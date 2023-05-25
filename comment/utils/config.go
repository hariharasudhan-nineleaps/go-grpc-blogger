package utils

import (
	"fmt"

	"github.com/hariharasudhan-nineleaps/go-grpc-blogger/comment/config"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*config.Config, error) {
	viper.AddConfigPath(path)   // path to config file
	viper.SetConfigName(".env") // name of config file
	viper.SetConfigType("env")  // type of config
	viper.AutomaticEnv()        // override the values from file from system env

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Invalid config.")
	}

	config := &config.Config{}
	err = viper.Unmarshal(config)

	if err != nil {
		return nil, fmt.Errorf("Invalid config Unmarshal failed.")
	}

	return config, nil
}
