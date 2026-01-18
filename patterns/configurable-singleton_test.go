package patterns

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// testConfig is a simple configuration struct used for testing.
type testConfig struct {
	value string
	port  int
}

// testService is a simple service struct used for testing the singleton pattern.
type testService struct {
	config *testConfig
	id     int
}

// TestNewConfigurableSingleton verifies that NewConfigurableSingleton creates a valid singleton manager.
//
// Test validates:
//   - NewConfigurableSingleton returns a non-nil ConfigurableSingleton manager
//   - Constructor function is properly stored
//   - Instance is not created eagerly (lazy initialization)
//   - Config is initially nil
func TestNewConfigurableSingleton(t *testing.T) {
	constructor := func(config *testConfig) (*testService, error) {
		return &testService{config: config, id: 1}, nil
	}

	singleton := NewConfigurableSingleton(constructor)

	if singleton == nil {
		t.Fatal("NewConfigurableSingleton returned nil")
	}

	if singleton.constructor == nil {
		t.Error("Constructor was not set")
	}

	if singleton.instance != nil {
		t.Error("Instance should not be created until Instance() is called")
	}

	if singleton.config != nil {
		t.Error("Config should be nil until WithConfig() is called")
	}
}

// TestConfigurableSingletonInstance verifies that Instance() creates and returns an instance.
//
// Test validates:
//   - Instance() creates the instance using the constructor
//   - Returned instance is not nil
//   - Instance fields are properly initialized with expected values
//   - Configuration is passed to the constructor
//   - No error occurs during successful initialization
func TestConfigurableSingletonInstance(t *testing.T) {
	config := &testConfig{value: "production", port: 8080}

	constructor := func(cfg *testConfig) (*testService, error) {
		return &testService{config: cfg, id: 42}, nil
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)
	instance := singleton.Instance()

	if instance == nil {
		t.Fatal("Instance() returned nil")
	}

	if singleton.Err() != nil {
		t.Fatalf("Unexpected error during initialization: %v", singleton.Err())
	}

	if instance.id != 42 {
		t.Errorf("Expected id 42, got %d", instance.id)
	}

	if instance.config != config {
		t.Error("Instance config should reference the provided config")
	}

	if instance.config.value != "production" {
		t.Errorf("Expected config value 'production', got '%s'", instance.config.value)
	}

	if instance.config.port != 8080 {
		t.Errorf("Expected config port 8080, got %d", instance.config.port)
	}
}

// TestConfigurableSingletonSameInstance verifies that multiple calls return the same instance.
//
// Test validates:
//   - Multiple calls to Instance() return the exact same instance
//   - Instance identity is preserved across all calls
//   - Singleton pattern guarantees only one instance exists
func TestConfigurableSingletonSameInstance(t *testing.T) {
	config := &testConfig{value: "test", port: 3000}

	constructor := func(cfg *testConfig) (*testService, error) {
		return &testService{config: cfg, id: 100}, nil
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)

	instance1 := singleton.Instance()
	instance2 := singleton.Instance()
	instance3 := singleton.Instance()

	if instance1 != instance2 {
		t.Error("First and second calls returned different instances")
	}

	if instance1 != instance3 {
		t.Error("First and third calls returned different instances")
	}

	if instance2 != instance3 {
		t.Error("Second and third calls returned different instances")
	}
}

// TestConfigurableConstructorCalledOnce verifies that the constructor is called exactly once.
//
// Test validates:
//   - Constructor is invoked exactly once despite multiple Instance() calls
//   - sync.Once mechanism properly prevents duplicate initialization
//   - Lazy initialization occurs only on first Instance() call
func TestConfigurableConstructorCalledOnce(t *testing.T) {
	var callCount atomic.Int32
	config := &testConfig{value: "test", port: 5000}

	constructor := func(cfg *testConfig) (*testService, error) {
		callCount.Add(1)
		return &testService{config: cfg}, nil
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)

	// Call Instance() multiple times
	for i := 0; i < 10; i++ {
		singleton.Instance()
	}

	if callCount.Load() != 1 {
		t.Errorf("Constructor was called %d times, expected 1", callCount.Load())
	}
}

