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

	rateLimiterService := rate_limiter.NewRateLimiterServiceImpl(nil)

	var rateLimiterItemRepositor rate_limiter.RateLimiterItemRepository

	if configs.DB_DRIVER == "in_memory" {
		rateLimiterItemRepositor = rate_limiter.NewRateLimiterItemRepositoryMemory()
	} else if configs.DB_DRIVER == "redis" {
		redis := infra.NewRedis(&redis.Options{
			Addr: configs.RedisHost + ":" + configs.RedisPort,
		})
		rateLimiterItemRepositor = infra.NewRateLimiterItemRedisRepository(redis)
	}

	rateLimiterService.SetRateLimiterItemRepository(rateLimiterItemRepositor)

	rateLimiterService.AddConfig("ip", rate_limiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultLimitOfRequestsPerInterval,
		IntervalInSeconds: configs.RateLimitDefaultIntervalInSeconds,
		BanTimeInSeconds:  configs.RateLimitDefaultBanTimeInSeconds,
	})

	rateLimiterService.AddConfig("token", rate_limiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultTokenLimitOfRequestsPerInterval,
		IntervalInSeconds: configs.RateLimitDefaultTokenIntervalInSeconds,
		BanTimeInSeconds:  configs.RateLimitDefaultTokenBanTimeInSeconds,
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
