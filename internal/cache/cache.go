package cache

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/yourusername/getemps-service/internal/models"
)

type Cache interface {
	Get(key string) (*models.EmployeeInfo, bool)
	Set(key string, value *models.EmployeeInfo, duration time.Duration)
	Delete(key string)
}

type inMemoryCache struct {
	cache *cache.Cache
}

func NewInMemoryCache(defaultExpiration, cleanupInterval time.Duration) Cache {
	return &inMemoryCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *inMemoryCache) Get(key string) (*models.EmployeeInfo, bool) {
	if item, found := c.cache.Get(key); found {
		if data, ok := item.([]byte); ok {
			var employeeInfo models.EmployeeInfo
			if err := json.Unmarshal(data, &employeeInfo); err == nil {
				return &employeeInfo, true
			}
		}
	}
	return nil, false
}

func (c *inMemoryCache) Set(key string, value *models.EmployeeInfo, duration time.Duration) {
	if data, err := json.Marshal(value); err == nil {
		c.cache.Set(key, data, duration)
	}
}

func (c *inMemoryCache) Delete(key string) {
	c.cache.Delete(key)
}

func GenerateCacheKey(nationalNumber string) string {
	return "emp_status:" + nationalNumber
}