package util

// 异步不等待模型
type IAsyncNoWaitModel interface {
	//
	ITaskQueue
}

// 异步模型
type IAsyncWaitModel interface {
	//
	IAsyncNoWaitModel

	// 等待, 根据语义表示阻塞还是不阻塞
	Wait()
}

func NewSeqNoWaitModel(chanCap int) IAsyncNoWaitModel {
	return &seqNoWaitModel{newAsyncTaskQueue(chanCap)}
}

func (sm *seqNoWaitModel) AddTask(t ITask) {
	sm.ITaskQueue.AddTask(t)
}

func (sm *seqNoWaitModel) Close() {
	sm.ITaskQueue.Close()
}

// 异步不等待顺序模型
type seqNoWaitModel struct {
	ITaskQueue
}

func NewSeqWaitModel(chanCap int) IAsyncWaitModel {
	return &seqWaitModel{newAsyncTaskQueue(chanCap), make(chan bool, chanCap), chanCap}
}

func newWaitTask(t ITask, ch chan bool) ITask {
	return &waitTask{t, ch}
}

type waitTask struct {
	ITask
	ch chan bool
}

func (t *waitTask) DoTask() (error) {
	err := t.ITask.DoTask()

	f := t.CallBack()
	defer f()

	return err
}

func (wt *waitTask) CallBack() (func()) {
	return func() {
		wt.ch <- true
	}
}

func (sm *seqWaitModel) AddTask(t ITask) {
	sm.ITaskQueue.AddTask(newWaitTask(t, sm.ch))
}

func (sm *seqWaitModel) Close() {
	sm.ITaskQueue.Close()
}

func (sm *seqWaitModel) Wait() {
	for i := 0; i < sm.chLen; i++ {
		<-sm.ch
	}
}

// 异步等待顺序模型
type seqWaitModel struct {
	ITaskQueue
	ch    chan bool
	chLen int
}
