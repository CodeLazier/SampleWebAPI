/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-22 22:32:33
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

type Eip struct{}

func (t *Eip) GetUnread() ([]Msg, error) {
	return nil, errors.New("not implement")
}

func (t *Eip) GetIndex(idx int) (*Msg, error) {
	return &Msg{Id: idx, Title: "Test Title"}, nil
}

func (t *Eip) OpenDB(ip string) error {
	//do somting...
	return nil
}

func (t *Eip) GetUnreadForAsync(ctx context.Context) <-chan *Msg {
	data := make(chan *Msg, 30) //buffer channel
	go func() {
		defer close(data)
		if err := t.OpenDB("xx.xx.xx.xx"); err != nil {
			data <- nil
		} else {
			for i := 0; i < 1000; i++ {
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