// TestConfigurableSingletonConcurrency verifies thread-safety with concurrent access.
//
// Test validates:
//   - ConfigurableSingleton is thread-safe when accessed from multiple goroutines
//   - Constructor is called exactly once even under concurrent load
//   - All goroutines receive the same instance
//   - No race conditions occur during initialization
//
// Uses 100 concurrent goroutines to stress-test the implementation.
func TestConfigurableSingletonConcurrency(t *testing.T) {
	var callCount atomic.Int32
	config := &testConfig{value: "concurrent", port: 9000}

	constructor := func(cfg *testConfig) (*testService, error) {
		callCount.Add(1)
		return &testService{config: cfg, id: 999}, nil
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	instances := make([]*testService, numGoroutines)

	// Launch multiple goroutines to call Instance() concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			instances[index] = singleton.Instance()
		}(i)
	}

	wg.Wait()

	// Verify constructor was called exactly once
	if callCount.Load() != 1 {
		t.Errorf("Constructor was called %d times, expected 1", callCount.Load())
	}

	// Verify all goroutines got the same instance
	firstInstance := instances[0]
	for i := 1; i < numGoroutines; i++ {
		if instances[i] != firstInstance {
			t.Errorf("Goroutine %d got a different instance", i)
		}
	}
}

// TestWithConfig verifies that configuration is passed to the constructor.
//
// Test validates:
//   - WithConfig() properly stores the configuration
//   - Constructor receives the exact configuration provided via WithConfig()
//   - Configuration values are preserved through the initialization process
func TestWithConfig(t *testing.T) {
	expectedConfig := &testConfig{value: "expected", port: 7000}
	var receivedConfig *testConfig

	constructor := func(cfg *testConfig) (*testService, error) {
		receivedConfig = cfg
		return &testService{config: cfg}, nil
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(expectedConfig)
	singleton.Instance()

	if receivedConfig != expectedConfig {
		t.Error("Constructor did not receive the expected config")
	}

	if receivedConfig.value != "expected" {
		t.Errorf("Expected config value 'expected', got '%s'", receivedConfig.value)
	}

	if receivedConfig.port != 7000 {
		t.Errorf("Expected config port 7000, got %d", receivedConfig.port)
	}
}

// TestWithConfigNil verifies behavior when no config is set.
//
// Test validates:
//   - Constructor is called with nil config when WithConfig() is not called
//   - Singleton initialization still succeeds with nil config
//   - Constructor can handle nil config gracefully
func TestWithConfigNil(t *testing.T) {
	var receivedConfig *testConfig

	constructor := func(cfg *testConfig) (*testService, error) {
		receivedConfig = cfg
		return &testService{config: cfg}, nil
	}

	singleton := NewConfigurableSingleton(constructor)
	singleton.Instance()

	if receivedConfig != nil {
		t.Error("Expected nil config when WithConfig is not called")
	}

	if singleton.Err() != nil {
		t.Errorf("Unexpected error: %v", singleton.Err())
	}
}

// TestWithConfigMethodChaining verifies that WithConfig returns the singleton for chaining.
//
// Test validates:
//   - WithConfig() returns the same ConfigurableSingleton instance
//   - Method chaining pattern is properly supported
//   - Fluent interface design is maintained
func TestWithConfigMethodChaining(t *testing.T) {
	constructor := func(cfg *testConfig) (*testService, error) {
		return &testService{config: cfg}, nil
	}

	config := &testConfig{value: "test", port: 4000}
	singleton := NewConfigurableSingleton(constructor)
	result := singleton.WithConfig(config)

	if result != singleton {
		t.Error("WithConfig should return the same singleton instance for method chaining")
	}
}

// TestWithConfigAfterInstance verifies that setting config after Instance() has no effect.
//
// Test validates:
//   - WithConfig() only affects initialization if called before Instance()
//   - Subsequent WithConfig() calls after Instance() are ignored
//   - First config set is preserved throughout the singleton's lifetime
//   - Immutability of singleton configuration after initialization
func TestWithConfigAfterInstance(t *testing.T) {
	var receivedConfig *testConfig

	constructor := func(cfg *testConfig) (*testService, error) {
		receivedConfig = cfg
		return &testService{config: cfg}, nil
	}

	config1 := &testConfig{value: "first", port: 1000}
	config2 := &testConfig{value: "second", port: 2000}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config1)
	singleton.Instance()

	// This should have no effect since instance is already created
	singleton.WithConfig(config2)

	if receivedConfig != config1 {
		t.Error("Expected config from first WithConfig call")
	}

	if receivedConfig.value != "first" {
		t.Errorf("Expected config value 'first', got '%s'", receivedConfig.value)
	}
}

