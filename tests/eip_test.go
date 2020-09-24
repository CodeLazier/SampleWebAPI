/*
 * @Author: your name
 * @Date: 2020-09-22 11:57:35
 * @LastEditTime: 2020-09-24 15:49:48
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\tests\msg_test.go
 */
package tests

import (
	"test/msg"
	"test/msg/orm"

	"testing"
)

func TestMockOrm(t *testing.T) {
	eip := msg.NewEip(msg.EipConfig{Orm: orm.NewOrmMock()})
	if msgs, err := eip.GetAll(); err != nil {
		t.Fail()
		t.Log(err)
	} else {
		t.Log(msgs)
	}

	if msg, err := eip.GetIndex(0); err != nil {
		t.Fail()
		t.Log(err)
	} else {
		t.Log(msg)
	}

	if err := eip.MarkRead(3); err != nil {
		t.Fail()
		t.Log(err)
	}
}

//need actual environment
func TestDBOrm(t *testing.T) {
	ormDB, err := orm.NewOrmDB(orm.OrmDBConfig{DBConn: "sqlserver://xx@xxx:9930?database=eip"})
	if err != nil {
		t.Fatal(err)
	}
	eip := msg.NewEip(msg.EipConfig{Orm: ormDB})
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
