// Package config provides utilities for loading application configuration from
// various sources such as JSON files, environment variables, and .env files.
// It supports automatic type conversion from string representations to Go types.
package env

import (
	"bufio"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/0x626f/go-kit/utils"
)

// Configuration file extensions and tag constants
var (
	// jsonExt is the file extension for JSON configuration files
	jsonExt = ".json"
	// envExt is the file extension for environment variable configuration files
	envExt = ".env"
	// tagEnv is the struct tag used to map struct fields to environment variables
	tagEnv     = "env"
	tagDefault = "default"
)

var prefix string

// SetEnvPrefix sets a global prefix for all environment variable lookups.
// When a prefix is set, all environment variable names will be automatically
// prefixed with the given name followed by an underscore.
//
// This is useful for namespacing environment variables in multi-tenant applications
// or when running multiple instances of the same application with different configurations.
//
// Parameters:
//   - name: The prefix to use for environment variable names
//
// Example:
//
//	// Set prefix to "MYAPP"
//	config.SetEnvPrefix("MYAPP")
//
//	// Now GetEnv("PORT", "8080") will look for "MYAPP_PORT" instead of "PORT"
//	port := config.GetEnv("PORT", "8080")
//
//	// Similarly, GetEnvAs will use the prefix
//	timeout := config.GetEnvAs("TIMEOUT", 30)  // Looks for "MYAPP_TIMEOUT"
func SetEnvPrefix(name string) {
	prefix = name
}

// GetEnvPrefix returns the current global environment variable prefix.
// If no prefix has been set, it returns an empty string.
//
// Returns:
//   - string: The current prefix, or empty string if not set
//
// Example:
//
//	config.SetEnvPrefix("MYAPP")
//	prefix := config.GetEnvPrefix() // Returns "MYAPP"
func GetEnvPrefix() string {
	return prefix
}

// getPrefixedEnv applies the global prefix to an environment variable name.
// If no prefix is set, the original environment variable name is returned unchanged.
// If a prefix is set, it returns the name in the format: "PREFIX_ENVNAME".
//
// This is an internal helper function used by all environment variable retrieval functions
// to ensure consistent prefix handling across the package.
//
// Parameters:
//   - env: The base environment variable name
//
// Returns:
//   - string: The prefixed environment variable name, or the original name if no prefix is set
func getPrefixedEnv(env string) string {
	if prefix == "" {
		return env
	}
	return prefix + "_" + env
}

// GetEnv retrieves the value of an environment variable with a fallback default.
// If the environment variable exists, its value is returned. Otherwise, the fallback
// value is returned.
//
// This is a convenience function that wraps os.LookupEnv and provides a simpler
// API for cases where a default value is acceptable.
//
// Parameters:
//   - name: The name of the environment variable to retrieve
//   - fallback: The default value to return if the environment variable is not set
//
// Returns:
//   - string: The value of the environment variable if it exists, otherwise the fallback value
//
// Example:
//
//	// Get database host with default
//	dbHost := config.GetEnv("DB_HOST", "localhost")
//
//	// Get port with default
//	port := config.GetEnv("PORT", "8080")
//
//	// Get optional feature flag
//	enableFeature := config.GetEnv("ENABLE_FEATURE_X", "false")
func GetEnv(name, fallback string) string {
	if value, exists := os.LookupEnv(getPrefixedEnv(name)); exists {
		return value
	}
	return fallback
}

// canConvertFromEnv checks if a reflect.Kind can be converted from a string environment variable value.
// This function supports primitive types and slices as convertible types.
//
// Supported kinds include:
//   - Numeric types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
//   - Floating point: float32, float64
//   - Boolean: bool
//   - String: string
//   - Slice: slice of supported types
//
// Parameters:
//   - kind: The reflect.Kind to check for conversion support
//
// Returns:
//   - bool: True if the kind can be converted from an environment variable string, false otherwise
func canConvertFromEnv(kind reflect.Kind) bool {
	casted := uint(kind)

	if casted > 0 && (casted <= 14 || casted == 23 || casted == 24) {
		return true
	}
	return false
}

