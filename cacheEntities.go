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
	Get(key string) (interface{}, bool)
	Delete(key string) error
	Count() int64
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
		// TODO: прописать удаление устаревших элементов
	}

	return &cache
}
