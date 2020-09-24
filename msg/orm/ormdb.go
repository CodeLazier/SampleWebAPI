/*
 * @Author: your name
 * @Date: 2020-09-24 14:20:01
 * @LastEditTime: 2020-09-24 15:35:42
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\msg\mock\ormmock.go
 */
package orm

import (
	"errors"
	"fmt"
	"sync"
	"test/msg"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type OrmDB struct {
	db *gorm.DB
	sync.RWMutex
}

func (t *OrmDB) OpenOrm(cfg ...string) error {
	t.Lock()
	defer t.Unlock()
	dsn := cfg[0] // "sqlserver://xx@xxx:9930?database=eip"
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

func (t *OrmDB) Select(query interface{}, args ...interface{}) (result interface{}, err error) {
	if t.db != nil {
		t.RLock()
		defer t.RUnlock()
		var tx *gorm.DB
		defer func() {
			var r interface{}
			if r = recover(); r != nil {
				err = r.(error)
				fmt.Println(r)
			}
		}()
		msgs := make([]msg.Msg, 0)
		if query != nil {
			tx = t.db.Where(query, args).Find(&msgs)
		} else {
			tx = t.db.Find(&msgs)
		}
		result = msgs
		err = tx.Error
		if err != nil {
			return result, err
		} else {
			return result, nil
		}
	}
	return nil, errors.New("Open db first")
}

func (t *OrmDB) Update(idx int, field string, value interface{}) (err error) {
	if t.db != nil {
		t.Lock()
		defer t.Unlock()
		defer func() {
			var r interface{}
			if r = recover(); r != nil {
				err = r.(error)
				fmt.Println(r)
			}
		}()

		tx := t.db.Model(&msg.Msg{}).Where("Id = ?", idx).Update(field, value)
		return tx.Error
	}
	return errors.New("Open db first")
}
