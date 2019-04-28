package util

type ITask interface {
	DoTask() (error)
}

type ITaskQueue interface {
	AddTask(t ITask)
	Close()
}

func newAsyncTaskQueue(chanCap int) ITaskQueue {
	queue := &asyncTaskQueue{ch: make(chan ITask, chanCap)}
	queue.StartRun()
	return queue
}

type asyncTaskQueue struct {
	ch chan ITask
}

func (tq *asyncTaskQueue) AddTask(t ITask) {
	tq.ch <- t
}

func (tq *asyncTaskQueue) Close() {
	close(tq.ch)
}

func (tq *asyncTaskQueue) StartRun() {
	go func() {
		for {
			task, ok := <-tq.ch
			if !ok {
				break
			}

			_ = task.DoTask()
		}
	}()
}
