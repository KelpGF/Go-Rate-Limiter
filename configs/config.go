package configs

import (
	"github.com/spf13/viper"
)

type configType struct {
	RateLimitDefaultIntervalInSeconds      int `mapstructure:"RATE_LIMIT_DEFAULT_INTERVAL_IN_SECONDS"`
	RateLimitDefaultRequestCountPerSeconds int `mapstructure:"RATE_LIMIT_DEFAULT_REQUEST_COUNT_PER_SECONDS"`
	RateLimitDefaultBanTimePerSeconds      int `mapstructure:"RATE_LIMIT_DEFAULT_REQUEST_BAN_TIME_PER_SECONDS"`

	RateLimitDefaultTokenIntervalInSeconds int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_INTERVAL_IN_SECONDS"`
	RateLimitDefaultTokenCountPerInterval  int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_REQUEST_COUNT_PER_SECONDS"`
	RateLimitDefaultTokenBanTimePerSeconds int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_BAN_TIME_PER_SECONDS"`

	RateLimitToken string `mapstructure:"RATE_LIMIT_TOKEN"`

	WebServerHost string `mapstructure:"WEB_SERVER_HOST"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
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

	if configs.RedisPort == "" {
		configs.RedisPort = "6379"
	}

	return configs
}
