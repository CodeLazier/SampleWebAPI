/*
 * @Author: your name
 * @Date: 2020-09-25 12:02:36
 * @LastEditTime: 2020-09-25 14:57:07
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\cache\simplecache.go
 */
package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheCtl interface {
	Add(key interface{}, item CacheItem) error
	Get(key interface{}) (CacheItem, error)
	Update(key interface{}, item CacheItem)
}

var ErrCacheNotFound error = fmt.Errorf("cache item is not found")
var ErrCacheExist error = fmt.Errorf("cache item is exist")

type CacheItem struct {
	Data interface{}

	Expire time.Duration //0 is infinity,Not less than 1 second

	createTime time.Time
}

type simpleCache struct {
	m  sync.RWMutex
	sm sync.Map
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

var sc *simpleCache

func GetInstance() *simpleCache {
	return sc
}

func init() {
	sc = &simpleCache{}
	go checkExpire()
}

func checkExpire() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		now := time.Now().Local()
		sc.sm.Range(func(key, value interface{}) bool {
			v := value.(CacheItem)
			if v.Expire > 0 {
				if now.Sub(v.createTime) > v.Expire {
					sc.sm.Delete(key)
				}
			}
			return true
		})
	}
}
