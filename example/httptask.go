package main

import (
	"fmt"
	"github.com/huizluo/threadpool/task"
	"time"
)

type HttpTaskFactory struct{}

func (ht *HttpTaskFactory) CreateTask() task.Task {
	return &HttpTask{}
}

type HttpTask struct {
	Id  int
}

func (h *HttpTask) SetID(id int) {
	h.Id = id
}

func (ht *HttpTask) Run() {
	fmt.Printf("task [ID = %d] download file \n", ht.Id)
	time.Sleep(time.Millisecond * 5)
}

