package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host string
	Port string

	// Client id for vk api
	VkClientId string

	// Client secret for vk api
	VkClientSecret string
}

func InitConfig() (*Config, error) {
	config := &Config{
		Host:           viper.GetString("Host"),
		Port:           viper.GetString("Port"),
		VkClientId:     viper.GetString("Vk.ClientId"),
		VkClientSecret: viper.GetString("Vk.ClientSecret"),
	}

	return config, nil
}
