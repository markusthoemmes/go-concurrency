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

// newCondBased returns a gate that is based sync.Cond.
func newCondBased() Interface {
	var mux sync.Mutex
	return &condBased{
		cond: sync.NewCond(&mux),
	}
}

type condBased struct {
	cond *sync.Cond
	open bool
}

// Open implements Interface.
func (c *condBased) Open() {
	c.cond.L.Lock()
	c.open = true
	c.cond.L.Unlock()
	c.cond.Broadcast()
}

// Close implements Interface.
func (c *condBased) Close() {
	c.cond.L.Lock()
	c.open = false
	c.cond.L.Unlock()
}

// Wait implements Interface.
// Note that this does not implement aborting of Wait if the context is done.
func (c *condBased) Wait(ctx context.Context) error {
	c.cond.L.Lock()
	for !c.open {
		c.cond.Wait()
	}
	c.cond.L.Unlock()
	return ctx.Err()
}
