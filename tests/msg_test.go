/*
 * @Author: your name
 * @Date: 2020-09-22 11:57:35
 * @LastEditTime: 2020-09-22 14:34:41
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\tests\msg_test.go
 */
package tests

import (
	"errors"
	"strconv"
	"testing"

	"test/msg"
	msgmock "test/msg/mocks"

	"github.com/golang/mock/gomock"
)

func check_idx(idx int) (*msg.Msg, error) {
	if idx < 0 {
		return nil, errors.New("id not found")
	}
	return &msg.Msg{Id: idx, Title: "test title"}, nil

}

func check_unread() ([]msg.Msg, error) {
	r := make([]msg.Msg, 0)
	for i := 0; i < 3; i++ {
		r = append(r, msg.Msg{Id: i, Title: "test title" + strconv.Itoa(i)})
	}

	return r, nil

}

func check_do_id(m *msgmock.MockMessages, t *testing.T, idx int) {
	if msg, err := m.GetIndex(idx); err != nil {
		t.Error(err)
	} else {
		t.Log(msg)
	}
}

func check_do_unread(m *msgmock.MockMessages, t *testing.T) {
	if msg, err := m.GetUnread(); err != nil {
		t.Error(err)
	} else {
		t.Log(msg)
	}
}

func Test_MSG_GetIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := msgmock.NewMockMessages(ctrl)

	m.EXPECT().GetIndex(gomock.Any()).DoAndReturn(check_idx)
	m.EXPECT().GetIndex(gomock.Any()).DoAndReturn(check_idx)

	check_do_id(m, t, 0)
	check_do_id(m, t, -10)

}

func Test_MSG_GetUnread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := msgmock.NewMockMessages(ctrl)

	m.EXPECT().GetUnread().DoAndReturn(check_unread)

	check_do_unread(m, t)

}
