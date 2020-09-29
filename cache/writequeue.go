/*
 * @Author: your name
 * @Date: 2020-09-29 21:14:24
 * @LastEditTime: 2020-09-29 22:23:16
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\cache\writequeue.go
 */
package cache

import (
	"container/list"
	"sync"
	"time"
)

type WriteActionsHandler func([]interface{}) error

var _queue *WriteQueue

func GetWriteQueueInstance() *WriteQueue {
	if _queue == nil {
		m := sync.Mutex{}
		m.Lock()
		if _queue == nil {
			//lazy load
			_queue = new(WriteQueue)
			_queue.l = new(list.List)
			go _queue.doAction()
		}
		m.Unlock()
	}
	return _queue
}

type WriteQueue struct {
	l          *list.List
	m          sync.RWMutex
	ActionsFun WriteActionsHandler
}

func (w *WriteQueue) Push(value ...interface{}) {
	w.m.Lock()
	defer w.m.Unlock()
	for _, v := range value {
		w.l.PushBack(v)
	}
}

func (w *WriteQueue) doAction() {
	if w.ActionsFun != nil {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			w.m.RLock()
			v := make([]interface{}, 0)
			for e := w.l.Front(); e != nil; e = e.Next() {
				v = append(v, e.Value)
			}
			if len(v) > 0 {
				w.ActionsFun(v)
			}
			w.m.RUnlock()
			//clear
			w.m.Lock()
			w.l = w.l.Init()
			w.m.Unlock()
		}
	}
}
