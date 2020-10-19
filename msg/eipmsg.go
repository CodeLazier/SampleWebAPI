package msg

import (
	"errors"
	"fmt"
	"log"
	"time"

	"test/cache"
	"test/handler"
	"test/queue"
)

//Eip impl
type EipMsgHandler struct {
	handler.Control
	useCache  bool
	cacheTime time.Duration
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

//NewEipDBHandler :fun helper wrap dbpool,new or get db conn in dbpool,automatically release after use
//but do not copy to external func to prevent the connection from to recycled
//call InitDB first if it has not been called
//提供自己的dbpool的根本原因是gorm或者其提供的db driver中的pool根本不靠谱(或许有某些神秘参数?),实测postgres/sqlite再大量并发模式下均有连接错误或写入失败问题
//另外一个原因是方便自己调度和控制.控制权再自己手里.
//并发测试见测试用例.1w并发下无写入和连接错误,再大也没问题.但量大效率会下降.需要db的横向扩展和多读单写策略,这里就暂不考虑了
func NewEipDBHandler(f func(*EipMsgHandler)) {
	dbctl := handler.GetMsgDB()
	if dbctl == nil {
		log.Println(fmt.Errorf("db conn is error"))
	} else {
		if f != nil {
			func() {
				defer func() {
					handler.PutMsgDB(dbctl)
				}()
				f(
					&EipMsgHandler{
						Control: dbctl,
					})
			}()
		}
	}
}

//UseCache if set use cache,then func use cache at call
func (t *EipMsgHandler) UseCache(useCache bool, cacheTime time.Duration, f func()) {
	t.useCache, t.cacheTime = useCache, cacheTime
	defer func() {
		t.useCache, t.cacheTime = false, 0
	}()
	if f != nil {
		f()
	}
}

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

func (t *EipMsgHandler) _getData(key cache_eipmsg, cmd handler.Cmd) (interface{}, error) {
	r := t.try_getCache(key)
	if r == nil {
		if r, err := t.Query(cmd); err != nil {
			return nil, err
		} else {
			return t.try_setCache(key, r), nil
		}
	}
	return r, nil
}

func (t *EipMsgHandler) try_getCache(key cache_eipmsg) interface{} {
	if t.useCache {
		if item, err := cache.GetInstance().Get(key); err != cache.ErrCacheNotFound {
			return item.Data
		}
	}
	return nil
}

func (t *EipMsgHandler) try_setCache(key cache_eipmsg, value interface{}) interface{} {
	if t.useCache {
		cache.GetInstance().Add(key, cache.NewCacheItem(value, t.cacheTime))
	}
	return value
}

//GetIndex get eipmsg for query id
func (t *EipMsgHandler) GetIndex(idx int) (interface{}, error) {
	return t._getData(cache_eipmsg{prefix: cache_prefix_GetIndex, start: idx, count: 1}, handler.NewRecord("Id = ?", idx))
}

//GetCount get all eipmsg items count
func (t *EipMsgHandler) GetCount() (int64, error) {
	cmd := handler.NewMultiRecords(0, -1)
	cmd.Order = ""
	cmd.CalcCount = true

	r, err := t._getData(cache_eipmsg{prefix: cache_prefix_GetCount, start: 0, count: -1}, cmd)
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
	r, err := t._getData(cache_eipmsg{prefix: cache_prefix_GetUnreadCount, start: 0, count: -1}, cmd)
	if r != nil {
		return r.(int64), nil
	}
	return -1, err
}

func (t *EipMsgHandler) GetUnread(start int, count int) (interface{}, error) {
	cmd := handler.NewMultiRecords(start, count)
	cmd.Query = "ReadTAG <> ? OR ReadTAG IS NULL"
	cmd.Args = []interface{}{ReadValue}

	return t._getData(cache_eipmsg{prefix: cache_prefix_GetUnread, start: start, count: count}, cmd)
}

func (t *EipMsgHandler) GetAll(start int, count int) (interface{}, error) {
	return t._getData(cache_eipmsg{prefix: cache_prefix_GetAll, start: start, count: count}, handler.NewMultiRecords(start, count))
}

func (t *EipMsgHandler) MarkRead(idx int) error {
	q := queue.GetInstance()
	q.SetDoFun(t._updateRead, false)
	q.Push(idx)
	//The front end always returns as correct
	//The backend needs to cooperate with the call link to be able to query result
	return nil
}

//New create a eipmsg write in db
func (t *EipMsgHandler) New(v interface{}) error {
	_, err := t.Insert(handler.NewInsertMsg(v))
	return err
}