// GetEnvAs retrieves an environment variable and converts it to the specified type T.
// If the environment variable doesn't exist or conversion fails, the fallback value is returned.
//
// This generic function provides type-safe environment variable retrieval with automatic
// type conversion. It supports all primitive types (int, uint, float, bool, string) and
// slices of these types.
//
// Type parameters:
//   - T: The target type for the environment variable value. Must be a convertible type.
//
// Parameters:
//   - name: The name of the environment variable to retrieve
//   - fallback: The default value to return if the variable is not set or conversion fails
//
// Returns:
//   - T: The converted value of the environment variable, or the fallback value
//
// Example:
//
//	// Get an integer port with default
//	port := config.GetEnvAs("PORT", 8080)
//
//	// Get a boolean flag with default
//	debugMode := config.GetEnvAs("DEBUG", false)
//
//	// Get a float64 timeout with default
//	timeout := config.GetEnvAs("TIMEOUT", 30.0)
//
//	// Get a slice of integers
//	ids := config.GetEnvAs("ALLOWED_IDS", []int{1, 2, 3})
func GetEnvAs[T any](name string, fallback T) T {
	if value, exists := os.LookupEnv(getPrefixedEnv(name)); exists {
		instance := utils.NewInstanceOf[T]()
		instanceType := reflect.TypeOf(instance).Elem()

		if !canConvertFromEnv(instanceType.Kind()) {
			return fallback
		}

		instanceValue := reflect.ValueOf(instance).Elem()
		err := mapPrimaryValue(instanceValue, value)

		if err == nil {
			return *instance
		}
	}
	return fallback
}

// GetEnvDuration retrieves an environment variable and parses it as a time.Duration.
// If the environment variable doesn't exist or parsing fails, the fallback duration is returned.
//
// The duration string should be in the format accepted by time.ParseDuration, such as
// "300ms", "1.5h", "2h45m", etc. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
//
// Parameters:
//   - name: The name of the environment variable to retrieve
//   - fallback: The default duration to return if the variable is not set or parsing fails
//
// Returns:
//   - time.Duration: The parsed duration from the environment variable, or the fallback value
//
// Example:
//
//	// Get request timeout with default
//	timeout := config.GetEnvDuration("REQUEST_TIMEOUT", 30*time.Second)
//
//	// Get retry delay with default
//	retryDelay := config.GetEnvDuration("RETRY_DELAY", 5*time.Second)
//
//	// Get connection idle timeout with default
//	idleTimeout := config.GetEnvDuration("IDLE_TIMEOUT", 90*time.Second)
//
//	// Environment variable examples:
//	// REQUEST_TIMEOUT=45s
//	// RETRY_DELAY=500ms
//	// IDLE_TIMEOUT=2m30s
func GetEnvDuration(name string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(getPrefixedEnv(name)); exists {
		duration, err := time.ParseDuration(value)

		if err == nil {
			return duration
		}
	}
	return fallback
}

// FromFile loads configuration from a file and maps it to a struct of type T.
// Supported file types are JSON (.json extension) and environment files (.env extension).
//
// Type parameters:
//   - T: The struct type to which the configuration will be mapped
//
// Parameters:
//   - filename: Path to the configuration file
//
// Returns:
//   - *T: A pointer to a struct of type T populated with configuration values
//   - error: An error if the file can't be read, has an unsupported extension, or if mapping fails
//
// Example:
//
//	type ServerConfig struct {
//	    Host string `json:"host" env:"SERVER_HOST"`
//	    Port int    `json:"port" env:"SERVER_PORT"`
//	}
//
//	// Load from JSON file
//	config, err := config.FromFile[ServerConfig]("server.json")
func FromFile[T any](filename string) (*T, error) {
	if !utils.IsObject[T]() {
		return nil, fmt.Errorf("underlying type must be a struct")
	}

	if isJson(filename) {
		return mapJSONConfig[T](filename)
	}

	if isEnv(filename) {
		return mapEnvConfig[T](filename)
	}

	return nil, fmt.Errorf("unsupported extension for %v", filename)
}

// FromEnvs loads configuration directly from environment variables and maps
// it to a struct of type T based on the "env" struct tag.
//
// Type parameters:
//   - T: The struct type to which the configuration will be mapped
//
// Returns:
//   - *T: A pointer to a struct of type T populated with configuration values
//   - error: An error if T is not a struct type
//
// Example:
//
//	type DatabaseConfig struct {
//	    Host string `env:"DB_HOST"`
//	    Port int `env:"DB_PORT"`
//	    User string `env:"DB_USER"`
//	    Password string `env:"DB_PASSWORD"`
//	}
//
//	// Load from environment variables
//	dbConfig, err: = config.FromEnvs[DatabaseConfig]()
func FromEnvs[T any]() (*T, error) {
	if !utils.IsObject[T]() {
		return nil, fmt.Errorf("underlying type must be a struct")
	}

	var err error
	instance := utils.NewInstanceOf[T]()
	instanceValue := reflect.ValueOf(instance).Elem()

	err = mapStructFromEnvs(instanceValue, "")

	return instance, err
}

