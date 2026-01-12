package patterns

import (
	"context"
	"sync"
)

// Constructor is a function type that creates and initializes an instance of type T.
// It receives a context that can be used for initialization and returns a pointer to T.
type Constructor[T any] func(ctx context.Context) *T

// Singleton provides a thread-safe generic implementation of the singleton pattern.
// It ensures that only one instance of type T is created, even when accessed
// concurrently from multiple goroutines.
//
// The instance is created lazily on the first call to Instance() using the
// provided constructor function. The sync.Once mechanism guarantees that the
// constructor is called exactly once.
//
// Example usage:
//
//	type Database struct {
//	    conn *sql.DB
//	}
//
//	dbConstructor := func(ctx context.Context) *Database {
//	    return &Database{conn: initConnection()}
//	}
//
//	dbSingleton := NewSingleton(dbConstructor).WithContext(context.Background())
//	db := dbSingleton.Instance() // First call creates the instance
//	db2 := dbSingleton.Instance() // Returns the same instance
type Singleton[T any] struct {
	once        sync.Once
	constructor Constructor[T]
	ctx         context.Context
	instance    *T
}

// NewSingleton creates a new Singleton manager for type T with the given constructor function.
// The constructor will be called lazily when Instance() is first invoked.
//
// Parameters:
//   - constructor: A function that creates and initializes an instance of T
//
// Returns:
//   - A pointer to a Singleton[T] manager
//
// Note: Use WithContext to set the context before calling Instance() if your
// constructor requires a context for initialization.
func NewSingleton[T any](constructor Constructor[T]) *Singleton[T] {
	return &Singleton[T]{
		constructor: constructor,
	}
}

// WithContext sets the context that will be passed to the constructor function.
// This method should be called before the first call to Instance().
//
// Parameters:
//   - ctx: The context to use when initializing the singleton instance
//
// Returns:
//   - The Singleton instance (for method chaining)
//
// Note: Calling this method after Instance() has been called has no effect,
// as the instance has already been created.
func (singleton *Singleton[T]) WithContext(ctx context.Context) *Singleton[T] {
	if singleton.ctx == nil {
		singleton.ctx = ctx
	}
	return singleton
}

// Instance returns the singleton instance of type T.
// On the first call, it creates the instance using the constructor function.
// Subsequent calls return the same instance without calling the constructor again.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// The sync.Once mechanism ensures the constructor is called exactly once.
//
// Returns:
//   - A pointer to the singleton instance of type T
func (singleton *Singleton[T]) Instance() *T {
	singleton.once.Do(func() {
		if singleton.ctx == nil {
			singleton.ctx = context.Background()
		}
		singleton.instance = singleton.constructor(singleton.ctx)
	})
	return singleton.instance
}
