//+build DBTest

package tests

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"test/handler"
	"test/msg"
	v1 "test/v1"
)

func init() {
	handler.InitDB("host=localhost user=postgres password=sasa dbname=postgres port=5432", true)
}

// * warning:will modify the database
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
		if msgs, err := eip.GetIndex(0); err != nil {
			t.Fail()
			t.Log(err)
		} else {
			t.Log(msgs)
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

func NewDBConn() (*handler.MsgDB, error) {
	if dbctrl, err := handler.NewMsgDB(
		handler.MsgDBConfig{
			DBConn: "host=localhost user=postgres password=sasa dbname=postgres port=5432", // "user=postgres password=sasa dbname=postgres port=5432",
			Debug:  false,                                                                  //true is output raw sql
		}); err != nil {
		return nil, err
	} else {
		return dbctrl, nil
	}
}

func Benchmark_DBHandler(b *testing.B) {
	msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
		for i := 0; i < b.N; i++ {
			eip.GetAll(0, -1)
			//eip.GetUnread(0, -1)
			eip.GetCount()
			//eip.GetUnreadCount()
		}
	})
}

func TestDBPostDBHandler_concurrency(t *testing.T) {
	wg := sync.WaitGroup{}
	count := 5000 //many concurrency?
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			v1.Test_PostNewEipMsg()
		}()
	}
	wg.Wait()
}

func TestDBInsertHandler_concurrency(t *testing.T) {
	wg := sync.WaitGroup{}
	count := 5000 //many concurrency?
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
				eip.New(handler.EipMsg{
					Title:   fmt.Sprintf("test%d", i),
					Content: fmt.Sprintf("content is %d", i),
				})
			})
		}(i)
	}
	wg.Wait()
}

func TestCreateTable(t *testing.T) {
	if db, err := NewDBConn(); err != nil {
		log.Fatalln(err)
	} else {
		if !db.RawDB.Migrator().HasTable(&handler.EipMsg{}) {
			if err = db.RawDB.AutoMigrate(&handler.EipMsg{}); err != nil {
				log.Fatalln(err)
			}
		}
	}

}

//set flag -v -count=1
func TestMSG_GetUnreadForAsync(t *testing.T) {
	rand.Seed(time.Now().Unix())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //chan is close then go func exit

	production := func(c int) <-chan *handler.EipMsg {
		return func() <-chan *handler.EipMsg {
			data := make(chan *handler.EipMsg, 30) //buffer channel
			go func() {
				defer close(data)
				var err error
				if err != nil { // t.OpenDB("xx.xx.xx.xx"); err != nil {
					data <- nil
				} else {
					for i := 0; i < c; i++ {
						select {
						case <-ctx.Done():
							return
						default:
							data <- &handler.EipMsg{Id: i, Title: fmt.Sprintf("test_%d", i)}
						}
					}
				}
			}()
			return data
		}()
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
					t.Log(data.Id, data.Title)
				} else {
					t.Error("read data is error")
				}
			}
			return count
		}(new(int))
	}

	//if max count is 1000
	var maxCount = 1000
	if consumer(production(maxCount)) < maxCount-1 {
		t.Log("data is cancel")
	} else {
		t.Log("data is full")
	}
}
