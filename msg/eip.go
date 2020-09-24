/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-24 15:14:10
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

type EipConfig struct {
	Orm OrmMsg
	//
}

//Eip impl
type Eip struct {
	cfg EipConfig
}

func NewEip(cfg EipConfig) *Eip {
	return &Eip{cfg: cfg}
}

func conv2Msg(i interface{}) ([]Msg, error) {
	if msgs, ok := i.([]Msg); ok {
		return msgs, nil
	}
	return nil, errors.New("return type is not incorrect")
}

func (t *Eip) GetUnread() ([]Msg, error) {
	if r, err := t.cfg.Orm.Select("read = ?", false); err != nil {
		return []Msg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *Eip) GetIndex(idx int) (*Msg, error) {
	r, err := t.cfg.Orm.Select("uniqueID = ?", idx)
	if err != nil {
		return nil, err
	}
	if msg, ok := r.(*Msg); ok {
		return msg, nil
	}
	return nil, fmt.Errorf("query result multiple count")

}

func (t *Eip) GetAll() ([]Msg, error) {
	if r, err := t.cfg.Orm.Select(nil); err != nil {
		return []Msg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *Eip) MarkRead(idx int) error {
	return t.cfg.Orm.Update(idx, "Read", true)
}

//For testing only
func (t *Eip) GetUnreadForAsync(ctx context.Context, maxCount int) <-chan *Msg {
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
