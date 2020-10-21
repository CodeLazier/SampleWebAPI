package tests

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
	"unsafe"

	"test/cache"
	"test/handler"
	"test/queue"
	"test/ratelimite"
)

type Snowflake struct {
	lastTimestamp int64
	index         int16
	machId        int16
}

func (s *Snowflake) Init(id int16) error {
	if id > 0xff {
		return errors.New("illegal machine id")
	}

	s.machId = id
	s.lastTimestamp = time.Now().UnixNano() / 1e6
	s.index = 0
	return nil
}

func (s *Snowflake) GetId() (int64, error) {
	curTimestamp := time.Now().UnixNano() / 1e6
	if curTimestamp == s.lastTimestamp {
		s.index++
		if s.index > 0xfff {
			s.index = 0xfff
			return -1, errors.New("out of range")
		}
	} else {
		//fmt.Printf("id/ms:%d -- %d\n", s.lastTimestamp, s.index)
		s.index = 0
		s.lastTimestamp = curTimestamp
	}
	return (0x1ffffffffff&s.lastTimestamp)<<22 + int64(0xff<<10) + int64(0xfff&s.index), nil
}

func Test_InterfaceCompatible(t *testing.T) {
	var _ handler.Control = &handler.MsgDB{}
	_ = &handler.CtlMock{}
}

func TestCache(t *testing.T) {
	c := cache.GetInstance()
	c2 := cache.GetInstance()
	if unsafe.Pointer(c) != unsafe.Pointer(c2) {
		log.Fatal("instance difference!")
	}

	c.TimeoutFunc = func(item cache.CacheItem) {
		fmt.Println("timeout removed:", item)
	}
	rand.Seed(time.Now().Unix())
	IDs := make([]int64, 0)
	snow := &Snowflake{}
	_ = snow.Init(1)
	m := sync.RWMutex{}
	wg := sync.WaitGroup{}
	getID := func() int64 {
		m.Lock()
		defer m.Unlock()
		id, _ := snow.GetId()
		IDs = append(IDs, id)
		return id
	}
	ch := make(chan struct{})
	f := func() {
		for {
			select {
			case <-ch:
				return
			default:
				newId := getID()
				if err := c.Add(fmt.Sprintf("test_%d", newId), cache.NewCacheItem(fmt.Sprintf("value_%d", newId), time.Duration(rand.Intn(3))*time.Second)); err != nil {
					t.Fatal(err)
				}

				time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
			}
		}
	}

	f2 := func() {
		defer wg.Done()
		tm := time.After(6 * time.Second)
		for {
			select {
			case <-tm:
				return
			default:
				m.RLock()
				//rand.Intn param is not zero
				var newId int64
				l := len(IDs)
				if l == 1 {
					newId = IDs[0]
				} else if l > 1 {
					newId = IDs[rand.Intn(l-1)]
				}
				m.RUnlock()

				if _, err := c.Get(fmt.Sprintf("test_%d", newId)); err == cache.ErrCacheNotFound {
					t.Log(err)
				}

				time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
			}
		}
	}

	_ = func(value interface{}, err error) {
		if err != nil {
			log.Print(err)
		} else {
			log.Print(value)
		}
	}

	const count = 100

	for i := 0; i < count; i++ {
		go f()
	}

	for i := 0; i < count/10; i++ {
		wg.Add(1)
		go f2()
	}

	wg.Wait()
	close(ch)
	//verify result
}

func TestBatchQueueCreate(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if queue.GetInstance() != queue.GetInstance() {
				t.Log(errors.New("instance error"))
			}

		}()
	}
	wg.Wait()
}

func TestBatchQueue(t *testing.T) {
	q := queue.GetInstance()
	q.Push(2, 30, 4, 5, 6, 7, 8)
	q.SetDoFun(func(v []interface{}) error {
		t.Log(v)
		return nil
	}, false)
	time.Sleep(2 * time.Second)
	q.Push("a", "b", "c", "d")
	time.Sleep(2 * time.Second)

	for i := 0; i < 100; i++ {
		q.Push(rand.Intn(999))
	}
	time.Sleep(2 * time.Second)

}

func TestRateLimite(t *testing.T) {
	r := ratelimite.NewTokenBucket(10)

	for i := 0; i < 10; i++ {
		r.RequestToken(3)
		log.Println("get token is ok")
	}
}
