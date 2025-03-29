package cache

import (
	"container/list"
	"errors"
	"log"
	"runtime"
	"sync"
	"syscall"
	"time"
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
	maxKeySize int
	maxValSize int
}

func NewCache(maxKeySize, maxValueSize int) *Cache {
	return &Cache{
		data:       make(map[string]*list.Element),
		lru:        list.New(),
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

func getTotalMemory() uint64 {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		log.Printf("Error retrieving system info: %v", err)
		return 2 * 1024 * 1024 * 1024
	}
	return info.Totalram * uint64(info.Unit)
}

func (c *Cache) evict(threshold uint64) {
	batchSize := 5
	for {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		if memStats.Alloc <= threshold-threshold/2 {
			batchSize = 5
			break
		}

		c.mu.Lock()
		removed := 0
		for removed < batchSize && c.lru.Len() > 0 {
			elem := c.lru.Back()
			if elem == nil {
				break
			}
			c.lru.Remove(elem)
			delete(c.data, elem.Value.(*cacheItem).key)
			removed++
		}
		c.mu.Unlock()

		time.Sleep(100 * time.Millisecond)
		batchSize *= 10
	}
}

func (c *Cache) MonitorMemoryUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	log.Println("Monitoring memory usage...")
	for range ticker.C {
		totalMem := getTotalMemory()
		threshold := totalMem * 70 / 100 // 70% threshold

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		log.Println("nowMemory:", memStats.Alloc/1024/1024)
		log.Println("threshold:", threshold/1024/1024)
		if memStats.Alloc > threshold {
			log.Printf("High memory usage detected: %d bytes allocated (threshold: %d bytes). Initiating eviction...", memStats.Alloc, threshold)
			c.evict(threshold)
		}
	}
}
