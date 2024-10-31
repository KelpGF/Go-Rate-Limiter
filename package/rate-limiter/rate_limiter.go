package ratelimiter

import "time"

type RateLimiterConfig struct {
	LimitPerInterval  int
	IntervalInSeconds int
}

type RateLimiterItem struct {
	Key             string
	LastRequestDate time.Time
	RequestsCount   int
}

type RateLimiterService interface {
	Execute(itemKey, configType string) bool
	GetConfig(name string) RateLimiterConfig
	SetConfig(configType string, config RateLimiterConfig)
}

type RateLimiterServiceImpl struct {
	configs                   map[string]RateLimiterConfig
	rateLimiterItemRepository RateLimiterItemRepository
}

func NewRateLimiterServiceImpl(rateLimiterItemRepository RateLimiterItemRepository) *RateLimiterServiceImpl {
	return &RateLimiterServiceImpl{
		configs:                   make(map[string]RateLimiterConfig),
		rateLimiterItemRepository: rateLimiterItemRepository,
	}
}

func (r *RateLimiterServiceImpl) Execute(itemKey, configType string) bool {
	config := r.GetConfig(configType)

	item := r.rateLimiterItemRepository.Find(itemKey)
	if item == nil {
		item = &RateLimiterItem{
			Key:             itemKey,
			LastRequestDate: time.Now(),
			RequestsCount:   0,
		}
	}

	if time.Since(item.LastRequestDate).Seconds() > float64(config.IntervalInSeconds) {
		item.RequestsCount = 0
		item.LastRequestDate = time.Now()
	}

	item.RequestsCount++
	r.rateLimiterItemRepository.Save(item)

	return item.RequestsCount <= config.LimitPerInterval
}

func (r *RateLimiterServiceImpl) GetConfig(name string) RateLimiterConfig {
	return r.configs[name]
}

func (r *RateLimiterServiceImpl) SetConfig(configType string, config RateLimiterConfig) {
	r.configs[configType] = config
}
