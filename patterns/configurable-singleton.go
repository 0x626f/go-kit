package patterns

import (
	"sync"
)

// ConfigurableConstructor is a function type that creates and initializes an instance of type Type
// using the provided configuration of type Config.
// It returns a pointer to Type and an error if initialization fails.
type ConfigurableConstructor[Type any, Config any] func(config *Config) (*Type, error)

// ConfigurableSingleton provides a thread-safe generic implementation of the singleton pattern
// with configuration support. It ensures that only one instance of type Type is created,
// even when accessed concurrently from multiple goroutines.
//
// The instance is created lazily on the first call to Instance() using the provided
// constructor function and configuration. The sync.Once mechanism guarantees that the
// constructor is called exactly once.
//
// Example usage:
//
//	type Database struct {
//	    conn *sql.DB
//	    maxConnections int
//	}
//
//	type DBConfig struct {
//	    Host string
//	    Port int
//	    MaxConnections int
//	}
//
//	dbConstructor := func(config *DBConfig) (*Database, error) {
//	    conn, err := initConnection(config.Host, config.Port)
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &Database{conn: conn, maxConnections: config.MaxConnections}, nil
//	}
//
//	config := &DBConfig{Host: "localhost", Port: 5432, MaxConnections: 100}
//	dbSingleton := NewConfigurableSingleton(dbConstructor).WithConfig(config)
//	db := dbSingleton.Instance() // First call creates the instance
//	if dbSingleton.Err() != nil {
//	    // Handle initialization error
//	}
//	db2 := dbSingleton.Instance() // Returns the same instance
type ConfigurableSingleton[Type any, Config any] struct {
	once        sync.Once
	constructor ConfigurableConstructor[Type, Config]
	instance    *Type
	config      *Config
	err         error
}

// NewConfigurableSingleton creates a new ConfigurableSingleton manager for type Type
// with the given constructor function that accepts configuration of type Config.
// The constructor will be called lazily when Instance() is first invoked.
//
// Parameters:
//   - constructor: A function that creates and initializes an instance of Type using Config
//
// Returns:
//   - A pointer to a ConfigurableSingleton[Type, Config] manager
func NewConfigurableSingleton[Type any, Config any](constructor ConfigurableConstructor[Type, Config]) *ConfigurableSingleton[Type, Config] {
	return &ConfigurableSingleton[Type, Config]{
		constructor: constructor,
	}
}

// WithConfig sets the configuration that will be passed to the constructor function.
// This method should be called before the first call to Instance().
//
// Parameters:
//   - config: The configuration to use when initializing the singleton instance
//
// Returns:
//   - The ConfigurableSingleton instance (for method chaining)
//
// Note: Calling this method after Instance() has been called has no effect,
// as the instance has already been created with the previously set configuration.
// If called multiple times before Instance(), only the first config is used.
func (singleton *ConfigurableSingleton[Type, Config]) WithConfig(config *Config) *ConfigurableSingleton[Type, Config] {
	if singleton.config == nil {
		singleton.config = config
	}
	return singleton
}

// Instance returns the singleton instance of type Type.
// On the first call, it creates the instance using the constructor function
// with the configured settings. Subsequent calls return the same instance
// without calling the constructor again.
//
// This method is thread-safe and can be called concurrently from multiple goroutines.
// The sync.Once mechanism ensures the constructor is called exactly once.
//
// Returns:
//   - A pointer to the singleton instance of type Type
func (singleton *ConfigurableSingleton[Type, Config]) Instance() *Type {
	singleton.once.Do(func() {
		singleton.instance, singleton.err = singleton.constructor(singleton.config)
	})
	return singleton.instance
}

// Config returns the configuration used for singleton initialization.
// This can be useful for inspecting what configuration was applied.
//
// Returns:
//   - A pointer to the Config used, or nil if WithConfig was not called
func (singleton *ConfigurableSingleton[Type, Config]) Config() *Config {
	return singleton.config
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
func (singleton *ConfigurableSingleton[Type, Config]) Err() error {
	return singleton.err
}