// LoadEnvs loads environment variables from a file and sets them
// in the current process environment using os.Setenv.
//
// Parameters:
//   - filename: Path to the .env file
//
// Returns:
//   - error: An error if the file can't be read or if setting any environment variable fails
//
// Example:
//
//	// Load environment variables from .env file
//	err := config.LoadEnvs(".env")
func LoadEnvs(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		separator := strings.Index(line, "=")

		if separator == -1 {
			continue
		}

		key, value := line[:separator], line[separator+1:]
		err = os.Setenv(key, value)

		if err != nil {
			break
		}
	}

	err = errors.Join(err, file.Close())

	return err
}

// isJson checks if the given filename has a JSON file extension (.json).
//
// Parameters:
//   - filename: The filename to check
//
// Returns:
//   - bool: True if the file has a .json extension, false otherwise
func isJson(filename string) bool {
	return strings.HasSuffix(filename, jsonExt)
}

// isEnv checks if the given filename has an environment file extension (.env).
//
// Parameters:
//   - filename: The filename to check
//
// Returns:
//   - bool: True if the file has a .env extension, false otherwise
func isEnv(filename string) bool {
	return strings.HasSuffix(filename, envExt)
}

// mapJSONConfig loads a JSON configuration file and maps it to a struct of type T.
//
// Type parameters:
//   - T: The struct type to which the configuration will be mapped
//
// Parameters:
//   - filename: Path to the JSON configuration file
//
// Returns:
//   - *T: A pointer to a struct of type T populated with configuration values
//   - error: An error if the file can't be read or if JSON unmarshaling fails
func mapJSONConfig[T any](filename string) (*T, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	instance := utils.NewInstanceOf[T]()

	err = json.Unmarshal(data, instance)

	if err != nil {
		return nil, err
	}

	return instance, nil
}

// mapEnvConfig loads an environment configuration file (.env) and maps it to a struct of type T
// based on the "env" struct tag.
//
// Type parameters:
//   - T: The struct type to which the configuration will be mapped
//
// Parameters:
//   - filename: Path to the .env configuration file
//
// Returns:
//   - *T: A pointer to a struct of type T populated with configuration values
//   - error: An error if the file can't be read or if mapping fails
func mapEnvConfig[T any](filename string) (*T, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	data := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		separator := strings.Index(line, "=")

		if separator == -1 {
			continue
		}

		key, value := line[:separator], line[separator+1:]
		data[key] = value
	}

	err = file.Close()

	if err != nil {
		return nil, err
	}

	instance := utils.NewInstanceOf[T]()
	instanceType := reflect.TypeOf(instance).Elem()
	instanceValue := reflect.ValueOf(instance).Elem()

	for index := 0; index < instanceType.NumField(); index++ {
		field := instanceType.Field(index)

		tag := field.Tag.Get(tagEnv)

		if tag == "" {
			continue
		}

		value, exists := data[getPrefixedEnv(tag)]

		if !exists {
			continue
		}

		ref := instanceValue.Field(index)

		if !ref.CanSet() {
			continue
		}

		err = errors.Join(err, mapPrimaryValue(ref, value))
	}

	return instance, err
}

