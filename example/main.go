package main

import (
	"github.com/huizluo/threadpool/pool"
	"github.com/huizluo/threadpool/task"
	"time"
)

func main() {
	var factory task.TaskFactory
	var p pool.Pool
	p = pool.Pool{}
	p.Init()
	factory = &HttpTaskFactory{}
	for i := 1; i < 1000; i++ {
		t := factory.CreateTask()
		t.SetID(i)
		p.Submit(t)
	}
	time.Sleep(time.Second)
}
