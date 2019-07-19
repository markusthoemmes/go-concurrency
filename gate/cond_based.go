package gate

import (
	"context"
	"sync"
)

func NewCondBased() Interface {
	var mux sync.Mutex
	return &condBased{
		cond: sync.NewCond(&mux),
	}
}

type condBased struct {
	cond *sync.Cond
	open bool
}

func (c *condBased) Open() {
	c.cond.L.Lock()
	c.open = true
	c.cond.L.Unlock()
	c.cond.Broadcast()
}

func (c *condBased) Close() {
	c.cond.L.Lock()
	c.open = false
	c.cond.L.Unlock()
}

func (c *condBased) Wait(ctx context.Context) error {
	c.cond.L.Lock()
	for !c.open {
		c.cond.Wait()
	}
	c.cond.L.Unlock()
	return ctx.Err()
}
