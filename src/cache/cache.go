package cache

import (
	"container/list"
	"errors"
	"sync"
)

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrKeyTooLarge   = errors.New("key exceeds maximum length")
	ErrValueTooLarge = errors.New("value exceeds maximum length")
)

type cacheItem struct {
	key   string
	value string
}

type Cache struct {
	data       map[string]*list.Element
	lru        *list.List
	mu         sync.RWMutex
	maxSize    int
	maxKeySize int
	maxValSize int
}

func NewCache(maxSize, maxKeySize, maxValueSize int) *Cache {
	return &Cache{
		data:       make(map[string]*list.Element),
		lru:        list.New(),
		maxSize:    maxSize,
		maxKeySize: maxKeySize,
		maxValSize: maxValueSize,
	}
}

func (c *Cache) Put(key, value string) error {
	if len(key) > c.maxKeySize {
		return ErrKeyTooLarge
	}

	if len(value) > c.maxValSize {
		return ErrValueTooLarge
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.data[key]; exists {
		c.lru.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
		return nil
	}

	if c.lru.Len() >= c.maxSize {
		c.evict()
	}

	elem := c.lru.PushFront(&cacheItem{key: key, value: value})
	c.data[key] = elem
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	c.mu.RLock()
	elem, ok := c.data[key]
	if !ok {
		c.mu.RUnlock()
		return "", ErrKeyNotFound
	}

	value := elem.Value.(*cacheItem).value
	c.mu.RUnlock()

	c.mu.Lock()
	c.lru.MoveToFront(elem)
	c.mu.Unlock()

	return value, nil
}

func (c *Cache) evict() {
	if elem := c.lru.Back(); elem != nil {
		c.lru.Remove(elem)
		delete(c.data, elem.Value.(*cacheItem).key)
	}
}
