package util

import (
	"log"
	"testing"
	"time"
)

func NewTask(v int) ITask {
	return &testTask{v}
}

type testTask struct {
	v int
}

func (t *testTask) DoTask() (error) {
	log.Println("value:", t.v)
	return nil
}

func TestSeqWaitModel(t *testing.T) {
	len := 100
	queue := NewSeqWaitModel(len)
	defer queue.Close()

	for i := 0; i < len; i++ {
		queue.AddTask(NewTask(i))
	}
	queue.Wait()
}

func TestSeqNoWaitModel(t *testing.T) {
	len := 100
	queue := NewSeqNoWaitModel(len)
	defer queue.Close()

	for i := 0; i < len; i++ {
		queue.AddTask(NewTask(i))
	}

	time.Sleep(time.Second * 1)
}

func TestNoSeqWaitModel(t *testing.T) {
	len := 100
	queue := NewNoSeqWaitModel(len)
	defer queue.Close()

	for i := 0; i < len; i++ {
		queue.AddTask(NewTask(i))
	}
	queue.Wait()
}

func TestNoSeqNoWaitModel(t *testing.T) {
	len := 100
	queue := NewNoSeqNoWaitModel(len)
	defer queue.Close()

	for i := 0; i < len; i++ {
		queue.AddTask(NewTask(i))
	}

	time.Sleep(time.Second * 1)
}
