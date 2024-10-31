package rate_limiter

type RateLimiterItemRepository interface {
	Find(key string) *RateLimiterItem
	Save(item *RateLimiterItem)
}

type RateLimiterItemRepositoryMemory struct {
	items map[string]*RateLimiterItem
}

func NewRateLimiterItemRepositoryMemory() *RateLimiterItemRepositoryMemory {
	return &RateLimiterItemRepositoryMemory{
		items: make(map[string]*RateLimiterItem),
	}
}

func (r *RateLimiterItemRepositoryMemory) Find(key string) *RateLimiterItem {
	return r.items[key]
}

func (r *RateLimiterItemRepositoryMemory) Save(item *RateLimiterItem) {
	r.items[item.Key] = item
}
