package context

import (
	"time"
)

type timerCtx struct {
	*cancelCtx
	timer    *time.Timer
	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return "context.timerCtx"
}

func (c *timerCtx) Value(key any) any {
	return nil
}

func (c *timerCtx) Done() <-chan struct{} {
	return c.cancelCtx.Done()
}

func (c *timerCtx) Err() error {
	return c.cancelCtx.Err()
}

func (c *timerCtx) Cancel(err error) {
	c.cancelCtx.Cancel(err)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
}

func (c *timerCtx) Add(child canceler) {
	c.cancelCtx.Add(child)
}

func (c *timerCtx) Delete(child canceler) {
	c.cancelCtx.Delete(child)
}

func (c *timerCtx) Get() map[canceler]struct{} {
	return c.cancelCtx.Get()
}
