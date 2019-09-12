package semaphore

import "context"

type Interface interface {
	// Acquire gets a token.
	Acquire(context.Context) error
	// TryAcquire returns true if it can get a token immediately.
	TryAcquire() bool
	// Release returns a token.
	Release()
}
