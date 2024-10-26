package configs

import (
	"github.com/spf13/viper"
)

type configType struct {
	RateLimitDefaultIntervalInSeconds      int `mapstructure:"RATE_LIMIT_DEFAULT_INTERVAL_IN_SECONDS"`
	RateLimitDefaultRequestCountPerSeconds int `mapstructure:"RATE_LIMIT_DEFAULT_REQUEST_COUNT_PER_SECONDS"`

	RateLimitDefaultTokenIntervalInSeconds int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_INTERVAL_IN_SECONDS"`
	RateLimitDefaultTokenCountPerInterval  int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_COUNT_PER_INTERVAL"`

	WebServerHost string `mapstructure:"WEB_SERVER_HOST"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) *configType {
	var configs *configType

	viper.SetConfigName("rate_limiter_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&configs)
	if err != nil {
		panic(err)
	}

	if configs.WebServerPort == "" {
		configs.WebServerPort = "3000"
	}

	return configs
}
