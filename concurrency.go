package structures

import (
	"context"
	"sync"
)

type ConcurrencyHandler struct {
	curTaskIdx    int
	taskCh        chan ConcurrencyTask
	concurrencyCh chan struct{}
	wg            *sync.WaitGroup
}

func NewConcurrencyHandler(maxConcurrentTasks int) *ConcurrencyHandler {
	wg := sync.WaitGroup{}
	wg.Add(1)
	return &ConcurrencyHandler{
		curTaskIdx:    0,
		taskCh:        make(chan ConcurrencyTask),
		concurrencyCh: make(chan struct{}, maxConcurrentTasks),
		wg:            &wg,
	}
}

func (c *ConcurrencyHandler) Done() {
	c.wg.Done()
}

func (c *ConcurrencyHandler) Wait() {
	c.wg.Wait()
}

func (c *ConcurrencyHandler) Close() {
	close(c.taskCh)
	close(c.concurrencyCh)
}

type ConcurrencyTask struct {
	f      func()
	idx    int
	ctx    context.Context
	cancel context.CancelFunc
}

func (c ConcurrencyTask) Ctx() context.Context {
	return c.ctx
}

func (c ConcurrencyTask) Idx() int {
	return c.idx
}

func (c *ConcurrencyTask) Cancel() {
	c.cancel()
}

func (c *ConcurrencyHandler) Start() {
	go func() {
		for task := range c.taskCh {
			c.concurrencyCh <- struct{}{}
			go func(task ConcurrencyTask) {
				select {
				case <-task.ctx.Done():
					break
				default:
					task.f()
				}
				c.Done()
				<-c.concurrencyCh
			}(task)
		}
	}()
}

func (c *ConcurrencyHandler) Enqueue(task func()) ConcurrencyTask {
	c.curTaskIdx++
	c.wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	cTask := ConcurrencyTask{
		f:      task,
		idx:    c.curTaskIdx,
		ctx:    ctx,
		cancel: cancel,
	}
	c.taskCh <- cTask
	return cTask
}
