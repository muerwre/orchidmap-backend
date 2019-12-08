package logger

import "github.com/spf13/viper"

type Config struct {
	// The number of proxies positioned in front of the API. This is used to interpret
	// X-Forwarded-For headers.
	ProxyCount int
}

func InitConfig() (*Config, error) {
	config := &Config{
		ProxyCount: viper.GetInt("ProxyCount"),
	}

	return config, nil
}
