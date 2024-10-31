package rate_limiter

import (
	"fmt"
	"time"
)

type RateLimiterConfig struct {
	LimitPerInterval  int
	IntervalInSeconds int
	BanTimeInSeconds  int
}

type RateLimiterItem struct {
	Key             string
	LastRequestDate time.Time
	RequestsCount   int
	BannedUntil     time.Time
}

type RateLimiterService interface {
	Execute(itemKey, configType string) bool
	GetConfig(name string) RateLimiterConfig
	AddConfig(configType string, config RateLimiterConfig)
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

func (r *RateLimiterServiceImpl) SetRateLimiterItemRepository(rateLimiterItemRepository RateLimiterItemRepository) {
	r.rateLimiterItemRepository = rateLimiterItemRepository
}

func (r *RateLimiterServiceImpl) Execute(itemKey, configType string) bool {
	config := r.GetConfig(configType)
	fmt.Printf("\n\nExecuting item %s with config %s \n", itemKey, configType)
	fmt.Println("Config:", config)

	timeNow := time.Now()
	fmt.Println("Time now:", timeNow)

	item := r.rateLimiterItemRepository.Find(itemKey)
	if item == nil {
		fmt.Println("Creating new item")
		item = &RateLimiterItem{
			Key:             itemKey,
			LastRequestDate: timeNow,
			RequestsCount:   0,
		}
	}

	isBanned := item.BannedUntil.After(timeNow)
	if isBanned {
		fmt.Printf("Item is banned: count %d \n", item.RequestsCount)
		return false
	}

	if time.Since(item.LastRequestDate).Seconds() > float64(config.IntervalInSeconds) {
		fmt.Println("Resetting requests count")
		item.RequestsCount = 0
		item.LastRequestDate = timeNow
	}

	item.RequestsCount++

	if item.RequestsCount > config.LimitPerInterval {
		fmt.Printf("Banning item. count %d. limit %d \n", item.RequestsCount, config.LimitPerInterval)
		item.BannedUntil = timeNow.Add(time.Duration(config.BanTimeInSeconds) * time.Second)
	}

	r.rateLimiterItemRepository.Save(item)
	fmt.Printf("Saving item. count %d. limit %d \n", item.RequestsCount, config.LimitPerInterval)

	return item.RequestsCount <= config.LimitPerInterval
}

func (r *RateLimiterServiceImpl) GetConfig(name string) RateLimiterConfig {
	return r.configs[name]
}

func (r *RateLimiterServiceImpl) AddConfig(configType string, config RateLimiterConfig) {
	r.configs[configType] = config
}
