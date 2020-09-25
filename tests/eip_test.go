/*
 * @Author: your name
 * @Date: 2020-09-22 11:57:35
 * @LastEditTime: 2020-09-25 14:32:58
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\tests\msg_test.go
 */
package tests

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"test/cache"
	"test/msg"
	"test/msg/orm"
	"testing"
	"time"
	"unsafe"
)

func TestCache(t *testing.T) {
	c := cache.GetInstance()
	c2 := cache.GetInstance()
	if unsafe.Pointer(c) != unsafe.Pointer(c2) {
		log.Fatal("instance difference!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errfunc := func(value interface{}, err error) {
		if err != nil {
			log.Print(err)
		} else {
			log.Print(value)
		}
	}

	go func() {
		i := 0
		for {
			select {
			case <-ctx.Done():
			default:
				i++
				if err := c.Add(fmt.Sprintf("test_%d", i), cache.CacheItem{Data: fmt.Sprintf("value_%d", i), Expire: time.Duration(rand.Intn(3)) * time.Second}); err != nil {
					log.Fatal(err)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
			default:
				errfunc(c.Get(fmt.Sprintf("test_%d", rand.Intn(1000))))
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	time.Sleep(10 * time.Second)

}

func TestMockOrm(t *testing.T) {
	eip := &msg.Eip{
		Control: orm.NewOrmMock(),
	}
	if msgs, err := eip.GetAll(); err != nil {
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
func TestDBOrm(t *testing.T) {
	ormDB, err := orm.NewOrmDB(orm.OrmDBConfig{
		DBConn: "sqlserver://xx@xxx:9930?database=eip",
	})
	if err != nil {
		t.Fatal(err)
	}
	eip := &msg.Eip{
		Control: ormDB,
	}
	if msgs, err := eip.GetAll(); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(msgs)
	}

	if msg, err := eip.GetIndex(0); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		t.Log(msg)
	}

	if err := eip.MarkRead(3); err != nil {
		t.Log(err)
		t.Fail()
	}
}
