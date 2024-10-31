package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KelpGF/Go-Rate-Limiter/configs"
	"github.com/KelpGF/Go-Rate-Limiter/internal/middlewares"
	"github.com/KelpGF/Go-Rate-Limiter/package/rate_limiter"
)

func main() {
	configs := configs.LoadConfig(".")

	mux := http.NewServeMux()

	rateLimiterItemRepositoryMemory := rate_limiter.NewRateLimiterItemRepositoryMemory()
	rateLimiterService := rate_limiter.NewRateLimiterServiceImpl(rateLimiterItemRepositoryMemory)

	rateLimiterService.SetConfig("ip", rate_limiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultRequestCountPerSeconds,
		IntervalInSeconds: configs.RateLimitDefaultIntervalInSeconds,
	})

	rateLimiterService.SetConfig("token", rate_limiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultTokenCountPerInterval,
		IntervalInSeconds: configs.RateLimitDefaultTokenIntervalInSeconds,
	})

	rateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(rateLimiterService, configs.RateLimitToken)

	mux.Handle("/", rateLimiterMiddleware.RateLimiter(HelloHandler()))

	http.ListenAndServe(configs.WebServerHost+":"+configs.WebServerPort, mux)
}

func HelloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		w.Write([]byte("Hello World"))
	})
}
