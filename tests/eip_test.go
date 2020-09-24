/*
 * @Author: your name
 * @Date: 2020-09-22 11:57:35
 * @LastEditTime: 2020-09-24 15:27:08
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
	eip := msg.NewEip(msg.EipConfig{Orm: &orm.OrmMock{}})
	if msgs, err := eip.GetAll(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(msgs)
	}

	if msg, err := eip.GetIndex(0); err != nil {
		t.Fatal(err)
	} else {
		t.Log(msg)
	}

	if err := eip.MarkRead(3); err != nil {
		t.Fatal(err)
	}
}

//need actual environment
func TestDBOrm(t *testing.T) {
	eip := msg.NewEip(msg.EipConfig{Orm: &orm.OrmDB{}})
	if msgs, err := eip.GetAll(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(msgs)
	}

	if msg, err := eip.GetIndex(0); err != nil {
		t.Fatal(err)
	} else {
		t.Log(msg)
	}

	if err := eip.MarkRead(3); err != nil {
		t.Fatal(err)
	}
}
