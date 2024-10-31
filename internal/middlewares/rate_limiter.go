package middlewares

import (
	"net"
	"net/http"

	"github.com/KelpGF/Go-Rate-Limiter/package/rate_limiter"
)

const RATE_LIMIT_ERROR_MESSAGE = "you have reached the maximum number of requests or actions allowed within a certain time frame"

type RateLimiterMiddleware struct {
	rateLimiterService rate_limiter.RateLimiterService
	rateLimiterToken   string
}

func NewRateLimiterMiddleware(rateLimiterService rate_limiter.RateLimiterService, rateLimiterToken string) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		rateLimiterService: rateLimiterService,
		rateLimiterToken:   rateLimiterToken,
	}
}

func (h *RateLimiterMiddleware) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		token := r.Header.Get("API_KEY")

		var ok bool

		if token == h.rateLimiterToken {
			ok = h.rateLimiterService.Execute(token, "token")
		} else {
			ok = h.rateLimiterService.Execute(ip, "ip")
		}

		if !ok {
			http.Error(w, RATE_LIMIT_ERROR_MESSAGE, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
