package pool

import (
	"errors"
	"fmt"
	"github.com/huizluo/threadpool/task"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrPoolClose = errors.New("pool close error")
)

type Pool struct {
	max     int32
	running int32
	free chan struct{}
	threads []*Thread
	lock    sync.Mutex
	once    sync.Once
}

func (p *Pool) Init() {
	p.once.Do(func() {
		p.max = 20
		p.running = 0
		p.free = make(chan struct{})
		p.threads = make([]*Thread, 0)
	})
}

//add task to pool
func (p *Pool) Submit(t task.Task) error {
	if len(p.free) > 0 {
		return ErrPoolClose
	}
	w := p.getThread()
	w.addTask(t)
	w.run()
	return nil
}

func (p *Pool) Release() {
	p.threads = nil
	p.free <- struct{}{}
}

//get worker from pool
func (p *Pool) getThread() *Thread {
	var t *Thread
RE_GET:
	p.lock.Lock()
	n := len(p.threads)
	fmt.Printf("pool [thread = %d] \n", len(p.threads))
	if n == 0 {
		if p.running == p.max {
			//TODO reset pool max size
			fmt.Println("wait for other thread done")
			time.Sleep(time.Millisecond * 10)
			p.lock.Unlock()
			goto RE_GET
		} else {
			t = &Thread{
				Pool: p,
				task: make(chan task.Task),
			}
			p.running++
			t.run()
		}
	} else {
		t = p.threads[n-1]
		p.threads[n-1] = nil
		p.threads = p.threads[:n-1]
	}
	p.lock.Unlock()

	return t
}

//
func (p *Pool) putThread(thread *Thread) {
	p.lock.Lock()
	p.threads = append(p.threads, thread)
	p.lock.Unlock()
}

//调整pool大小
func (p *Pool) reSize(size int32) {
	if size < p.max {
		diff := p.max - size
		var i int32
		for ; i < diff; i++ {
			//删除多余的thread
			p.getThread().stop()
		}
	} else if size == p.max {
		return
	}
	atomic.StoreInt32(&p.max, size)
}

//Thread
type Thread struct {
	Pool *Pool
	task chan task.Task
}

//stop 关闭通道
func (thread *Thread) stop() {
	close(thread.task)
}

func (thread *Thread) addTask(t task.Task) {
	thread.task <- t
}

func (thread *Thread) run() {
	go func() {
		for t := range thread.task {
			if t == nil {
				atomic.AddInt32(&thread.Pool.running, -1)
				break
			}
			fmt.Printf("pool [running = %d] \n", thread.Pool.running)
			t.Run()
			//这里归还了thread 但是当前这个协程还在running，从task chan读取阻赛，直到这个通道被关闭
			thread.Pool.putThread(thread)
		}
	}()
}
