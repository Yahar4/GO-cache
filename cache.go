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

// получение всех значений из кэша
func (c *Cache) GetAllItems() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()

	allItems := make(map[string]interface{})

	for key, item := range c.items {
		if time.Now().UnixNano() < item.Expiration {
			allItems[key] = item.Value
		}
	}

	return allItems
}

// получение конкретного значения по ключу
func (c *Cache) GetItem(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	// получаем значения по ключу
	item, found := c.items[key]
	// проверка на то, что значение существует
	if !found {
		return nil, false
	}

	// проверка на то, ликвидны ли наши данные
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

// подсчет кол-ва элементов в кэше
func (c *Cache) Count() int64 {
	c.RLock()
	defer c.RUnlock()

	numberOfItems := len(c.items)

	return int64(numberOfItems)
}

func (c *Cache) RenameKey(oldKey, newKey string) error {
	c.Lock()
	defer c.Unlock()

	// проверка существования ключа -> если существует, то меняем его и удаляем старый
	if item, exists := c.items[oldKey]; exists {
		// обработка случая, если наш новый ключ уже существует
		if _, exists := c.items[newKey]; exists {
			return errors.New("new key already exists in your cache")
		}

		c.items[newKey] = item
		delete(c.items, oldKey)
	} else {
		return errors.New("no such key in cache")
	}

	return nil
}
