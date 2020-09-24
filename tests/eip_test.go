/*
 * @Author: your name
 * @Date: 2020-09-22 11:57:35
 * @LastEditTime: 2020-09-24 21:13:38
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\tests\msg_test.go
 */
package tests

import (
	"test/msg"
	"test/msg/orm"
	"time"

	"testing"
)

func TestMockOrm(t *testing.T) {
	eip := &msg.Eip{
		Orm: orm.NewOrmMock(),
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
		Orm: ormDB,
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
