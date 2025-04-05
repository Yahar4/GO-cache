package cache

import (
	"sync"
	"time"
)

// тип данных для значений кэша
type CacheItem struct {
	// мы можем хранить любое значение, поэтому нам нужен тип interface
	Value interface{}
	// время создания кэша
	CreatedAt time.Time
	// истечение времени для кэша, unix
	Expiration int64
}

// in-memory кэш
type Cache struct {
	sync.RWMutex

	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]CacheItem
}

type CacheMethods interface {
	Set(key string, value interface{}, duration time.Duration)
	GetAllItems() map[string]interface{}
	GetItem(key string) (interface{}, bool)
	Delete(key string) error
	Count() int64
	RenameKey(oldKey, newKey string) error
}

// конструктор
func New(deafultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]CacheItem)

	cache := Cache{
		defaultExpiration: deafultExpiration,
		cleanupInterval:   cleanupInterval,
		items:             items,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

// сборка мусора
func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) == 0 {
			c.clearItems(keys)
		}
	}
}

// просроченные ключи
func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.Unlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// убираем просроченные элементы
func (c *Cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
