/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-28 22:18:09
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\eip.go
 */
package msg

import (
	"context"
	"errors"
	"fmt"

	"test/handler"
	. "test/handler"
)

type Messages interface {
	//get msgs for read,sort by senddate
	GetUnread(start int, count int) ([]EipMsg, error)
	//get msg for uniqueid
	GetIndex(id int) (*EipMsg, error)
	//get all msgs,sort by senddate
	GetAll(start int, count int) ([]EipMsg, error)
	//set read
	MarkRead(idx int) error

	GetUnradCount() (int, error)
	GetCount() (int, error)

	//...
}

//Eip impl
type EipMsgHandler struct {
	handler.Control
}

func conv2Msg(i interface{}) ([]EipMsg, error) {
	switch msgs := i.(type) {
	case []EipMsg:
		return msgs, nil
	case *[]EipMsg:
		return *msgs, nil

	default:
		return nil, errors.New("return type is not incorrect")
	}
}

func (t *EipMsgHandler) GetUnread(start int, count int) ([]EipMsg, error) {
	if r, err := t.Query(handler.DefaultCmd); err != nil {
		return []EipMsg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsgHandler) GetIndex(idx int) (*EipMsg, error) {
	r, err := t.Query(handler.DefaultCmd)
	if err != nil {
		return nil, err
	}
	if msg, ok := r.(*EipMsg); ok {
		return msg, nil
	}
	return nil, fmt.Errorf("query result multiple count")

}

func (t *EipMsgHandler) GetCount() ([]EipMsg, error) {
	if r, err := t.Query(handler.DefaultCmd); err != nil {
		return []EipMsg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsgHandler) GetUnreadCount() ([]EipMsg, error) {
	if r, err := t.Query(handler.DefaultCmd); err != nil {
		return []EipMsg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsgHandler) GetAll(start int, count int) ([]EipMsg, error) {
	var msgs []EipMsg
	if r, err := t.Query(handler.Cmd{Model: &msgs, Start: start, Count: count}); err != nil {
		return []EipMsg{}, err
	} else {
		return conv2Msg(r)
	}
}

func (t *EipMsgHandler) MarkRead(idx int) error {
	return t.Update(handler.DefaultCmd)
}

//For testing only
func (t *EipMsgHandler) GetUnreadForAsync(ctx context.Context, maxCount int) <-chan *EipMsg {
	data := make(chan *EipMsg, 30) //buffer channel
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
					data <- &EipMsg{UniqueID: i, Subject: fmt.Sprintf("test_%d", i)}
				}
			}
		}
	}()
	return data
}
