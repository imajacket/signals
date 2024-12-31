package signals

import "context"

// keyedListener represents a combination of a listener and an optional key used for identification.
type keyedListener[T any] struct {
	key      string
	listener ListenerExtended[T]
}

type ListenerExtended[T any] struct {
	handler  SignalListener[T]
	rollback SignalListener[T]
}

// BaseSignal provides the base implementation of the Signal interface.
// It is intended to be used as an abstract base for underlying signal mechanisms.
//
// Example:
//
//	type MyDerivedSignal[T any] struct {
//		BaseSignal[T]
//		// Additional fields or methods specific to MyDerivedSignal
//	}
//
//	func (s *MyDerivedSignal[T]) Emit(ctx context.Context, payload T) error {
//		// Custom implementation for emitting the signal
//	}
type BaseSignal[T any] struct {
	subscribers    []keyedListener[T]
	subscribersMap map[string]ListenerExtended[T]
}

// AddListener adds a handler and rollback (the listener) to the signal. The listener will be called
// whenever the signal is emitted. It returns the number of subscribers after
// the listener was added. It accepts an optional key that can be used to remove
// the listener later or to check if the listener was already added. It returns
// -1 if the listener with the same key was already added to the signal.
//
// Example:
//
//	signal := signals.New[int]()
//	count := signal.AddListener(func(ctx context.Context, payload int) error {
//		// Listener implementation
//		// ...
//	}, nil, "key1")
//	fmt.Println("Number of subscribers after adding listener:", count)
func (s *BaseSignal[T]) AddListener(handler, rollback SignalListener[T], key ...string) int {
	listener := ListenerExtended[T]{
		handler,
		rollback,
	}

	if len(key) > 0 {
		if _, ok := s.subscribersMap[key[0]]; ok {
			return -1
		}

		s.subscribersMap[key[0]] = listener
		s.subscribers = append(s.subscribers, keyedListener[T]{
			key:      key[0],
			listener: listener,
		})
	} else {
		s.subscribers = append(s.subscribers, keyedListener[T]{
			listener: listener,
		})
	}

	return len(s.subscribers)
}

// RemoveListener removes a listener from the signal. It returns the number
// of subscribers after the listener was removed. It returns -1 if the
// listener was not found.
//
// Example:
//
//	signal := signals.New[int]()
//	signal.AddListener(func(ctx context.Context, payload int) error {
//		// Listener implementation
//		// ...
//	}, nil, "key1")
//	count := signal.RemoveListener("key1")
//	fmt.Println("Number of subscribers after removing listener:", count)
func (s *BaseSignal[T]) RemoveListener(key string) int {
	if _, ok := s.subscribersMap[key]; ok {
		delete(s.subscribersMap, key)

		for i, sub := range s.subscribers {
			if sub.key == key {
				s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
				break
			}
		}
		return len(s.subscribers)
	}

	return -1
}

// Reset resets the signal by removing all subscribers from the signal,
// effectively clearing the list of subscribers.
// This can be used when you want to stop all listeners from receiving
// further signals.
//
// Example:
//
//	signal := signals.New[int]()
//	signal.AddListener(func(ctx context.Context, payload int) error {
//		// Listener implementation
//		// ...
//	}, nil)
//	signal.Reset() // Removes all listeners
//	fmt.Println("Number of subscribers after resetting:", signal.Len())
func (s *BaseSignal[T]) Reset() {
	s.subscribers = make([]keyedListener[T], 0)
	s.subscribersMap = make(map[string]ListenerExtended[T])
}

// Len returns the number of listeners subscribed to the signal.
// This can be used to check how many listeners are currently waiting for a signal.
// The returned value is of type int.
//
// Example:
//
//	signal := signals.New[int]()
//	signal.AddListener(func(ctx context.Context, payload int) {
//		// Listener implementation
//		// ...
//	})
//	fmt.Println("Number of subscribers:", signal.Len())
func (s *BaseSignal[T]) Len() int {
	return len(s.subscribers)
}

// IsEmpty checks if the signal has any subscribers.
// It returns true if the signal has no subscribers, and false otherwise.
// This can be used to check if there are any listeners before emitting a signal.
//
// Example:
//
//	signal := signals.New[int]()
//	fmt.Println("Is signal empty?", signal.IsEmpty()) // Should print true
//	signal.AddListener(func(ctx context.Context, payload int) {
//		// Listener implementation
//		// ...
//	})
//	fmt.Println("Is signal empty?", signal.IsEmpty()) // Should print false
func (s *BaseSignal[T]) IsEmpty() bool {
	return len(s.subscribers) == 0
}

// Emit is not implemented in BaseSignal and panics if called. It should be
// implemented by a derived type.
//
// Example:
//
//	type MyDerivedSignal[T any] struct {
//		BaseSignal[T]
//		// Additional fields or methods specific to MyDerivedSignal
//	}
//
//	func (s *MyDerivedSignal[T]) Emit(ctx context.Context, payload T) error {
//		// Custom implementation for emitting the signal
//	}
func (s *BaseSignal[T]) Emit(ctx context.Context, payload T) error {
	panic("implement me in derived type")
}
