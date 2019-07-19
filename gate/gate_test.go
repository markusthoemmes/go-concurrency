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
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	implementations = map[string]func() Interface{
		"cond-based":           newCondBased,
		"channel-atomic-based": newChannelBasedWithAtomics,
		"channel-lock-based":   newChannelBasedWithLock,
		"lock-based":           newLockBased,
	}

	parallelisms = []int{1, 10, 100, 1000}

	bg = context.Background()
)

func TestGateCorrectness(t *testing.T) {
	for name, impl := range implementations {
		t.Run(name, func(t *testing.T) {
			gate := impl()
			calls := 100
			var grp sync.WaitGroup

			for i := 0; i < calls; i++ {
				grp.Add(1)
				go func() {
					gate.Wait(bg)
					grp.Done()
				}()
			}

			time.Sleep(10 * time.Millisecond)
			gate.Open()
			grp.Wait()
		})
	}
}

func BenchmarkGateOpen(b *testing.B) {
	for implName, impl := range implementations {
		for _, parallelism := range parallelisms {
			benchName := fmt.Sprintf("%s-%d", implName, parallelism)
			b.Run(benchName, func(b *testing.B) {
				b.SetParallelism(parallelism)
				g := impl()
				g.Open()

				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						g.Wait(bg)
					}
				})
			})
		}
	}
}

func BenchmarkGateOpenCloseFrequently(b *testing.B) {
	for implName, impl := range implementations {
		for _, parallelism := range parallelisms {
			benchName := fmt.Sprintf("%s-%d", implName, parallelism)
			b.Run(benchName, func(b *testing.B) {
				b.SetParallelism(parallelism)
				g := impl()
				go func() {
					for {
						time.Sleep(1 * time.Millisecond)
						g.Open()
						time.Sleep(1 * time.Millisecond)
						g.Close()
					}
				}()

				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						g.Wait(bg)
					}
				})
			})
		}
	}
}
