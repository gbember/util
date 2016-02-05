// gtimer.go
package util

import (
	"container/heap"
	"sync"
	"time"
)

const (
	max_d time.Duration = 1 << 32
)

type GTimer struct {
	mut    sync.Mutex
	tmr    *time.Timer
	ref    int64
	tfl    timer_fun_list
	funBuf []func()
	min    int64
	ctr    chan bool
}

func NewGTimer() *GTimer {
	gtmr := new(GTimer)
	gtmr.tmr = time.NewTimer(max_d)
	gtmr.min = time.Now().Add(max_d).UnixNano()
	gtmr.tfl = timer_fun_list(make([]*timer_fun, 0, 100))
	gtmr.funBuf = make([]func(), 0, 1)
	gtmr.ctr = make(chan bool)
	go gtmr.loop()
	return gtmr
}

type timer_fun_list []*timer_fun

type timer_fun struct {
	t   int64
	ref int64
	fun func()
}

//加一个延迟执行任务函数
func (gtmr *GTimer) AddAfter(t time.Duration, fun func()) (ref int64) {
	gtmr.mut.Lock()
	defer gtmr.mut.Unlock()
	if gtmr.tmr == nil {
		panic("GTimer is closeds")
	}

	gtmr.ref++

	tf := &timer_fun{
		t:   time.Now().Add(t).UnixNano(),
		ref: gtmr.ref,
		fun: fun,
	}

	heap.Push(&gtmr.tfl, tf)

	if tf.t < gtmr.min {
		gtmr.min = tf.t
		gtmr.tmr.Reset(t)
	}

	return gtmr.ref
}

//取消一个延迟任务函数
func (gtmr *GTimer) Cancel(ref int64) bool {
	gtmr.mut.Lock()
	defer gtmr.mut.Unlock()
	if gtmr.tmr == nil {
		panic("GTimer is closeds")
	}
	index := -1
	for i := 0; i < len(gtmr.tfl); i++ {
		if gtmr.tfl[i].ref == ref {
			index = i
			break
		}
	}
	if index != -1 {
		heap.Remove(&gtmr.tfl, index)
		return true
	}
	return false
}

func (gtmr *GTimer) Stop() {
	gtmr.mut.Lock()
	defer gtmr.mut.Unlock()
	if gtmr != nil {
		gtmr.tmr.Stop()
		close(gtmr.ctr)
		gtmr.tmr = nil
		gtmr.tfl = nil
	}
}

func (gtmr *GTimer) loop() {
	for gtmr.tmr != nil {
		select {
		case <-gtmr.tmr.C:
			gtmr.run()
		case <-gtmr.ctr:
			return
		}
	}
}
func (gtmr *GTimer) run() {
	gtmr.mut.Lock()
	defer gtmr.mut.Unlock()
	if gtmr == nil {
		return
	}
	now := time.Now().UnixNano()

	for gtmr.tfl.Len() > 0 {
		tf := heap.Pop(&gtmr.tfl).(*timer_fun)
		if tf.t > now {
			heap.Push(&gtmr.tfl, tf)
			gtmr.min = tf.t
			gtmr.tmr.Reset(time.Duration(tf.t - now))
			break
		}
		//TODO 处理panic
		gtmr.funBuf = append(gtmr.funBuf, tf.fun)
	}
	if gtmr.tfl.Len() == 0 {
		gtmr.min = time.Now().Add(max_d).UnixNano()
		gtmr.tmr.Reset(max_d)
	}
	for _, fun := range gtmr.funBuf {
		fun()
	}
	gtmr.funBuf = gtmr.funBuf[0:0]
}

//
func (tfl timer_fun_list) Len() int           { return len(tfl) }
func (tfl timer_fun_list) Less(i, j int) bool { return tfl[i].t < tfl[j].t }
func (tfl timer_fun_list) Swap(i, j int)      { tfl[i], tfl[j] = tfl[j], tfl[i] }
func (tfl *timer_fun_list) Push(x interface{}) {
	*tfl = append(*tfl, x.(*timer_fun))
}
func (tfl *timer_fun_list) Pop() interface{} {
	old := *tfl
	length := len(old)
	x := old[length-1]
	*tfl = old[:length-1]
	return x
}

func (tfl *timer_fun_list) getMin() *timer_fun {
	old := *tfl
	length := len(old)
	if length > 0 {
		return old[length-1]
	}
	return nil
}