// TestWithConfigMultipleCalls verifies that only the first WithConfig call takes effect.
//
// Test validates:
//   - Multiple WithConfig() calls before Instance() only use the first config
//   - Subsequent config attempts are ignored
//   - First-wins semantics for configuration
func TestWithConfigMultipleCalls(t *testing.T) {
	var receivedConfig *testConfig

	constructor := func(cfg *testConfig) (*testService, error) {
		receivedConfig = cfg
		return &testService{config: cfg}, nil
	}

	config1 := &testConfig{value: "first", port: 1111}
	config2 := &testConfig{value: "second", port: 2222}
	config3 := &testConfig{value: "third", port: 3333}

	singleton := NewConfigurableSingleton(constructor)
	singleton.WithConfig(config1)
	singleton.WithConfig(config2)
	singleton.WithConfig(config3)
	singleton.Instance()

	if receivedConfig != config1 {
		t.Error("Expected first config to be used")
	}

	if receivedConfig.value != "first" {
		t.Errorf("Expected config value 'first', got '%s'", receivedConfig.value)
	}
}

// TestConfigMethod verifies that Config() returns the configured value.
//
// Test validates:
//   - Config() returns the configuration set via WithConfig()
//   - Config() returns nil if WithConfig() was not called
//   - Config retrieval works before and after Instance() call
func TestConfigMethod(t *testing.T) {
	constructor := func(cfg *testConfig) (*testService, error) {
		return &testService{config: cfg}, nil
	}

	// Test with config set
	config := &testConfig{value: "test", port: 6000}
	singleton1 := NewConfigurableSingleton(constructor).WithConfig(config)

	if singleton1.Config() != config {
		t.Error("Config() should return the configured value")
	}

	singleton1.Instance()

	if singleton1.Config() != config {
		t.Error("Config() should still return the configured value after Instance()")
	}

	// Test without config set
	singleton2 := NewConfigurableSingleton(constructor)

	if singleton2.Config() != nil {
		t.Error("Config() should return nil when no config is set")
	}
}

// TestConfigurableMultipleTypes verifies that ConfigurableSingleton works with different types.
//
// Test validates:
//   - Generic implementation supports different concrete types
//   - Multiple singleton instances can coexist with different type parameters
//   - Type safety is maintained for each singleton instance
//   - Different types and configs don't interfere with each other
func TestConfigurableMultipleTypes(t *testing.T) {
	type stringConfig struct {
		prefix string
	}

	type stringContainer struct {
		data string
	}

	type intConfig struct {
		multiplier int
	}

	type intContainer struct {
		data int
	}

	stringConstructor := func(cfg *stringConfig) (*stringContainer, error) {
		prefix := ""
		if cfg != nil {
			prefix = cfg.prefix
		}
		return &stringContainer{data: prefix + "hello"}, nil
	}

	intConstructor := func(cfg *intConfig) (*intContainer, error) {
		multiplier := 1
		if cfg != nil {
			multiplier = cfg.multiplier
		}
		return &intContainer{data: 123 * multiplier}, nil
	}

	strCfg := &stringConfig{prefix: "test-"}
	intCfg := &intConfig{multiplier: 2}

	stringSingleton := NewConfigurableSingleton(stringConstructor).WithConfig(strCfg)
	intSingleton := NewConfigurableSingleton(intConstructor).WithConfig(intCfg)

	str := stringSingleton.Instance()
	num := intSingleton.Instance()

	if str.data != "test-hello" {
		t.Errorf("Expected string 'test-hello', got '%s'", str.data)
	}

	if num.data != 246 {
		t.Errorf("Expected int 246, got %d", num.data)
	}
}

