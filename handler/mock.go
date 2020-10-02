/*
 * @Author: your name
 * @Date: 2020-09-24 14:20:01
 * @LastEditTime: 2020-09-30 21:23:51
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\msg\mock\ormmock.go
 */
package handler

import (
	"fmt"
)

type CtlMock struct {
}

func NewCtlMock() *CtlMock {
	return &CtlMock{}
}

//var mockData []msg.Msg

func (t *CtlMock) OpenOrm(...string) error {
	// for i := 0; i < 100; i++ {
	// 	var r bool = rand.Intn(2) == 0
	// 	mockData = append(mockData, msg.Msg{Id: 0, Title: fmt.Sprintf("Test_%d", i), Read: r})
	// }
	return nil
}

func (t *CtlMock) Insert(cmd Cmd) (interface{}, error) {
	return nil, nil
}

func (t *CtlMock) Query(cmd Cmd) (r interface{}, err error) {
	if cmd.Query == nil {
		//mock some data
		r := make([]EipMsg, 0)
		for i := 0; i < 10; i++ {
			r = append(r, EipMsg{Id: i, Title: fmt.Sprintf("test_%d", i)})
		}
		return r, nil
	} else {
		//mock one data
		return &EipMsg{Id: 0, Title: "test0"}, nil
	}
}

func (t *CtlMock) Update(cmd Cmd) error {
	return nil
}
