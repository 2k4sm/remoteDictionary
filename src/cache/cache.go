package cache

import (
	"container/list"
	"errors"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
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
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	value := elem.Value.(*cacheItem).value

	c.lru.MoveToFront(elem)

	return value, nil
}

func getTotalMemory() uint64 {
	const defaultMemory uint64 = 2 * 1024 * 1024 * 1024
	
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error retrieving memory info: %v", err)
		return defaultMemory
	}
	
	return vmStat.Total
}

func (c *Cache) evict(threshold uint64) {
	batchSize := 5
	for {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		if memStats.Alloc <= threshold/2 {
			log.Printf("Memory reduced to %d MB (below target %d MB), stopping eviction",
				memStats.Alloc/1024/1024, threshold/2/1024/1024)
			break
		}

		c.mu.Lock()
		if c.lru.Len() == 0 {
			c.mu.Unlock()
			log.Println("Cache empty, cannot evict further")
			break
		}

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
		batchSize = min(batchSize*2, 1000)

		runtime.GC()
		time.Sleep(50 * time.Millisecond)
	}
}

func (c *Cache) MonitorMemoryUsage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	log.Println("Monitoring memory usage...")
	for range ticker.C {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		totalMem := getTotalMemory()
		threshold := totalMem * 70 / 100 // 70% threshold

		memUsageMB := memStats.Alloc / 1024 / 1024
		thresholdMB := threshold / 1024 / 1024

		if memStats.Alloc > threshold {
			log.Printf("Memory usage Critical: %d MB used (threshold: %d MB).", memUsageMB, thresholdMB)
			go c.evict(threshold)
		}
	}
}
