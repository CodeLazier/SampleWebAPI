/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-30 12:41:34
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\eip.go
 */
package msg

import (
	"context"
	"errors"
	"fmt"
	"test/cache"
	"test/handler"
	"test/queue"
	"time"
)

//Eip impl
type EipMsgHandler struct {
	handler.Control
	UseCache bool
}

type cache_eipmsg struct {
	prefix string
	start  int
	count  int
	extra  interface{}
}

const (
	cache_prefix_GetAll         string = "GetAll()"
	cache_prefix_GetUnread      string = "GetUnread()"
	cache_prefix_GetIndex       string = "GetIndex()"
	cache_prefix_GetCount       string = "GetCount()"
	cache_prefix_GetUnreadCount string = "GetUnreadCount()"

	ReadValue int = 1
)

func (t *EipMsgHandler) _updateRead(v []interface{}) error {
	cmd := handler.NewUpdateRecord("ReadTAG", 0, "UniqueID in ?", v)
	return t.Update(cmd)
}

func conv2Msg(i interface{}) ([]handler.EipMsg, error) {
	switch msgs := i.(type) {
	case []handler.EipMsg:
		return msgs, nil
	case *[]handler.EipMsg:
		return *msgs, nil

	default:
		return nil, errors.New("return type is not incorrect")
	}
}

func (t *EipMsgHandler) _getData(key cache_eipmsg, cmd handler.Cmd, d time.Duration) (interface{}, error) {
	r := t.try_getCache(key)
	if r == nil {
		if r, err := t.Query(cmd); err != nil {
			return nil, err
		} else {
			return t.try_setCache(key, d, r), nil
		}
	}
	return r, nil
}

func (t *EipMsgHandler) try_getCache(key cache_eipmsg) interface{} {
	if t.UseCache {
		if item, err := cache.GetInstance().Get(key); err != cache.ErrCacheNotFound {
			return item.Data
		}
	}
	return nil
}

func (t *EipMsgHandler) try_setCache(key cache_eipmsg, exp time.Duration, value interface{}) interface{} {
	if t.UseCache {
		cache.GetInstance().Add(key, cache.NewCacheItem(value, exp))
	}
	return value
}

func (t *EipMsgHandler) GetIndex(idx int) (interface{}, error) {
	//note:cache is 0.5 hour,if content doesn’t change
	return t._getData(cache_eipmsg{prefix: cache_prefix_GetIndex, start: idx, count: 1}, handler.NewRecord("UniqueID = ?", idx), 30*time.Minute)
}

func (t *EipMsgHandler) GetCount() (int64, error) {
	cmd := handler.NewMultiRecords(0, -1)
	cmd.Order = ""
	cmd.CalcCount = true

	r, err := t._getData(cache_eipmsg{prefix: cache_prefix_GetCount, start: 0, count: -1}, cmd, 10*time.Second)
	if r != nil {
		return r.(int64), nil
	}
	return -1, err
}

func (t *EipMsgHandler) GetUnreadCount() (int64, error) {
	cmd := handler.NewMultiRecords(0, -1)
	cmd.Order = ""
	cmd.Query = "ReadTAG <> ? OR ReadTAG IS NULL"
	cmd.Args = []interface{}{ReadValue}
	cmd.CalcCount = true
	r, err := t._getData(cache_eipmsg{prefix: cache_prefix_GetUnreadCount, start: 0, count: -1}, cmd, 10*time.Second)
	if r != nil {
		return r.(int64), nil
	}
	return -1, err
}

func (t *EipMsgHandler) GetUnread(start int, count int) (interface{}, error) {
	cmd := handler.NewMultiRecords(start, count)
	cmd.Query = "ReadTAG <> ? OR ReadTAG IS NULL"
	cmd.Args = []interface{}{ReadValue}

	return t._getData(cache_eipmsg{prefix: cache_prefix_GetUnread, start: start, count: count}, cmd, 30*time.Second)
}

func (t *EipMsgHandler) GetAll(start int, count int) (interface{}, error) {
	return t._getData(cache_eipmsg{prefix: cache_prefix_GetAll, start: start, count: count}, handler.NewMultiRecords(start, count), 30*time.Second)
}

func (t *EipMsgHandler) MarkRead(idx int) error {
	q := queue.GetInstance()
	q.SetDoFun(t._updateRead, false)
	q.Push(idx)
	//The front end always returns as correct
	//The backend needs to cooperate with the call link to be able to query result
	return nil
}

//For testing only
func (t *EipMsgHandler) GetUnreadForAsync(ctx context.Context, maxCount int) <-chan *handler.EipMsg {
	data := make(chan *handler.EipMsg, 30) //buffer channel
	go func() {
		defer close(data)
		var err error
		if err != nil { // t.OpenDB("xx.xx.xx.xx"); err != nil {
			data <- nil
		} else {
			for i := 0; i < maxCount; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					data <- &handler.EipMsg{UniqueID: i, Subject: fmt.Sprintf("test_%d", i)}
				}
			}
		}
	}()
	return data
}
