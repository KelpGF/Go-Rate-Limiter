package configs

import (
	"github.com/spf13/viper"
)

type configType struct {
	RateLimitDefaultIntervalInSeconds          int `mapstructure:"RATE_LIMIT_DEFAULT_INTERVAL_IN_SECONDS"`
	RateLimitDefaultLimitOfRequestsPerInterval int `mapstructure:"RATE_LIMIT_DEFAULT_LIMIT_OF_REQUEST_PER_INTERVAL"`
	RateLimitDefaultBanTimeInSeconds           int `mapstructure:"RATE_LIMIT_DEFAULT_REQUEST_BAN_TIME_IN_SECONDS"`

	RateLimitDefaultTokenIntervalInSeconds          int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_INTERVAL_IN_SECONDS"`
	RateLimitDefaultTokenLimitOfRequestsPerInterval int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_LIMIT_OF_REQUEST_PER_INTERVAL"`
	RateLimitDefaultTokenBanTimeInSeconds           int `mapstructure:"RATE_LIMIT_DEFAULT_TOKEN_BAN_TIME_IN_SECONDS"`

	RateLimitToken string `mapstructure:"RATE_LIMIT_TOKEN"`

	WebServerHost string `mapstructure:"WEB_SERVER_HOST"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	DB_DRIVER string `mapstructure:"DB_DRIVER"`
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

	if configs.DB_DRIVER == "" {
		configs.DB_DRIVER = "in_memory"
	}

	return configs
}
