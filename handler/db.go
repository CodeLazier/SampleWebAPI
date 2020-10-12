package handler

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"test/pool"
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

const WAIT_DBCONN_TIMEOUT = 6 * time.Second

var ERR_NO_AFFECTED = fmt.Errorf("no affected rows")

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

//When it comes to parameter settings, the closure mode cannot be used
var _one = sync.Once{}
var _instancePool *pool.LimitPool

func createPool() *pool.LimitPool {
	_one.Do(func() {
		_instancePool = pool.New(func() interface{} {
			if msgDB, err := NewMsgDB(MsgDBConfig{
				DBConn: _DB_DNS,
				Debug:  _DB_DEBUG_MODE,
			}); err != nil {
				log.Println(err)
				return nil
			} else {
				return msgDB
			}
		})

	})
	return _instancePool
}

func PutMsgDB(m *MsgDB) {
	createPool().IN <- m
}

func GetMsgDB() *MsgDB {
	//singleton instance mode
	dbpool := createPool()
	a := time.NewTimer(WAIT_DBCONN_TIMEOUT)
	defer a.Stop()
	//timeout
	select {
	case result := <-dbpool.OUT:
		if result != nil {
			//TODO 每次获取均检查心跳,可能再大量并发获取的状况下损害效能.考虑添加时间判断,超过一定时间才会判断
			if err := result.(*MsgDB).Ping(); err != nil {
				log.Println("db conn is deaded")
				//retry
				dbpool.Deinc()
				return GetMsgDB()
			} else {
				return result.(*MsgDB)
			}
		}

	case <-a.C:
		log.Println("db is busy,please wait or retry")
	}
	return nil
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
			db.SetMaxIdleConns(runtime.NumCPU()*2 + 1)
			//default inifinte
			db.SetMaxOpenConns(runtime.NumCPU()*2 + 1)
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
	return nil, errors.New("open db first")
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
	return nil, errors.New("open db first")
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
	return errors.New("open db first")
}
