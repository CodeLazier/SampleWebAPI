/*
 * @Author: your name
 * @Date: 2020-09-24 14:20:01
 * @LastEditTime: 2020-09-28 14:13:36
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\msg\mock\ormmock.go
 */
package orm

import (
	"fmt"
	"test/msg"
)

type OrmMock struct {
}

func NewOrmMock() *OrmMock {
	return &OrmMock{}
}

//var mockData []msg.Msg

func (t *OrmMock) OpenOrm(...string) error {
	// for i := 0; i < 100; i++ {
	// 	var r bool = rand.Intn(2) == 0
	// 	mockData = append(mockData, msg.Msg{Id: 0, Title: fmt.Sprintf("Test_%d", i), Read: r})
	// }
	return nil
}

func (t *OrmMock) Select(q interface{}, a ...interface{}) (r interface{}, err error) {
	if q == nil {
		//mock some data
		r := make([]msg.Msg, 0)
		for i := 0; i < 10; i++ {
			r = append(r, msg.Msg{UniqueID: i, Subject: fmt.Sprintf("test_%d", i)})
		}
		return r, nil
	} else {
		//mock one data
		return &msg.Msg{UniqueID: 0, Subject: "test0"}, nil
	}
}

func (t *OrmMock) Update(idx int, field string, value interface{}) error {
	return nil
}
