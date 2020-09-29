/*
 * @Author: your name
 * @Date: 2020-09-24 14:20:01
 * @LastEditTime: 2020-09-29 12:01:15
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\msg\mock\ormmock.go
 */
package handler

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MsgDB struct {
	db *gorm.DB
	sync.RWMutex
	Cfg MsgDBConfig
}

type MsgDBConfig struct {
	DBConn string
	Debug  bool
	//
}

func NewMsgDB(cfg MsgDBConfig) (*MsgDB, error) {
	r := &MsgDB{Cfg: cfg}
	if err := r.OpenOrm(cfg.DBConn); err != nil {
		return nil, err
	}
	return r, nil
}

func (t *MsgDB) buildSql(cmd Cmd) *gorm.DB {
	if t.db != nil {
		//optimization depends on the order
		statement := t.db
		if cmd.CalcCount {
			statement = statement.Model(cmd.Model)
		}
		if cmd.Order != "" {
			statement = statement.Order(cmd.Order)
		}
		if cmd.Start > 0 {
			statement = statement.Limit(cmd.Start)
		}
		if cmd.Count > 0 {
			statement = statement.Offset(cmd.Count)
		}
		if cmd.Query != nil {
			statement = statement.Where(cmd.Query, cmd.Args...)
		}

		return statement
	}
	return t.db
}

func (t *MsgDB) OpenOrm(cfg ...string) error {
	t.Lock()
	defer t.Unlock()
	if len(cfg) > 0 {
		dsn := cfg[0]
		var err error
		if t.db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{}); err != nil {
			return err
		} else {
			db, err := t.db.DB()
			if err != nil {
				return err
			}
			if t.Cfg.Debug {
				t.db = t.db.Debug()
			}
			return db.Ping()
		}
	}
	return fmt.Errorf("config is null")
}

func (t *MsgDB) Query(cmd Cmd) (result interface{}, err error) {
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

		tx = t.buildSql(cmd)
		if cmd.CalcCount {
			var c int64 = -1
			tx = tx.Count(&c)
			result = c
		} else {
			tx = tx.Find(cmd.Model)
			result = cmd.Model
		}
		err = tx.Error
		if err != nil {
			return result, err
		} else {
			return result, nil
		}
	}
	return nil, errors.New("Open db first")
}

func (t *MsgDB) Update(cmd Cmd) (err error) {
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

		//tx := t.db.Model(&msg.Msg{}).Where("Id = ?", ).Update(field, value)
		//return tx.Error
		return nil
	}
	return errors.New("Open db first")
}
