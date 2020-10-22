package memory

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
	"github.com/zfd81/magpie/config"
)

type Cache struct {
	c *cache.Cache
}

func (c *Cache) Set(key string, value interface{}) *Cache {
	c.c.Set(key, value, cache.NoExpiration)
	return c
}

func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) *Cache {
	c.c.Set(key, value, ttl)
	return c
}

func (c *Cache) Remove(key string) *Cache {
	c.c.Delete(key)
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.c.Get(key)
}

func (c *Cache) GetWithExpiration(key string) (interface{}, time.Time, bool) {
	return c.c.GetWithExpiration(key)
}

func (c *Cache) GetString(key string) string {
	value, found := c.Get(key)
	if found {
		return cast.ToString(value)
	}
	return ""
}

func (c *Cache) GetInt(key string) int {
	value, found := c.Get(key)
	if found {
		return cast.ToInt(value)
	}
	return 0
}

func (c *Cache) GetBool(key string) bool {
	value, found := c.Get(key)
	if found {
		return cast.ToBool(value)
	}
	return false
}

func (c *Cache) Count() int {
	return c.c.ItemCount()
}

func New() *Cache {
	conf := config.GetConfig()
	c := &Cache{
		c: cache.New(conf.Memory.ExpirationTime, conf.Memory.CleanupInterval),
	}
	return c
}
