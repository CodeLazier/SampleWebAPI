/*
 * @Author: your name
 * @Date: 2020-09-22 11:57:35
 * @LastEditTime: 2020-09-28 22:19:54
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\tests\msg_test.go
 */
package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"test/cache"
	"test/handler"
	"test/msg"

	"testing"
	"time"
	"unsafe"
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
	return int64((0x1ffffffffff&s.lastTimestamp)<<22) + int64(0xff<<10) + int64(0xfff&s.index), nil
}

func Test_InterfaceCompatible(t *testing.T) {
	// var _ handler.Control = &handler.Eip{}
	// _ = &handler.CtlMock{}
}

func Benchmark_Add_Get(b *testing.B) {
	c := cache.GetInstance()

	wg := sync.WaitGroup{}
	w := func(key int, ca chan int) <-chan int {
		ca <- key
		if err := c.Add(key, cache.CacheItem{}); err == cache.ErrCacheExist {
			b.Fatal(err)
		}
		return ca
	}

	r := func(ca <-chan int) {
		if _, err := c.Get(<-ca); err == cache.ErrCacheNotFound {
			b.Fatal(err)
		}
	}

	wg.Add(b.N)

	ch := make(chan int, 1)
	for i := 0; i < b.N; i++ {
		go func(i int) {
			defer wg.Done()
			r(w(i, ch))
		}(i)
	}

	wg.Wait()

	c.Clear()
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
				newId := IDs[rand.Intn(len(IDs))]
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

	const count = 1000

	for i := 0; i < count; i++ {
		go f()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go f2()
	}

	wg.Wait()
	close(ch)
	//verify result
}

func TestMockOrm(t *testing.T) {
	eip := &msg.EipMsgHandler{
		Control: handler.NewCtlMock(),
	}
	if msgs, err := eip.GetAll(0, -1); err != nil {
		t.Fail()
		t.Log(err)
	} else {
		t.Log(msgs)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	count := 0
	defer ticker.Stop()
	getIndex := func() {
		if msg, err := eip.GetIndex(0); err != nil {
			t.Fail()
			t.Log(err)
		} else {
			t.Log(msg)
		}
	}
	for {
		<-ticker.C
		getIndex()
		count++
		if count > 5 {
			break
		}
	}

	if err := eip.MarkRead(3); err != nil {
		t.Fail()
		t.Log(err)
	}
}

//need actual environment
func TestDBHandler(t *testing.T) {
	eipDB, err := handler.NewMsgDB(handler.MsgDBConfig{
		DBConn: "sqlserver://sa:sasa@localhost?database=WebEIP5",
	})
	if err != nil {
		t.Fatal(err)
	}
	eip := &msg.EipMsgHandler{
		Control: eipDB,
	}
	if msgs, err := eip.GetAll(0, -1); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(msgs)
	}

	// if msg, err := eip.GetIndex(0); err != nil {
	// 	t.Log(err)
	// 	t.Fail()
	// } else {
	// 	t.Log(msg)
	// }

	// if err := eip.MarkRead(3); err != nil {
	// 	t.Log(err)
	// 	t.Fail()
	// }
}

func TestMSG_GetUnreadForAsync(t *testing.T) {
	rand.Seed(time.Now().Unix())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //chan is close then go func exit

	production := func(m *msg.EipMsgHandler, c int) <-chan *handler.EipMsg {
		return m.GetUnreadForAsync(ctx, c)
	}

	consumer := func(msgs <-chan *handler.EipMsg) int {
		go func() {
			<-time.After(time.Duration(rand.Intn(runtime.NumCPU())) * time.Millisecond)
			cancel()
		}()

		return *func(count *int) *int {
			for data := range msgs {
				if data != nil {
					*count++
					t.Log(data.UniqueID, data.Subject)
				} else {
					t.Error("read data is error")
				}
			}
			return count
		}(new(int))
	}

	//if max count is 1000
	var maxCount int = 1000
	if consumer(production(&msg.EipMsgHandler{}, maxCount)) < maxCount-1 {
		t.Log("data is cancel")
	} else {
		t.Log("data is full")
	}
}
