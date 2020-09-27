/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-25 21:33:21
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

func conv2Msg(i interface{}) ([]Msg, error) {
	if msgs, ok := i.([]Msg); ok {
		return msgs, nil
	}
	return nil, errors.New("return type is not incorrect")
}

func (t *EipMsg) GetUnread() ([]Msg, error) {
	if r, err := t.Select("read = ?", false); err != nil {
		return []Msg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsg) GetIndex(idx int) (*Msg, error) {
	r, err := t.Select("uniqueID = ?", idx)
	if err != nil {
		return nil, err
	}
	if msg, ok := r.(*Msg); ok {
		return msg, nil
	}
	return nil, fmt.Errorf("query result multiple count")

}

func (t *EipMsg) GetAll() ([]Msg, error) {
	if r, err := t.Select(nil); err != nil {
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
					data <- &Msg{Id: i, Title: fmt.Sprintf("test_%d", i)}
				}
			}
		}
	}()
	return data
}
