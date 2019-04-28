package util

type ITask interface {
	DoTask() (error)
}

type ITaskQueue interface {
	AddTask(t ITask)
	Close()
}

func newAsyncSeqTaskQueue(chanCap int) ITaskQueue {
	queue := &asyncTaskQueue{ch: make(chan ITask, chanCap)}
	queue.StartRun()
	return queue
}

func newAsyncNoSeqTaskQueue(chanCap int) ITaskQueue {
	queue := &asyncNoSeqTaskQueue{&asyncTaskQueue{ch: make(chan ITask, chanCap)}}
	queue.StartRun()
	return queue
}

type asyncTaskQueue struct {
	ch chan ITask
}

func (tq *asyncTaskQueue) AddTask(t ITask) {
	tq.ch <- t
}

func (tq *asyncTaskQueue) StartRun() {
	go func() {
		for {
			task, ok := <-tq.ch
			if !ok {
				break
			}

			task.DoTask()
		}
	}()
}

func (tq *asyncTaskQueue) Close() {
	close(tq.ch)
}

type asyncNoSeqTaskQueue struct {
	*asyncTaskQueue
}

func (tq *asyncNoSeqTaskQueue) StartRun() {
	go func() {
		for {
			task, ok := <-tq.ch
			if !ok {
				break
			}

			go task.DoTask()
		}
	}()
}
