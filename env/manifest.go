package env

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Manifest is a fluent builder for loading configuration from multiple sources
// into a strongly-typed struct. It provides a chainable API for configuring
// application settings from environment files, variables, and custom properties.
//
// Type parameters:
//   - T: The target configuration struct type. Must be a struct with appropriate
//     "env" tags on its fields for environment variable mapping.
//
// Example:
//
//	type AppConfig struct {
//	    Host string `env:"HOST"`
//	    Port int    `env:"PORT"`
//	}
//
//	config, err := NewManifest[AppConfig]().
//	    WithPrefix("MYAPP").
//	    WithSource(".env").
//	    Load()
type Manifest[T any] struct{}

// NewManifest creates a new Manifest instance for loading configuration of type T.
// This is the entry point for the fluent configuration builder API.
//
// Type parameters:
//   - T: The target configuration struct type
//
// Returns:
//   - *Manifest[T]: A new Manifest instance ready to be configured
//
// Example:
//
//	manifest := NewManifest[AppConfig]()
func NewManifest[T any]() *Manifest[T] {
	return &Manifest[T]{}
}

// WithPrefix sets a global environment variable prefix based on the value
// of the specified environment variable. This is useful for multi-tenant
// applications or when running multiple instances with different configurations.
//
// The value retrieved from the environment variable will be used as a prefix
// for all subsequent environment variable lookups in the format "PREFIX_VARNAME".
//
// Parameters:
//   - variable: The name of the environment variable containing the application name/prefix
//
// Returns:
//   - *Manifest[T]: The Manifest instance for method chaining
//
// Example:
//
//	// If APP_NAME=myservice, all env vars will be prefixed with "myservice_"
//	manifest := NewManifest[Config]().
//	    WithPrefix("APP_NAME").
//	    Load()
//	// Will look for "myservice_HOST", "myservice_PORT", etc.
func (manifest *Manifest[T]) WithPrefix(variable string) *Manifest[T] {
	appName := GetEnv(variable, "")
	SetEnvPrefix(appName)

	return manifest
}

// WithSource loads environment variables from a file at the specified path.
// The file should be in .env format (KEY=VALUE pairs, one per line).
//
// This method will panic if the file cannot be read or parsed. Use this when
// the configuration file is required for the application to run.
//
// Parameters:
//   - filepath: The path to the .env file to load
//
// Returns:
//   - *Manifest[T]: The Manifest instance for method chaining
//
// Panics:
//   - If the file cannot be read or contains invalid format
//
// Example:
//
//	manifest := NewManifest[Config]().
//	    WithSource(".env").
//	    WithSource("config/.env.local").
//	    Load()
func (manifest *Manifest[T]) WithSource(filepath string) *Manifest[T] {
	err := LoadEnvs(filepath)
	if err != nil {
		panic(err)
	}

	return manifest
}

// WithAbsoluteSource loads environment variables from a file relative to the
// directory where this package is located. This is useful for loading
// configuration files that are bundled with your application code.
//
// Parameters:
//   - filename: The filename (not full path) of the .env file to load
//
// Returns:
//   - *Manifest[T]: The Manifest instance for method chaining
//
// Panics:
//   - If the runtime caller information cannot be obtained
//   - If the file cannot be read or contains invalid format
//
// Example:
//
//	// Loads config.env from the same directory as the env package
//	manifest := NewManifest[Config]().
//	    WithAbsoluteSource("config.env").
//	    Load()
func (manifest *Manifest[T]) WithAbsoluteSource(filename string) *Manifest[T] {
	_, path, _, ok := runtime.Caller(1)
	if !ok {
		panic("couldn't get path to executable")
	}

	err := LoadEnvs(fmt.Sprintf("%v/%v", filepath.Dir(path), filename))
	if err != nil {
		panic(err)
	}

	return manifest
}

// WithRelativeSource loads environment variables from a file relative to a path
// specified in an environment variable. The environment variable should contain
// a full file path, and the filename will be resolved relative to that path's directory.
//
// This is useful when the configuration file location depends on runtime
// environment settings or deployment configurations.
//
// Parameters:
//   - variable: The name of the environment variable containing a reference path
//   - filename: The filename to load relative to the variable's path directory
//
// Returns:
//   - *Manifest[T]: The Manifest instance for method chaining
//
// Panics:
//   - If the file cannot be read or contains invalid format
//
// Example:
//
//	// If CONFIG_PATH=/app/config/main.conf
//	// This will load /app/config/secrets.env
//	manifest := NewManifest[Config]().
//	    WithRelativeSource("CONFIG_PATH", "secrets.env").
//	    Load()
func (manifest *Manifest[T]) WithRelativeSource(variable, filename string) *Manifest[T] {
	path := GetEnv(variable, "")

	err := LoadEnvs(fmt.Sprintf("%v/%v", filepath.Dir(path), filename))
	if err != nil {
		panic(err)
	}

	return manifest
}

// WithProperty sets a single environment variable programmatically.
// This is useful for setting default values, overriding values, or injecting
// configuration at runtime.
//
// The property is set in the current process environment and will be available
// to all subsequent operations including Load().
//
// Parameters:
//   - key: The environment variable name to set
//   - value: The value to assign to the environment variable
//
// Returns:
//   - *Manifest[T]: The Manifest instance for method chaining
//
// Panics:
//   - If the environment variable cannot be set (rare, usually OS-level restrictions)
//
// Example:
//
//	manifest := NewManifest[Config]().
//	    WithProperty("DEBUG", "true").
//	    WithProperty("LOG_LEVEL", "info").
//	    Load()
func (manifest *Manifest[T]) WithProperty(key, value string) *Manifest[T] {
	err := SetEnv(key, value)
	if err != nil {
		panic(err)
	}
	return manifest
}

// Load finalizes the configuration loading process and maps all environment
// variables to a struct of type T based on the "env" struct tags.
//
// This method should be called after all configuration sources have been added
// via the With* methods. It reads the current process environment (which includes
// all loaded files and properties) and maps them to the target struct.
//
// Returns:
//   - *T: A pointer to a populated struct of type T
//   - error: An error if T is not a struct type or if mapping fails
//
// Example:
//
//	type Config struct {
//	    Host string `env:"HOST"`
//	    Port int    `env:"PORT"`
//	}
//
//	config, err := NewManifest[Config]().
//	    WithSource(".env").
//	    WithProperty("DEBUG", "true").
//	    Load()
//	if err != nil {
//	    log.Fatal(err)
//	}
func (manifest *Manifest[T]) Load() (*T, error) {
	return FromEnvs[T]()
}
