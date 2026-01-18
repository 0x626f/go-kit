package patterns

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// testStruct is a simple struct used for testing the singleton pattern.
// It contains basic fields to verify instance creation and initialization.
type testStruct struct {
	value string
	id    int
}

// TestNewSingleton verifies that NewSingleton creates a valid singleton manager.
//
// Test validates:
//   - NewSingleton returns a non-nil Singleton manager
//   - Constructor function is properly stored
//   - Instance is not created eagerly (lazy initialization)
func TestNewSingleton(t *testing.T) {
	constructor := func() (*testStruct, error) {
		return &testStruct{value: "test", id: 1}, nil
	}

	singleton := NewSingleton(constructor)

	if singleton == nil {
		t.Fatal("NewSingleton returned nil")
	}

	if singleton.constructor == nil {
		t.Error("Constructor was not set")
	}

	if singleton.instance != nil {
		t.Error("Instance should not be created until Instance() is called")
	}
}

// TestSingletonInstance verifies that Instance() creates and returns an instance.
//
// Test validates:
//   - Instance() creates the instance using the constructor
//   - Returned instance is not nil
//   - Instance fields are properly initialized with expected values
//   - No error occurs during successful initialization
func TestSingletonInstance(t *testing.T) {
	constructor := func() (*testStruct, error) {
		return &testStruct{value: "test", id: 42}, nil
	}

	singleton := NewSingleton(constructor)
	instance := singleton.Instance()

	if instance == nil {
		t.Fatal("Instance() returned nil")
	}

	if singleton.Err() != nil {
		t.Fatalf("Unexpected error during initialization: %v", singleton.Err())
	}

	if instance.value != "test" {
		t.Errorf("Expected value 'test', got '%s'", instance.value)
	}

	if instance.id != 42 {
		t.Errorf("Expected id 42, got %d", instance.id)
	}
}

// TestSingletonSameInstance verifies that multiple calls return the same instance.
//
// Test validates:
//   - Multiple calls to Instance() return the exact same instance
//   - Instance identity is preserved across all calls
//   - Singleton pattern guarantees only one instance exists
func TestSingletonSameInstance(t *testing.T) {
	constructor := func() (*testStruct, error) {
		return &testStruct{value: "test", id: 100}, nil
	}

	singleton := NewSingleton(constructor)

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

// TestConstructorCalledOnce verifies that the constructor is called exactly once.
//
// Test validates:
//   - Constructor is invoked exactly once despite multiple Instance() calls
//   - sync.Once mechanism properly prevents duplicate initialization
//   - Lazy initialization occurs only on first Instance() call
func TestConstructorCalledOnce(t *testing.T) {
	var callCount atomic.Int32

	constructor := func() (*testStruct, error) {
		callCount.Add(1)
		return &testStruct{value: "test"}, nil
	}

	singleton := NewSingleton(constructor)

	// Call Instance() multiple times
	for i := 0; i < 10; i++ {
		singleton.Instance()
	}

	if callCount.Load() != 1 {
		t.Errorf("Constructor was called %d times, expected 1", callCount.Load())
	}
}

// TestSingletonConcurrency verifies thread-safety with concurrent access.
//
// Test validates:
//   - Singleton is thread-safe when accessed from multiple goroutines
//   - Constructor is called exactly once even under concurrent load
//   - All goroutines receive the same instance
//   - No race conditions occur during initialization
//
// Uses 100 concurrent goroutines to stress-test the implementation.
func TestSingletonConcurrency(t *testing.T) {
	var callCount atomic.Int32

	constructor := func() (*testStruct, error) {
		callCount.Add(1)
		return &testStruct{value: "concurrent", id: 999}, nil
	}

	singleton := NewSingleton(constructor)

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	instances := make([]*testStruct, numGoroutines)

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

// TestMultipleTypes verifies that Singleton works with different types.
//
// Test validates:
//   - Generic implementation supports different concrete types
//   - Multiple singleton instances can coexist with different type parameters
//   - Type safety is maintained for each singleton instance
//   - Different types don't interfere with each other
func TestMultipleTypes(t *testing.T) {
	type stringContainer struct {
		data string
	}

	type intContainer struct {
		data int
	}

	stringConstructor := func() (*stringContainer, error) {
		return &stringContainer{data: "hello"}, nil
	}

	intConstructor := func() (*intContainer, error) {
		return &intContainer{data: 123}, nil
	}

	stringSingleton := NewSingleton(stringConstructor)
	intSingleton := NewSingleton(intConstructor)

	str := stringSingleton.Instance()
	num := intSingleton.Instance()

	if str.data != "hello" {
		t.Errorf("Expected string 'hello', got '%s'", str.data)
	}

	if num.data != 123 {
		t.Errorf("Expected int 123, got %d", num.data)
	}
}

// TestConstructorError verifies error handling when constructor fails.
//
// Test validates:
//   - Constructor errors are properly captured and stored
//   - Err() method returns the constructor's error
//   - Instance() returns nil when constructor fails
//   - Error is available immediately after Instance() call
func TestConstructorError(t *testing.T) {
	expectedErr := errors.New("initialization failed")

	constructor := func() (*testStruct, error) {
		return nil, expectedErr
	}

	singleton := NewSingleton(constructor)
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

// TestErrorPersistence verifies that errors persist across multiple calls.
//
// Test validates:
//   - Error is set once and persists for the singleton's lifetime
//   - Multiple calls to Instance() don't change the error
//   - Multiple calls to Err() return the same error
//   - Constructor is still called only once even when it fails
func TestErrorPersistence(t *testing.T) {
	var callCount atomic.Int32
	expectedErr := errors.New("persistent error")

	constructor := func() (*testStruct, error) {
		callCount.Add(1)
		return nil, expectedErr
	}

	singleton := NewSingleton(constructor)

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

// TestErrorConcurrency verifies error handling is thread-safe.
//
// Test validates:
//   - Error handling works correctly under concurrent load
//   - All goroutines see the same error
//   - Constructor is called exactly once even when failing
//   - No race conditions occur during error initialization
func TestErrorConcurrency(t *testing.T) {
	var callCount atomic.Int32
	expectedErr := errors.New("concurrent error")

	constructor := func() (*testStruct, error) {
		callCount.Add(1)
		return nil, expectedErr
	}

	singleton := NewSingleton(constructor)

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

// TestNoErrorWhenNil verifies Err() returns nil for successful initialization.
//
// Test validates:
//   - Err() returns nil when constructor succeeds
//   - Successful initialization doesn't set any error
//   - Error state remains nil across multiple Err() calls
func TestNoErrorWhenNil(t *testing.T) {
	constructor := func() (*testStruct, error) {
		return &testStruct{value: "success"}, nil
	}

	singleton := NewSingleton(constructor)
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
