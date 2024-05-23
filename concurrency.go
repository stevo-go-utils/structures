package structures

import "context"

type ConcurrencyHandler struct {
	curTaskIdx    int
	taskCh        chan ConcurrencyTask
	concurrencyCh chan struct{}
}

func NewConcurrencyHandler(maxConcurrentTasks int) *ConcurrencyHandler {
	return &ConcurrencyHandler{
		curTaskIdx:    0,
		taskCh:        make(chan ConcurrencyTask),
		concurrencyCh: make(chan struct{}, maxConcurrentTasks),
	}
}

func (c *ConcurrencyHandler) Close() {
	close(c.taskCh)
	close(c.concurrencyCh)
}

type ConcurrencyTask struct {
	f   func()
	idx int
	ctx context.Context
}

func (c ConcurrencyTask) Ctx() context.Context {
	return c.ctx
}

func (c ConcurrencyTask) Idx() int {
	return c.idx
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
				<-c.concurrencyCh
			}(task)
		}
	}()
}

func (c *ConcurrencyHandler) Enqueue(task func()) context.CancelFunc {
	c.curTaskIdx++
	ctx, cancel := context.WithCancel(context.Background())
	c.taskCh <- ConcurrencyTask{
		f:   task,
		idx: c.curTaskIdx,
		ctx: ctx,
	}
	return cancel
}
