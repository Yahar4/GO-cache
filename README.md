# GO-cache 🚀

[![Go Reference](https://pkg.go.dev/badge/github.com/Yahar4/GO-cache.svg)](https://pkg.go.dev/github.com/Yahar4/GO-cache)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Yahar4/GO-cache/blob/main/LICENSE)

**GO-cache** – это простой и эффективный кэш на языке Go, предназначенный для хранения временных данных в памяти с поддержкой TTL (Time-To-Live). Подходит для использования в веб-приложениях, микросервисах и других проектах, где требуется быстрое кэширование данных.

## Установка


```bash
go get -u github.com/Yahar4/GO-cache
```

## Пример использования
```go
import (
	"github.com/Yahar4/GO-cache"
)

func main() {
	// создание кэша с дефолтным временем хранения
	// и временем удаления раз в 10 минут
	c := cache.New(5*time.Minute, 10*time.Minute)

	// установка значений в кэш по ключу "foo"
	c.Set("foo", "bar", cache.DefaultExpiration)

	// получение значения из кэша по ключу
	value := c.Get("foo")
}
```
