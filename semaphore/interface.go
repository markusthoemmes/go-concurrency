package semaphore

import "context"

type Interface interface {
	// Acquire gets a token.
	Acquire(context.Context) error
	TryAcquire() bool
	// Release returns a token.
	Release()
}
