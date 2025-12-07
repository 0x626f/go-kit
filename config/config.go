// Package config provides utilities for loading application configuration from
// various sources such as JSON files, environment variables, and .env files.
// It supports automatic type conversion from string representations to Go types.
package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/0x626f/go-kit/utils"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Configuration file extensions and tag constants
var (
	// jsonExt is the file extension for JSON configuration files
	jsonExt = ".json"
	// envExt is the file extension for environment variable configuration files
	envExt = ".env"
	// tagEnv is the struct tag used to map struct fields to environment variables
	tagEnv = "env"
)

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
	instanceType := reflect.TypeOf(instance).Elem()
	instanceValue := reflect.ValueOf(instance).Elem()

	for index := 0; index < instanceType.NumField(); index++ {
		field := instanceType.Field(index)

		tag := field.Tag.Get(tagEnv)

		if tag == "" {
			continue
		}

		value, exists := os.LookupEnv(tag)

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

		value, exists := data[tag]

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
// Parameters:
//   - ref: The reflect.Value to set
//   - value: The string value to convert and set
func mapPrimaryValue(ref reflect.Value, value string) error {
	var aggError error

	switch ref.Kind() {
	case reflect.String:
		ref.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		ref.SetInt(num)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(value, 10, 64)
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
	case reflect.Float32, reflect.Float64:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		ref.SetFloat(num)
	case reflect.Slice:
		elemType := ref.Type().Elem()

		if elemType.Kind() == reflect.Slice {
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
