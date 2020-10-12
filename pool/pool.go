package pool

import (
	"log"
	"runtime"
	"sync/atomic"
)

var MAX_POOL_SIZE int32 = int32(runtime.NumCPU()*2 + 1)

//限制池,包装一个缓存通道,今后可以添加更多控制
//DB链接时珍贵的有限资源,此限制池严格限制,超出链接获取将等待直到获得(可以有timeout机制)
type LimitPool struct {
	pool chan interface{}
	c    int32
	n    func() interface{}

	IN  chan interface{}
	OUT chan interface{}
}

func (l *LimitPool) Deinc() {
	atomic.AddInt32(&l.c, -1)
}

func (l *LimitPool) inc() {
	atomic.AddInt32(&l.c, 1)
}

func (l *LimitPool) put() {
	for {
		if int32(len(l.pool)) == MAX_POOL_SIZE {
			log.Print("Pool is full,will discard the first")
			<-l.pool //discard
		}
		l.pool <- <-l.IN
		//fmt.Println("Put", data)
	}
}

func (l *LimitPool) get() {
	for {
		n := len(l.pool)
		if n == 0 && l.c < MAX_POOL_SIZE {
			l.inc()
			l.pool <- l.n()
		}
		//timeout need maybe
		l.OUT <- <-l.pool
		//fmt.Println("Get", data)
	}

}

func New(n func() interface{}) *LimitPool {
	l := LimitPool{}
	l.pool = make(chan interface{}, MAX_POOL_SIZE)
	l.IN = make(chan interface{})
	l.OUT = make(chan interface{})
	l.n = n
	go l.put()
	go l.get()
	return &l
}
