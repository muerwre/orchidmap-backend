package api

import "github.com/spf13/viper"

type Config struct {
	// The port to bind the web application server to
	Port int

	// The number of proxies positioned in front of the API. This is used to interpret
	// X-Forwarded-For headers.
	ProxyCount int

	// Output debugging messages
	Debug bool

	// TlsHosts
	TlsHosts []string
}

func InitConfig() (*Config, error) {
	config := &Config{
		Port:       viper.GetInt("Port"),
		ProxyCount: viper.GetInt("ProxyCount"),
		Debug:      viper.GetBool("API.Debug"),
		TlsHosts:   viper.GetStringSlice("API.TlsHosts"),
	}

	if config.Port == 0 {
		config.Port = 9092
	}

	return config, nil
}
