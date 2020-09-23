/*
 * @Author: your name
 * @Date: 2020-09-22 11:57:35
 * @LastEditTime: 2020-09-23 11:17:41
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\tests\msg_test.go
 */
package tests

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

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

func check_do_id(m *msgmock.MockMessages, t *testing.T, idx int) (*msg.Msg, error) {
	if msg, err := m.GetIndex(idx); err != nil {
		t.Error(err)
		return nil, err
	} else {
		return msg, err
	}

}

func check_do_unread(m *msgmock.MockMessages, t *testing.T) {
	if msg, err := m.GetUnread(); err != nil {
		t.Error(err)
	} else {
		t.Log(msg)
	}
}

func TestMSG_GetIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := msgmock.NewMockMessages(ctrl)

	m.EXPECT().GetIndex(gomock.Any()).DoAndReturn(check_idx)
	m.EXPECT().GetIndex(gomock.Any()).DoAndReturn(check_idx)

	check_do_id(m, t, 0)
	check_do_id(m, t, -10)

	for i := 0; i < 10; i++ {
		msg := &msg.Eip{}
		if ms, err := msg.GetIndex(i); err != nil {
			t.Error(err)
		} else if ms.Id != i {
			t.Errorf("id is difference for:%d", i)
		}
	}
}

//add flag -v and clear test cache (set flag -count=1)
//production consumer
func TestMSG_GetUnreadForAsync(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := &msg.Eip{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() //chan is close then go func exit
	retChan := m.GetUnreadForAsync(ctx)
	go func() {
		select {
		//set n mill then will timeout
		case <-time.After(time.Duration(rand.Intn(8)) * time.Millisecond):
			cancel()
		}

	}()

	c := *func(count *int32) *int32 {
		for data := range retChan {
			if data != nil {
				atomic.AddInt32(count, 1) // No need for this env
				t.Log(data.Id, data.Title)
			} else {
				t.Error("read data is error")
			}
		}
		return count
	}(new(int32))

	//if max count is 1000
	if c < 1000-1 {
		t.Log("data is cancel")
	} else {
		t.Log("data is full")
	}

}

func TestMSG_GetUnread(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := msgmock.NewMockMessages(ctrl)

	m.EXPECT().GetUnread().DoAndReturn(check_unread)

	check_do_unread(m, t)

}

func TestMSG_GetIndex2(t *testing.T) {
	eip := &msg.Eip{}
	eip.GetIndex(2)

}
