/*
 * @Author: your name
 * @Date: 2020-09-22 11:20:05
 * @LastEditTime: 2020-09-22 14:32:43
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\db\eip.go
 */
package msg

import "errors"

type Eip struct{}

func (t *Eip) GetUnread() ([]Msg, error) {
	return nil, errors.New("not implement")
}

func (t *Eip) GetIndex(idx int) (*Msg, error) {
	return &Msg{Id: 0, Title: "Test Title"}, nil
}
