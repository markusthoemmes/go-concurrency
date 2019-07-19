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

// newLockBased returns a gate that is based on sync.RWMutex.
// The gate is open when the lock is unlocked, close if it locked.
// Concurrent waits are possible via only acquiring read locks.
func newLockBased() Interface {
	gate := &lockBased{}
	gate.mux.Lock()
	return gate
}

type lockBased struct {
	mux          sync.RWMutex
	openCloseMux sync.Mutex
	open         bool
}

// Open implements interface.
func (c *lockBased) Open() {
	c.openCloseMux.Lock()
	defer c.openCloseMux.Unlock()
	if !c.open {
		c.open = true
		c.mux.Unlock()
	}
}

// Close implements interface.
func (c *lockBased) Close() {
	c.openCloseMux.Lock()
	defer c.openCloseMux.Unlock()
	if c.open {
		c.open = false
		c.mux.Lock()
	}
}

// Wait implements Interface.
// Note that this does not implement aborting of Wait if the context is done.
func (c *lockBased) Wait(ctx context.Context) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return ctx.Err()
}