// TestConfigurableConstructorError verifies error handling when constructor fails.
//
// Test validates:
//   - Constructor errors are properly captured and stored
//   - Err() method returns the constructor's error
//   - Instance() returns nil when constructor fails
//   - Error is available immediately after Instance() call
func TestConfigurableConstructorError(t *testing.T) {
	expectedErr := errors.New("initialization failed")
	config := &testConfig{value: "error", port: 9999}

	constructor := func(cfg *testConfig) (*testService, error) {
		return nil, expectedErr
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)
	instance := singleton.Instance()

	if instance != nil {
		t.Error("Instance should be nil when constructor returns error")
	}

	if singleton.Err() == nil {
		t.Fatal("Err() should return the constructor error")
	}

	if singleton.Err() != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, singleton.Err())
	}
}

// TestConfigurableErrorPersistence verifies that errors persist across multiple calls.
//
// Test validates:
//   - Error is set once and persists for the singleton's lifetime
//   - Multiple calls to Instance() don't change the error
//   - Multiple calls to Err() return the same error
//   - Constructor is still called only once even when it fails
func TestConfigurableErrorPersistence(t *testing.T) {
	var callCount atomic.Int32
	expectedErr := errors.New("persistent error")
	config := &testConfig{value: "error", port: 8888}

	constructor := func(cfg *testConfig) (*testService, error) {
		callCount.Add(1)
		return nil, expectedErr
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)

	// Call Instance() multiple times
	for i := 0; i < 5; i++ {
		instance := singleton.Instance()
		if instance != nil {
			t.Errorf("Call %d: Instance should be nil", i)
		}

		err := singleton.Err()
		if err != expectedErr {
			t.Errorf("Call %d: Expected error %v, got %v", i, expectedErr, err)
		}
	}

	if callCount.Load() != 1 {
		t.Errorf("Constructor was called %d times, expected 1", callCount.Load())
	}
}

// TestConfigurableErrorConcurrency verifies error handling is thread-safe.
//
// Test validates:
//   - Error handling works correctly under concurrent load
//   - All goroutines see the same error
//   - Constructor is called exactly once even when failing
//   - No race conditions occur during error initialization
func TestConfigurableErrorConcurrency(t *testing.T) {
	var callCount atomic.Int32
	expectedErr := errors.New("concurrent error")
	config := &testConfig{value: "error", port: 7777}

	constructor := func(cfg *testConfig) (*testService, error) {
		callCount.Add(1)
		return nil, expectedErr
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errs := make([]error, numGoroutines)

	// Launch multiple goroutines to call Instance() concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			singleton.Instance()
			errs[index] = singleton.Err()
		}(i)
	}

	wg.Wait()

	// Verify constructor was called exactly once
	if callCount.Load() != 1 {
		t.Errorf("Constructor was called %d times, expected 1", callCount.Load())
	}

	// Verify all goroutines got the same error
	for i := 0; i < numGoroutines; i++ {
		if errs[i] != expectedErr {
			t.Errorf("Goroutine %d got error %v, expected %v", i, errs[i], expectedErr)
		}
	}
}

// TestConfigurableNoErrorWhenNil verifies Err() returns nil for successful initialization.
//
// Test validates:
//   - Err() returns nil when constructor succeeds
//   - Successful initialization doesn't set any error
//   - Error state remains nil across multiple Err() calls
func TestConfigurableNoErrorWhenNil(t *testing.T) {
	config := &testConfig{value: "success", port: 6666}

	constructor := func(cfg *testConfig) (*testService, error) {
		return &testService{config: cfg}, nil
	}

	singleton := NewConfigurableSingleton(constructor).WithConfig(config)
	instance := singleton.Instance()

	if instance == nil {
		t.Fatal("Instance should not be nil for successful initialization")
	}

	if singleton.Err() != nil {
		t.Errorf("Err() should return nil for successful initialization, got %v", singleton.Err())
	}

	// Call Err() multiple times to ensure consistency
	for i := 0; i < 3; i++ {
		if err := singleton.Err(); err != nil {
			t.Errorf("Call %d: Err() should return nil, got %v", i, err)
		}
	}
}
