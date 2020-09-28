/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-28 15:47:54
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\eip.go
 */
package msg

import (
	"context"
	"errors"
	"fmt"
)

//Eip impl
type EipMsg struct {
	Control
}

func (m Msg) TableName() string {
	return "EIP_MessageMaster"
}

func conv2Msg(i interface{}) ([]Msg, error) {
	if msgs, ok := i.([]Msg); ok {
		return msgs, nil
	}
	return nil, errors.New("return type is not incorrect")
}

func (t *EipMsg) GetUnread(where CustomWhere, start int, count int) ([]Msg, error) {
	if r, err := t.Select("read = ?", false); err != nil {
		return []Msg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsg) GetIndex(idx int) (*Msg, error) {
	r, err := t.Select("UniqueID = ?", idx)
	if err != nil {
		return nil, err
	}
	if msg, ok := r.(*Msg); ok {
		return msg, nil
	}
	return nil, fmt.Errorf("query result multiple count")

}

func (t *EipMsg) GetCount(where CustomWhere) ([]Msg, error) {
	if r, err := t.Select(where.Query, where.Args...); err != nil {
		return []Msg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsg) GetUnreadCount(where CustomWhere) ([]Msg, error) {
	if r, err := t.Select(where.Query, where.Args...); err != nil {
		return []Msg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsg) GetAll(where CustomWhere, start int, count int) ([]Msg, error) {
	args := make([]interface{}, 0)
	args = append(args, start, count)
	args = append(args, where.Args...)

	if r, err := t.Select(where.Query, args...); err != nil {
		return []Msg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsg) MarkRead(idx int) error {
	return t.Update(idx, "Read", true)
}

//For testing only
func (t *EipMsg) GetUnreadForAsync(ctx context.Context, maxCount int) <-chan *Msg {
	data := make(chan *Msg, 30) //buffer channel
	go func() {
		defer close(data)
		var err error
		if err != nil { // t.OpenDB("xx.xx.xx.xx"); err != nil {
			data <- nil
		} else {
			for i := 0; i < maxCount; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					data <- &Msg{UniqueID: i, Subject: fmt.Sprintf("test_%d", i)}
				}
			}
		}
	}()
	return data
}
