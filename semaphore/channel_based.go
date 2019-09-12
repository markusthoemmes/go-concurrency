package semaphore

import "context"

type semChan chan struct{}

func newChannelBased(n int64) Interface {
	return semChan(make(chan struct{}, n))
}

func (s semChan) Acquire(ctx context.Context) error {
	s <- struct{}{}
	return ctx.Err()
}

func (s semChan) TryAcquire() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s semChan) Release() {
	<-s
}
