/*
 * @Author: your name
 * @Date: 2020-09-24 14:20:01
 * @LastEditTime: 2020-09-28 16:01:39
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
	Cfg OrmDBConfig
}

type OrmDBConfig struct {
	DBConn string
	//
}

func NewOrmDB(cfg OrmDBConfig) (*OrmDB, error) {
	r := &OrmDB{Cfg: cfg}
	if err := r.OpenOrm(cfg.DBConn); err != nil {
		return nil, err
	}
	return r, nil
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
		start := 0
		count := -1
		cargs := make([]interface{}, 0)
		cargs = append(cargs, args...)

		if len(cargs) > 1 {
			start = cargs[0].(int)
			count = cargs[1].(int)

			cargs = cargs[2:]

		}

		var msgs []msg.Msg
		if query != nil {
			tx = t.db.Order("SendDate desc").Limit(count).Offset(start).Where(query, cargs).Find(&msgs)
		} else {
			tx = t.db.Order("SendDate desc").Limit(count).Offset(start).Find(&msgs)
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
