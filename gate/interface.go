package gate

import "context"

type Interface interface {
	// Open opens the gate.
	Open()
	// Close closes the gate.
	Close()
	// Wait blocks while the gate is closed.
	Wait(context.Context) error
}

func New() Interface {
	return newChannelBasedWithAtomics()
}
