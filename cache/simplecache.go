package cache

import (
	"fmt"
	"sync"
	"time"
)

//一个简单的cache实现,够用就好.根据需求实现LRU/LFU
type CacheCtl interface {
	Add(key interface{}, item CacheItem) error
	Get(key interface{}) (CacheItem, error)
	Update(key interface{}, item CacheItem)
	Clear()
}

var ErrCacheNotFound error = fmt.Errorf("cache item is not found")
var ErrCacheExist error = fmt.Errorf("cache item is exist")

type CacheItem struct {
	Data interface{}

	expire time.Duration //0 is infinity,Not less than 1 second

	createTime time.Time
}

func NewCacheItem(data interface{}, expire time.Duration) CacheItem {
	r := CacheItem{Data: data}
	return r.Expire(expire)
}

type simpleCache struct {
	m           sync.RWMutex
	sm          sync.Map
	TimeoutFunc func(CacheItem)
}

func (t CacheItem) Expire(d time.Duration) CacheItem {
	if d > 0 {
		v := d
		if v < time.Second {
			v = time.Second
		}
		t.expire = v
	}
	return t
}

func (t *simpleCache) Add(key interface{}, item CacheItem) error {
	if _, ok := t.sm.Load(key); !ok {
		item.createTime = time.Now().Local()
		t.sm.Store(key, item)
	} else {
		return ErrCacheExist
	}
	return nil
}

func (t *simpleCache) Get(key interface{}) (CacheItem, error) {
	if v, ok := t.sm.Load(key); ok {
		ci := v.(CacheItem)
		ci.createTime = time.Now().Local()
		return ci, nil
	} else {
		return CacheItem{}, ErrCacheNotFound
	}
}

func (t *simpleCache) Update(key interface{}, item CacheItem) {
	if _, ok := t.sm.LoadOrStore(key, item); ok {
		//if loaded then update
		item.createTime = time.Now().Local()
		t.sm.Store(key, item)
	}
}

func (t *simpleCache) Clear() {
	t.sm = sync.Map{}
}

var sc *simpleCache

func GetInstance() *simpleCache {
	return sc
}

func init() {
	sc = &simpleCache{}
	go checkExpire(sc)
}

func checkExpire(c *simpleCache) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		now := time.Now().Local()
		c.sm.Range(func(key, value interface{}) bool {
			v := value.(CacheItem)
			if v.expire > 0 {
				if now.Sub(v.createTime) > v.expire {
					if c.TimeoutFunc != nil {
						c.TimeoutFunc(v)
					}
					c.sm.Delete(key)
				}
			}
			return true
		})
	}
}
