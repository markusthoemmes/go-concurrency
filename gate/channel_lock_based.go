/*
Copyright 2019 Markus Thömmes

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gate

import (
	"context"
	"sync"
)

// newChannelBasedWithLock returns a gate that is based on a broadcast channel and
// an atomically updated open/closed state, used as a fast-path.
// Once the gate closes a channel is generated which concurrent Wait calls
// read from. Once the gate opens again, that channel is closed, notifying
// all waiters to move on.
func newChannelBasedWithLock() Interface {
	return &channelBasedWithLock{
		broadcast: make(chan struct{}),
	}
}

type channelBasedWithLock struct {
	broadcast chan struct{}
	open      uint32

	// mux guards Open and Close from being called concurrently and the state
	// mutations in general.
	mux sync.RWMutex
}

func (c *channelBasedWithLock) Open() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.open == 0 {
		c.open = 1
		close(c.broadcast)
	}
}

func (c *channelBasedWithLock) Close() {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.open == 1 {
		c.broadcast = make(chan struct{})
		c.open = 0
	}
}

func (c *channelBasedWithLock) Wait(ctx context.Context) error {
	c.mux.RLock()
	open := c.open
	c.mux.RUnlock()
	if open > 0 {
		return ctx.Err()
	}
	select {
	case <-c.broadcast:
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}
