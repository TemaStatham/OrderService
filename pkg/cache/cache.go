package cache

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/TemaStatham/OrderService/pkg/service"
)


const (
	countOfRestoredItems       = 10
	lifetimeElementInsideCache = 60
)

// Cache : структура описывающуя наш контейнер-хранилище
type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
	service           *service.Service
}

// Item : структура для элемента
type Item struct {
	Value      *model.OrderClient
	Created    time.Time
	Expiration int64
}

// New : инициализация нового контейнера-хранилища
func New(defaultExpiration, cleanupInterval time.Duration, service *service.Service) *Cache {
	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		service:           service,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC()
	}

	// Инициализация кэша из базы данных
    if err := cache.RestoreCache(); err != nil {
        log.Printf("Failed to restore data from database: %v", err)
    }

	return &cache
}

// Set : возможность записывать данные в кэш
func (c *Cache) Set(key string, value *model.OrderClient, duration time.Duration) {
	var expiration int64

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	if _, exists := c.items[key]; exists {
		log.Printf("key : %s, is exist\n", key)
		return
	}

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

// Get : метод для получения значений
func (c *Cache) Get(key string) (*model.OrderClient, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	return item.Value, true
}

// Delete : возможность удалить кеш
func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

// StartGC : поиск просроченных ключей с последующей очисткой (GC)
// Для этого напишем метод StartGC, который запускается при инициализация нового экземпляра кеша New и работает пока программа не будет завершена.
func (c *Cache) StartGC() {
	go c.GC()
}

// GC проверяет истекшие ключи, записывает их в БД и удаляет из кеша
func (c *Cache) GC() {
	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.writeToDatabase(keys)
			//c.clearItems(keys)
		}

	}
}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// writeToDatabase записывает просроченные ключи в БД
func (c *Cache) writeToDatabase(keys []string) {

	c.Lock()

	defer c.Unlock()

	fmt.Println("write")

	for _, key := range keys {
		value, found := c.items[key]
		if !found {
			continue
		}

		_, err := c.service.AddOrder(value.Value)
		if err != nil {
			log.Print(err)
		}

		delete(c.items, key)
	}
}

// // clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
// func (c *Cache) clearItems(keys []string) {

// 	c.Lock()

// 	defer c.Unlock()

// 	for _, k := range keys {
// 		delete(c.items, k)
// 	}
// }

// RestoreCache получает данные из базы данных и инициализирует кэш
func (c *Cache) RestoreCache() error {
	orders, err := c.service.GetRecentOrders(countOfRestoredItems)
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	for _, order := range orders {
		c.items[order.OrderUID] = Item{
			Value:      order,
			Expiration: lifetimeElementInsideCache,
			Created:    time.Now(),
		}
	}

	return nil
}
