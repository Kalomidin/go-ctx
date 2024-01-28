package context

import (
	"sync"
	"sync/atomic"
	"time"
)

type cancelCtx struct {
	mu sync.Mutex // protects following fields

	done     atomic.Value
	err      error
	children map[canceler]struct{}
}

func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
	}()

	done := c.done.Load()
	if done != nil {
		return done.(chan struct{})
	} else {
		done = make(chan struct{})
		c.done.Store(done)
		return done.(chan struct{})
	}
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.err
	return err
}

func (c *cancelCtx) Value(key any) any {
	return nil
}

func (c *cancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *cancelCtx) Cancel(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return // already canceled
	}

	for child := range c.children {
		child.Cancel(err)
		<-child.Done()
		delete(c.children, child)
	}
	c.err = err
	done := c.done.Load()
	if done != nil {
		close(done.(chan struct{}))
	} else {
		c.done.Store(make(chan struct{}))
		close(c.done.Load().(chan struct{}))
	}
}

func (c *cancelCtx) Add(child canceler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.children == nil {
		c.children = make(map[canceler]struct{})
	}
	c.children[child] = struct{}{}
}

func (c *cancelCtx) Delete(child canceler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.children, child)
}

func (c *cancelCtx) Get() map[canceler]struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.children
}
