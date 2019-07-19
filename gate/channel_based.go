package gate

import (
	"context"
	"sync"
	"sync/atomic"
)

func NewChannelBased() Interface {
	return &channelBased{
		broadcast: make(chan struct{}),
	}
}

type channelBased struct {
	broadcast chan struct{}
	open      int32

	// mux guards Open and Close from being called concurrently.
	mux sync.Mutex
}

func (c *channelBased) Open() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if atomic.CompareAndSwapInt32(&c.open, 0, 1) {
		close(c.broadcast)
	}
}

func (c *channelBased) Close() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if atomic.CompareAndSwapInt32(&c.open, 1, 0) {
		c.broadcast = make(chan struct{})
	}
}

func (c *channelBased) Wait(ctx context.Context) error {
	if atomic.LoadInt32(&c.open) > 0 {
		return ctx.Err()
	}
	select {
	case <-c.broadcast:
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}
