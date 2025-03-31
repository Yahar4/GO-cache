package cache

import (
	"errors"
	"time"
)

// установка значений в хранилище
func (c *Cache) Set(key string, val interface{}, duration time.Duration) {
	var exp int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		exp = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = CacheItem{
		Value:      val,
		CreatedAt:  time.Now(),
		Expiration: exp,
	}
}

// получение значений из хранилища
func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	// получаем значения по ключу
	item, found := c.items[key]
	// проверка на то, что значение существует
	if !found {
		return nil, false
	}

	// проверка на то, не истекли ли наши данные
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	// возвращаем все значения
	return item.Value, true
}

// удаление кэша
func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}
