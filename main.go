package main

import (
	"net/http"

	"github.com/KelpGF/Go-Rate-Limiter/configs"
	"github.com/KelpGF/Go-Rate-Limiter/internal/infra"
	"github.com/KelpGF/Go-Rate-Limiter/internal/middlewares"
	"github.com/KelpGF/Go-Rate-Limiter/package/rate_limiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	configs := configs.LoadConfig(".")

	redis := infra.NewRedis(&redis.Options{
		Addr: configs.RedisHost + ":" + configs.RedisPort,
	})
	rateLimiterItemRedisRepository := infra.NewRateLimiterItemRedisRepository(redis)
	rateLimiterService := rate_limiter.NewRateLimiterServiceImpl(rateLimiterItemRedisRepository)

	rateLimiterService.SetConfig("ip", rate_limiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultRequestCountPerSeconds,
		IntervalInSeconds: configs.RateLimitDefaultIntervalInSeconds,
		BanTimeInSeconds:  configs.RateLimitDefaultBanTimePerSeconds,
	})

	rateLimiterService.SetConfig("token", rate_limiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultTokenCountPerInterval,
		IntervalInSeconds: configs.RateLimitDefaultTokenIntervalInSeconds,
		BanTimeInSeconds:  configs.RateLimitDefaultTokenBanTimePerSeconds,
	})

	rateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(rateLimiterService, configs.RateLimitToken)

	mux := http.NewServeMux()
	mux.Handle("/", rateLimiterMiddleware.RateLimiter(HelloHandler()))

	http.ListenAndServe(configs.WebServerHost+":"+configs.WebServerPort, mux)
}

func HelloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
}
