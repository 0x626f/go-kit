package patterns

import (
	"sync"
)

// Constructor is a function type that creates and initializes an instance of type Type.
// It receives a context that can be used for initialization and returns a pointer to Type.
type Constructor[Type any] func() (*Type, error)

// Singleton provides a thread-safe generic implementation of the singleton pattern.
// It ensures that only one instance of type Type is created, even when accessed
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
//	dbConstructor := func() (*Database, error) {
//	    conn, err := initConnection()
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &Database{conn: conn}, nil
//	}
//
//	dbSingleton := NewSingleton(dbConstructor)
//	db := dbSingleton.Instance() // First call creates the instance
//	if dbSingleton.Err() != nil {
//	    // Handle initialization error
//	}
//	db2 := dbSingleton.Instance() // Returns the same instance
type Singleton[Type any] struct {
	once        sync.Once
	constructor Constructor[Type]
	instance    *Type
	err         error
}

// NewSingleton creates a new Singleton manager for type Type with the given constructor function.
// The constructor will be called lazily when Instance() is first invoked.
//
// Parameters:
//   - constructor: A function that creates and initializes an instance of Type
//
// Returns:
//   - A pointer to a Singleton[Type] manager
func NewSingleton[Type any](constructor Constructor[Type]) *Singleton[Type] {
	return &Singleton[Type]{
		constructor: constructor,
	}
}

// Instance returns the singleton instance of type Type.
// On the first call, it creates the instance using the constructor function.
// Subsequent calls return the same instance without calling the constructor again.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// The sync.Once mechanism ensures the constructor is called exactly once.
//
// Returns:
//   - A pointer to the singleton instance of type Type
func (singleton *Singleton[Type]) Instance() *Type {
	singleton.once.Do(func() {
		singleton.instance, singleton.err = singleton.constructor()
	})
	return singleton.instance
}

// Err returns any error that occurred during instance creation.
// This method should be called after Instance() to check if the constructor
// returned an error during initialization.
//
// Returns:
//   - The error returned by the constructor, or nil if creation was successful
//
// Note: The error is only set once, during the first call to Instance().
// Subsequent calls to Err() will return the same error value.
func (singleton *Singleton[Type]) Err() error {
	return singleton.err
}
