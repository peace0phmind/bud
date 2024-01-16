package util

import (
	"github.com/peace0phmind/bud/stream"
	"sync"
)

type Cache[K comparable, V any] struct {
	cacheMap sync.Map
}

type CachePair[K comparable, V any] struct {
	K K
	V V
}

type cacheItem[V any] struct {
	value      V
	valueValid bool
	lock       sync.Mutex
}

func (ci *cacheItem[V]) getValue() (v V, got bool) {
	if !ci.valueValid {
		got = false
		return
	}

	return ci.value, true
}

func (c *Cache[K, V]) Get(k K) (v V, got bool) {
	value, loaded := c.cacheMap.Load(k)
	if !loaded {
		got = false
		return
	} else {
		ci := value.(*cacheItem[V])
		return ci.getValue()
	}
}

func (c *Cache[K, V]) GetOrDefault(k K, def V) (v V, got bool) {
	value, got := c.cacheMap.Load(k)
	if got {
		ci := value.(*cacheItem[V])
		return ci.getValue()
	} else {
		return def, false
	}
}

func (c *Cache[K, V]) GetOrStore(k K, v V) (ov V, got bool) {
	ci, got := c.cacheMap.LoadOrStore(k, &cacheItem[V]{value: v, valueValid: true})
	if got {
		return ci.(*cacheItem[V]).getValue()
	} else {
		got = false
		return
	}
}

func (c *Cache[K, V]) GetAndDelete(k K) (v V, ok bool) {
	value, loaded := c.cacheMap.LoadAndDelete(k)
	if !loaded {
		ok = false
		return
	} else {
		ci := value.(*cacheItem[V])
		return ci.getValue()
	}
}

func (c *Cache[K, V]) Set(k K, v V) (oldValue V, ok bool) {
	actual, _ := c.cacheMap.LoadOrStore(k, &cacheItem[V]{})
	ci := actual.(*cacheItem[V])
	ci.lock.Lock()
	defer ci.lock.Unlock()

	oldValue, ok = ci.getValue()
	ci.value = v
	ci.valueValid = true
	return
}

func (c *Cache[K, V]) Delete(k K) {
	c.cacheMap.Delete(k)
}

func (c *Cache[K, V]) Range(rangeFunc func(k K, v V) bool) {
	c.cacheMap.Range(func(key, value any) bool {
		k := key.(K)
		ci := value.(*cacheItem[V])
		if ci.valueValid {
			return rangeFunc(k, ci.value)
		}
		return true
	})
}

func (c *Cache[K, V]) Filter(filterFunc func(k K, v V) bool) *Cache[K, V] {
	filteredCache := &Cache[K, V]{cacheMap: sync.Map{}}
	c.cacheMap.Range(func(key, value any) bool {
		k := key.(K)
		ci := value.(*cacheItem[V])
		if ci.valueValid && filterFunc(k, ci.value) {
			filteredCache.Set(k, ci.value)
		}
		return true
	})
	return filteredCache
}

func (c *Cache[K, V]) FilterToStream(filterFunc func(k K, v V) bool) stream.Stream[*CachePair[K, V]] {
	result := stream.Stream[*CachePair[K, V]]{}

	c.cacheMap.Range(func(key, value any) bool {
		k := key.(K)
		ci := value.(*cacheItem[V])
		if ci.valueValid && filterFunc(k, ci.value) {
			result = result.Append(&CachePair[K, V]{K: k, V: ci.value})
		}
		return true
	})

	return result
}

func (c *Cache[K, V]) ToStream() stream.Stream[*CachePair[K, V]] {
	result := stream.Stream[*CachePair[K, V]]{}

	c.cacheMap.Range(func(key, value any) bool {
		k := key.(K)
		ci := value.(*cacheItem[V])
		if ci.valueValid {
			result = result.Append(&CachePair[K, V]{K: k, V: ci.value})
		}
		return true
	})

	return result
}

func (c *Cache[K, V]) Size() int {
	size := 0
	c.cacheMap.Range(func(key, value any) bool {
		size += 1
		return true
	})

	return size
}

func (c *Cache[K, V]) GetOrNew(k K, newFunc func() (V, error)) (v V, err error) {
	if newFunc == nil {
		panic("new function must not be null")
	}

	actual, loaded := c.cacheMap.LoadOrStore(k, &cacheItem[V]{})
	ci := actual.(*cacheItem[V])
	if loaded {
		if ci.valueValid {
			return ci.value, nil
		}
	}

	ci.lock.Lock()
	defer ci.lock.Unlock()

	if ci.valueValid {
		return ci.value, nil
	}

	defer func() {
		if err == nil {
			ci.value = v
			ci.valueValid = true
		}
	}()

	return newFunc()
}

func (c *Cache[K, V]) ToMap() map[K]V {
	m := make(map[K]V)
	c.cacheMap.Range(func(key, value any) bool {
		k := key.(K)
		ci := value.(*cacheItem[V])
		if ci.valueValid {
			m[k] = ci.value
		}
		return true
	})
	return m
}
