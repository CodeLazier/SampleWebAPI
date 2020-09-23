/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-23 13:43:41
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\eip.go
 */
package msg

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

//Eip impl
type Eip struct {
	db *gorm.DB
	sync.RWMutex
}

func getData(db *gorm.DB, query interface{}, args ...interface{}) (msgs []Msg, err error) {
	var tx *gorm.DB
	defer func() {
		var r interface{}
		if r = recover(); r != nil {
			err = r.(error)
			fmt.Println(r)
		}
	}()
	if query != nil {
		tx = db.Where(query, args).Find(&msgs)
	} else {
		tx = db.Find(&msgs)
	}
	err = tx.Error
	if err != nil {
		return msgs, err
	} else {
		return msgs, nil
	}
}

func (t *Eip) GetUnread() ([]Msg, error) {
	t.RLock()
	defer t.RUnlock()
	return getData(t.db, "read = ?", false)
}

func (t *Eip) GetIndex(idx int) (*Msg, error) {
	t.RLock()
	defer t.RUnlock()
	msgs, err := getData(t.db, nil)
	if err != nil {
		return nil, err
	}
	if len(msgs) > 0 {
		return &msgs[0], nil
	} else {
		return nil, errors.New("idx is not found msg")
	}
}

func (t *Eip) GetAll() ([]Msg, error) {
	t.RLock()
	defer t.RUnlock()
	msgs := make([]Msg, 100)
	r := t.db.Find(&msgs)
	if r.Error != nil {
		return []Msg{}, r.Error
	} else {
		return msgs, nil
	}
}

func (t *Eip) MarkRead(idx int) error {
	t.Lock()
	defer t.Unlock()
	return t.db.Model(&Msg{}).Where("Id = ?", idx).Update("Read", true).Error
}

func (t *Eip) OpenDB(cfg string) error {
	dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
	var err error
	if t.db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{}); err != nil {
		return err
	} else {
		db, err := t.db.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	}
}

func (t *Eip) GetUnreadForAsync(ctx context.Context, maxCount int) <-chan *Msg {
	data := make(chan *Msg, 30) //buffer channel
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
					data <- &Msg{Id: i, Title: fmt.Sprintf("test_%d", i)}
				}
			}
		}
	}()
	return data
}
