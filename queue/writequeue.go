package queue

import (
	"sync"
	"time"
)

type WriteActionsHandler func([]interface{}) error

var _queue *WriteQueue

const maxCount int = 3000

func GetInstance() *WriteQueue {
	if _queue == nil {
		m := sync.Mutex{}
		m.Lock()
		if _queue == nil {
			//lazy load
			_queue = new(WriteQueue)
			_queue.q = make([]interface{}, 0, 3000)
			go _queue.doAction()
		}
		m.Unlock()
	}
	return _queue
}

type WriteQueue struct {
	q          []interface{}
	m          sync.RWMutex
	actionsFun WriteActionsHandler
}

func (w *WriteQueue) SetDoFun(handler WriteActionsHandler, force bool) {
	if w.actionsFun != nil {
		if force {
			w.m.Lock()
			w.actionsFun = handler
			w.m.Unlock()
		}
	} else {
		w.actionsFun = handler
	}
}

func (w *WriteQueue) Push(value ...interface{}) {
	w.m.Lock()
	defer w.m.Unlock()
	w.q = append(w.q, value...)
}

func (w *WriteQueue) doAction() {
	if w.actionsFun != nil {
		timer := time.NewTimer(time.Second)
		for {
			<-timer.C
			func() {
				defer timer.Reset(time.Second)
				w.m.RLock()
				l := len(w.q)
				if l > 0 {
					if l > maxCount {
						l = maxCount
					}
					w.actionsFun(w.q[:l])
				}
				w.m.RUnlock()
				//clear
				w.m.Lock()
				w.q = w.q[l:]
				w.m.Unlock()
			}()
		}
	}
}
