package infra

import (
	"encoding/json"

	"github.com/KelpGF/Go-Rate-Limiter/package/rate_limiter"
)

type RateLimiterItemRedisRepository struct {
	redis *Redis
}

func NewRateLimiterItemRedisRepository(redis *Redis) *RateLimiterItemRedisRepository {
	return &RateLimiterItemRedisRepository{
		redis: redis,
	}
}

func (r *RateLimiterItemRedisRepository) Find(key string) *rate_limiter.RateLimiterItem {
	val, err := r.redis.Get(r.makeKey(key))

	if err != nil {
		return nil
	}

	var item rate_limiter.RateLimiterItem

	err = json.Unmarshal([]byte(val), &item)

	if err != nil {
		return nil
	}

	return &item

}

func (r *RateLimiterItemRedisRepository) Save(item *rate_limiter.RateLimiterItem) {
	key := r.makeKey(item.Key)
	val, err := json.Marshal(item)

	if err != nil {
		return
	}

	r.redis.Set(key, val)
}

func (r *RateLimiterItemRedisRepository) makeKey(itemKey string) string {
	return "rate_limiter:" + itemKey
}
