package memory

import (
	"strings"
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

func (c *Cache) RemoveWithPrefix(prefix string) int {
	cnt := 0
	for k, _ := range c.c.Items() {
		if strings.HasPrefix(k, prefix) {
			c.c.Delete(k)
			cnt++
		}
	}
	return cnt
}

func (c *Cache) Iterator(f func(k string, v interface{}) error) error {
	for k, v := range c.c.Items() {
		err := f(k, v.Object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.c.Get(key)
}

func (c *Cache) GetWithExpiration(key string) (interface{}, time.Time, bool) {
	return c.c.GetWithExpiration(key)
}

func (c *Cache) GetByte(key string) []byte {
	value, found := c.Get(key)
	if found {
		return value.([]byte)
	}
	return nil
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

func (c *Cache) GetSlice(key string) []interface{} {
	value, found := c.Get(key)
	if found {
		return cast.ToSlice(value)
	}
	return []interface{}{}
}

func (c *Cache) Clear() {
	c.c.Flush()
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

func Size(data interface{}) int {
	switch data.(type) {
	case string:
		return len(data.(string))
	case int64:
		return 8
	case int32:
		return 4
	case int16:
		return 2
	case float64:
		return 8
	case float32:
		return 4
	case bool:
		return 1
	default:
		return 1
	}
}
