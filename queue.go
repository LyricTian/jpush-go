package jpush

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// Jober 任务执行接口
type Jober interface {
	Call()
}

// NewQueue 创建队列实例
func NewQueue(maxThread int) *Queue {
	return &Queue{
		maxWorker:  maxThread,
		workers:    make([]*queueWorker, maxThread),
		workerPool: make(chan chan Jober, maxThread),
		list:       list.New(),
		lock:       new(sync.RWMutex),
		wg:         new(sync.WaitGroup),
	}
}

// Queue 任务队列
type Queue struct {
	maxWorker  int
	workers    []*queueWorker
	workerPool chan chan Jober
	list       *list.List
	lock       *sync.RWMutex
	wg         *sync.WaitGroup
	running    uint32
}

// Start 开始执行
func (q *Queue) Start() {
	if atomic.LoadUint32(&q.running) == 1 {
		return
	}
	atomic.StoreUint32(&q.running, 1)

	// 启动工作线程
	for i := 0; i < q.maxWorker; i++ {
		q.workers[i] = newQueueWorker(q.workerPool, q.wg)
		q.workers[i].Start()
	}

	go q.dispatcher()
}

func (q *Queue) dispatcher() {
	for {
		q.lock.RLock()
		if atomic.LoadUint32(&q.running) != 1 && q.list.Len() == 0 {
			q.lock.RUnlock()
			break
		}
		ele := q.list.Front()
		q.lock.RUnlock()

		if ele == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		worker := <-q.workerPool
		worker <- ele.Value.(Jober)

		q.lock.Lock()
		q.list.Remove(ele)
		q.lock.Unlock()
	}
}

// Push 放入队列，等待执行
func (q *Queue) Push(param Jober) {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	q.wg.Add(1)
	q.lock.Lock()
	q.list.PushBack(param)
	q.lock.Unlock()
}

// Stop 停止
func (q *Queue) Stop() {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	atomic.StoreUint32(&q.running, 0)
	q.wg.Wait()

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}

func newQueueWorker(pool chan chan Jober, wg *sync.WaitGroup) *queueWorker {
	return &queueWorker{
		pool:    pool,
		wg:      wg,
		jobChan: make(chan Jober),
		quit:    make(chan struct{}),
	}
}

// 工作线程
type queueWorker struct {
	pool    chan chan Jober
	wg      *sync.WaitGroup
	jobChan chan Jober
	quit    chan struct{}
}

// 开启工作线程
func (w *queueWorker) Start() {
	w.pool <- w.jobChan

	go func() {
		for {
			select {
			case j := <-w.jobChan:
				j.Call()
				close(w.jobChan)

				w.jobChan = make(chan Jober)
				w.pool <- w.jobChan
				w.wg.Done()
			case <-w.quit:
				<-w.pool
				close(w.jobChan)
				return
			}
		}
	}()
}

// 停止工作线程
func (w *queueWorker) Stop() {
	close(w.quit)
}
