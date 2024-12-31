package signals

import (
	"context"
	"errors"
)

// SyncSignal is a struct that implements the Signal interface.
// It provides a synchronous way of notifying all subscribers of a signal.
// The type parameter `T` is a placeholder for any type.
type SyncSignal[T any] struct {
	BaseSignal[T]
}

// Emit notifies all subscribers of the signal and passes the payload in a
// synchronous way.
//
// The payload is of the same type as the SyncSignal's type parameter `T`.
// The method iterates over the subscribers slice of the SyncSignal,
// and for each subscriber, it calls the subscriber's listener function,
// passing the context and the payload.
// If the context has a deadline or cancellable property, the listeners
// must respect it. This means that the listeners should stop processing when
// the context is cancelled. Unlike the AsyncSignal's Emit method, this method
// does not call the listeners in separate goroutines, so the listeners are
// called synchronously, one after the other.
//
// Example:
//
//	signal := signals.NewSync[string]()
//	signal.AddListener(func(ctx context.Context, payload string) {
//		// Listener implementation
//		// ...
//	})
//
//	signal.Emit(context.Background(), "Hello, world!")
func (s *SyncSignal[T]) Emit(ctx context.Context, payload T) error {
	count := -1
	for index, sub := range s.subscribers {
		err := sub.listener.handler(ctx, payload)
		if err != nil {
			count = index
			break
		}
	}

	if count == -1 {
		return nil
	}

	for index, sub := range s.subscribers {
		_ = sub.listener.rollback(ctx, payload)

		if index == count {
			break
		}
	}

	return errors.New("there was an error emitting signal, rollback was called")
}
