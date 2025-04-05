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

func (c *Cache) Increment(key string, number int64) error {
	c.Lock()
	val, found := c.items[key]

	if !found {
		c.Unlock()
		return errors.New("element to increment not found")
	}

	switch val.Value.(type) {
	case int:
		val.Value = val.Value.(int) + int(number)
	case int8:
		val.Value = val.Value.(int8) + int8(number)
	case int16:
		val.Value = val.Value.(int16) + int16(number)
	case int32:
		val.Value = val.Value.(int32) + int32(number)
	case int64:
		val.Value = val.Value.(int64) + int64(number)
	case uint:
		val.Value = val.Value.(uint) + uint(number)
	case uint8:
		val.Value = val.Value.(uint8) + uint8(number)
	case uint16:
		val.Value = val.Value.(uint16) + uint16(number)
	case uint32:
		val.Value = val.Value.(uint32) + uint32(number)
	case uint64:
		val.Value = val.Value.(uint64) + uint64(number)
	case float32:
		val.Value = val.Value.(float32) + float32(number)
	case float64:
		val.Value = val.Value.(float64) + float64(number)
	default:
		c.Unlock()
		return errors.New("the value is not and integer/float")
	}

	c.items[key] = val

	c.Unlock()

	return nil
}

// функция для копирования элемента
// TODO: написать ее
// func (c *Cache) Copy(key string) (interface{}, bool) {
//
// }

// функция для проверки существования элемента в кэше
func (c *Cache) Exist(key string) bool {
	c.RLock()
	defer c.Unlock()

	_, exists := c.items[key]

	return exists
}

// очистка всех файлов из кэша
func (c *Cache) FlushAll() {
	c.Lock()
	defer c.Unlock()

	c.items = map[string]CacheItem{}
}
