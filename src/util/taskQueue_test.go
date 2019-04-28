package util

import (
	"log"
	"testing"
)

func NewTask(v int) ITask {
	return &testTask{v}
}

type testTask struct {
	v int
}

func (t *testTask) DoTask() (error) {
	log.Println("value:", t.v, ", goRouting:")
	return nil
}

func TestTaskQueue(t *testing.T) {
	len := 1000
	queue := NewSeqWaitModel(len)
	defer queue.Close()

	for i := 0; i < len; i++ {
		queue.AddTask(NewTask(i))
	}
	queue.Wait()
}
