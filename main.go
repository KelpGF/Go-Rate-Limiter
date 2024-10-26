package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/KelpGF/Go-Rate-Limiter/configs"
	ratelimiter "github.com/KelpGF/Go-Rate-Limiter/package/rate-limiter"
)

const RATE_LIMIT_ERROR_MESSAGE = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func main() {
	configs := configs.LoadConfig(".")

	mux := http.NewServeMux()

	rateLimiterItemRepositoryMemory := ratelimiter.NewRateLimiterItemRepositoryMemory()
	rateLimiterService := ratelimiter.NewRateLimiterServiceImpl(rateLimiterItemRepositoryMemory)

	rateLimiterService.SetConfig("ip", ratelimiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultRequestCountPerSeconds,
		IntervalInSeconds: configs.RateLimitDefaultIntervalInSeconds,
	})

	rateLimiterService.SetConfig("token", ratelimiter.RateLimiterConfig{
		LimitPerInterval:  configs.RateLimitDefaultTokenCountPerInterval,
		IntervalInSeconds: configs.RateLimitDefaultTokenIntervalInSeconds,
	})

	mux.Handle("/", RateLimiter(HelloHandler(), rateLimiterService))

	http.ListenAndServe(configs.WebServerHost+":"+configs.WebServerPort, mux)
}

func RateLimiter(next http.Handler, rateLimiterService ratelimiter.RateLimiterService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		ok := rateLimiterService.Execute(ip, "ip")
		if !ok {
			http.Error(w, RATE_LIMIT_ERROR_MESSAGE, http.StatusTooManyRequests)
			return
		}

		// salvar request no REDIS
		next.ServeHTTP(w, r)
	})
}

func HelloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		w.Write([]byte("Hello World"))
	})
}
