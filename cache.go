package cache

import "time"

func (c *Cache) Set(key string, val interface{}, duration time.Duration) {}

func (c *Cache) Get(key string) (interface{}, bool) {}

func (c *Cache) Delete(key string) {}
