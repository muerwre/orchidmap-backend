package app

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Host   string
	Port   string
	HasTls bool

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
		HasTls:         len(viper.GetStringSlice("API.TlsFiles")) == 2,
	}

	fmt.Printf("TLS FILES LENGTH IS %v", len(viper.GetStringSlice("TlsFiles")))
	fmt.Printf("TLS ENABLED? %v", config.HasTls)

	return config, nil
}
