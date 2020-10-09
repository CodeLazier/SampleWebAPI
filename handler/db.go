package handler

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"gorm.io/driver/postgres"
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

var ERR_NO_AFFECTED error = fmt.Errorf("No affected rows")

var _one sync.Once
var _msgdb *MsgDB

var _DB_DNS string
var _DB_DEBUG_MODE bool

func InitDB(dns string, debug bool) {
	if _DB_DNS != "" {
		//DB dns will only take valid the first time you set it up
	} else {
		_DB_DNS = dns
		_DB_DEBUG_MODE = debug
	}
}

func GetInstance() (*MsgDB, error) {
	var err error
	_one.Do(func() {
		cf := func() error {
			if _msgdb, err = NewMsgDB(MsgDBConfig{
				DBConn: _DB_DNS,        // "user=postgres password=sasa dbname=postgres port=5432",
				Debug:  _DB_DEBUG_MODE, //true is output raw sql
			}); err != nil {
				return err
			}
			return nil
		}

		checkf := func() {
			ticker := time.NewTicker(10 * time.Second)
			for {
				<-ticker.C
				if _msgdb != nil {
					if err := _msgdb.Ping(); err != nil {
						cf()
					}
				}
			}
		}

		if cf() == nil {
			go checkf()
		}
	})
	return _msgdb, err
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
		statement = statement.Limit(cmd.Count)
		statement = statement.Offset(cmd.Start)

		if cmd.Query != nil {
			statement = statement.Where(cmd.Query, cmd.Args...)
		}

		return statement
	}
	return t.db
}

func (t *MsgDB) _recover() error {
	var err error
	var r interface{}
	if r = recover(); r != nil {
		err = r.(error)
		fmt.Println(r)
	}
	return err
}

func (t *MsgDB) Ping() error {
	if t.db != nil {
		if db, err := t.db.DB(); err != nil {
			return err
		} else {
			return db.Ping()
		}
	}
	return fmt.Errorf("db is null")
}

func (t *MsgDB) OpenOrm(cfg ...string) error {
	t.Lock()
	defer t.Unlock()
	if len(cfg) > 0 {
		dsn := cfg[0]
		var err error
		if t.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
			return err
		} else {
			db, err := t.db.DB()
			if err != nil {
				return err
			}
			if t.Cfg.Debug {
				t.db = t.db.Debug()
			}
			db.SetMaxIdleConns(runtime.NumCPU() * 2)
			//default inifinte
			//db.SetConnMaxLifetime(time.Hour)
			return t.Ping()
		}
	}
	return fmt.Errorf("config is null")
}

func (t *MsgDB) Insert(cmd Cmd) (interface{}, error) {
	if t.db != nil {
		t.Lock()
		defer t.Unlock()
		defer t._recover()
		if msg, ok := cmd.Update.Value.(EipMsg); ok {
			tx := t.db.Create(&msg)
			return msg, tx.Error
		} else {
			return nil, fmt.Errorf("interface is not compatible")
		}
	}
	return nil, errors.New("Open db first")
}

func (t *MsgDB) Query(cmd Cmd) (result interface{}, err error) {
	if t.db != nil {
		t.RLock()
		defer t.RUnlock()

		var tx *gorm.DB
		defer t._recover()

		tx = t.buildSql(cmd)
		if cmd.CalcCount {
			var c int64 = -1
			tx = tx.Count(&c)
			result = c
		} else {
			tx = tx.Find(cmd.Model)
			if tx.RowsAffected > 0 {
				result = cmd.Model
			} else {
				err = errors.New("the query didn't return any data")
			}
		}
		if err == nil {
			err = tx.Error
		}
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
		defer t._recover()

		if cmd.Update.Value != nil {
			var tx *gorm.DB
			if cmd.Update.Field == "^All" {
				tx = t.db.Model(cmd.Model).Where(cmd.Query, cmd.Args...).Updates(cmd.Model)
			}
			tx = t.db.Model(cmd.Model).Where(cmd.Query, cmd.Args...).Update(cmd.Update.Field, cmd.Update.Value)

			if tx.Error != nil {
				return tx.Error
			}
			if tx.RowsAffected == 0 {
				return ERR_NO_AFFECTED
			}
			return nil
		}
	}
	return errors.New("Open db first")
}
