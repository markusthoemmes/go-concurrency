package gate

import (
	"context"
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
	for name, impl := range implementations {
		b.Run(name, func(b *testing.B) {
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

func BenchmarkGateOpenCloseFrequently(b *testing.B) {
	for name, impl := range implementations {
		b.Run(name, func(b *testing.B) {
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
