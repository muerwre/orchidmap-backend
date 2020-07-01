package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host           string
	Port           string
	HasTls         bool
	VkClientId     string
	VkClientSecret string
	VkCallbackUrl  string
}

func InitConfig() (*Config, error) {
	config := &Config{
		Host:           viper.GetString("Host"),
		Port:           viper.GetString("Port"),
		VkClientId:     viper.GetString("Vk.ClientId"),
		VkClientSecret: viper.GetString("Vk.ClientSecret"),
		VkCallbackUrl:  viper.GetString("Vk.CallbackUrl"),
		HasTls:         len(viper.GetStringSlice("API.TlsFiles")) == 2,
	}

	return config, nil
}