// mapPrimaryValue converts a string value to the appropriate type and sets it in the given reflect.Value.
// This function handles primitive types (string, numeric types, bool) and slices of these types.
// For slices, it splits the string by comma and converts each element.
//
// Supported types:
//   - String: Direct assignment
//   - Numeric types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
//   - Floating point: float32, float64
//   - Boolean: true/false or 1/0
//   - Slices: Comma-separated values (e.g., "1,2,3" for []int)
//
// Parameters:
//   - ref: The reflect.Value to set
//   - value: The string value to convert and set
//
// Returns:
//   - error: An error if conversion fails or if the type is unsupported
func mapPrimaryValue(ref reflect.Value, value string) error {
	refType := ref.Type()
	if utils.IsInstanceOf[time.Duration](refType) {
		dur, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		ref.SetInt(int64(dur))
		return nil
	}

	if utils.Implements[encoding.TextUnmarshaler](refType) {
		ptr := ref.Addr()
		m := ptr.MethodByName("UnmarshalText")
		result := m.Call([]reflect.Value{reflect.ValueOf([]byte(value))})
		if !result[0].IsNil() {
			ref.Set(result[0].Elem())
		}
		return nil
	}

	var aggError error
	switch ref.Kind() {
	case reflect.String:
		ref.SetString(value)
	case reflect.Int, reflect.Int64:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		ref.SetInt(num)
	case reflect.Int8:
		num, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}
		ref.SetInt(num)
	case reflect.Int16:
		num, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
		ref.SetInt(num)
	case reflect.Int32:
		num, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		ref.SetInt(num)
	case reflect.Uint, reflect.Uint64:
		num, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		ref.SetUint(num)
	case reflect.Uint8:
		num, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return err
		}
		ref.SetUint(num)
	case reflect.Uint16:
		num, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}
		ref.SetUint(num)
	case reflect.Uint32:
		num, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		ref.SetUint(num)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		ref.SetBool(b)
	case reflect.Float64:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		ref.SetFloat(num)
	case reflect.Float32:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		ref.SetFloat(num)
	case reflect.Slice:
		elemType := refType.Elem()

		if elemType.Kind() == reflect.Slice || elemType.Kind() == reflect.Array {
			ref.SetZero()
			return fmt.Errorf("couldn't map dimensional arrays from .env")
		}

		values := strings.Split(value, ",")

		slice := reflect.MakeSlice(ref.Type(), len(values), len(values))

		for index, item := range values {
			err := mapPrimaryValue(slice.Index(index), item)
			if err != nil {
				aggError = errors.Join(aggError, err)
			}
		}
		ref.Set(slice)
	default:
		ref.SetZero()
	}
	return aggError
}

// addNestedPrefix adds a prefix to an environment variable name for nested struct field mapping.
// If the prefix is empty, the envName is returned unchanged.
// This is used internally to construct hierarchical environment variable names for nested structs.
//
// Parameters:
//   - envName: The environment variable name to prefix
//   - prefix: The prefix to prepend
//
// Returns:
//   - string: The prefixed environment variable name in the format "PREFIX_ENVNAME"
//
// Example:
//
//	addNestedPrefix("HOST", "DB") // Returns "DB_HOST"
//	addNestedPrefix("HOST", "")   // Returns "HOST"
func addNestedPrefix(envName, prefix string) string {
	if prefix == "" {
		return envName
	}
	return prefix + "_" + envName
}

// mapStructFromEnvs recursively maps environment variables to struct fields.
// This function supports nested structs, struct pointers, and all primitive types.
//
// For nested structs, it builds hierarchical environment variable names by combining
// prefixes from parent structs with child field names. For example, a nested struct
// with env tag "DB" containing a field with env tag "HOST" will look for "DB_HOST".
//
// Parameters:
//   - ref: The reflect.Value of the struct to populate
//   - prefix: The accumulated prefix for nested struct fields
//
// Returns:
//   - error: An aggregated error if any field mapping fails
//
// Example struct mapping:
//
//	type Config struct {
//	    Database struct {
//	        Host string `env:"HOST"`
//	        Port int    `env:"PORT"`
//	    } `env:"DB"`
//	}
//	// Will look for environment variables: DB_HOST, DB_PORT
func mapStructFromEnvs(ref reflect.Value, prefix string) (err error) {
	refType := ref.Type()
	for index := 0; index < ref.NumField(); index++ {
		field := refType.Field(index)

		tag := field.Tag.Get(tagEnv)

		if tag == "-" {
			continue
		}

		fieldRef := ref.Field(index)

		if !fieldRef.CanSet() {
			continue
		}

		if fieldRef.Kind() == reflect.Struct {
			err = errors.Join(err, mapStructFromEnvs(fieldRef, addNestedPrefix(tag, prefix)))
		} else if fieldRef.Kind() == reflect.Pointer && fieldRef.Type().Elem().Kind() == reflect.Struct {
			fieldRef.Set(reflect.New(fieldRef.Type().Elem()))
			err = errors.Join(err, mapStructFromEnvs(fieldRef.Elem(), addNestedPrefix(tag, prefix)))
		} else {
			if tag == "" {
				continue
			}

			// Try to get value from environment variable first
			value, exists := os.LookupEnv(getPrefixedEnv(addNestedPrefix(tag, prefix)))

			// If env var doesn't exist, try to use default tag value
			if !exists {
				value = field.Tag.Get(tagDefault)
				// If no default either, skip this field
				if value == "" {
					continue
				}
			}

			err = errors.Join(err, mapPrimaryValue(fieldRef, value))
		}
	}
	return
}
