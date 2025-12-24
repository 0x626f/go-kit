package patterns

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

// testStruct is a simple struct used for testing
type testStruct struct {
	value string
	id    int
}

// TestNewSingleton verifies that NewSingleton creates a valid singleton instance
func TestNewSingleton(t *testing.T) {
	constructor := func(ctx context.Context) *testStruct {
		return &testStruct{value: "test", id: 1}
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

// TestSingletonInstance verifies that Instance() creates and returns an instance
func TestSingletonInstance(t *testing.T) {
	constructor := func(ctx context.Context) *testStruct {
		return &testStruct{value: "test", id: 42}
	}

	singleton := NewSingleton(constructor)
	instance := singleton.Instance()

	if instance == nil {
		t.Fatal("Instance() returned nil")
	}

	if instance.value != "test" {
		t.Errorf("Expected value 'test', got '%s'", instance.value)
	}

	if instance.id != 42 {
		t.Errorf("Expected id 42, got %d", instance.id)
	}
}

// TestSingletonSameInstance verifies that multiple calls return the same instance
func TestSingletonSameInstance(t *testing.T) {
	constructor := func(ctx context.Context) *testStruct {
		return &testStruct{value: "test", id: 100}
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

// TestConstructorCalledOnce verifies that the constructor is called exactly once
func TestConstructorCalledOnce(t *testing.T) {
	var callCount atomic.Int32

	constructor := func(ctx context.Context) *testStruct {
		callCount.Add(1)
		return &testStruct{value: "test"}
	}

	singleton := NewSingleton(constructor)

	for i := 0; i < 10; i++ {
		singleton.Instance()
	}

	if callCount.Load() != 1 {
		t.Errorf("Constructor was called %d times, expected 1", callCount.Load())
	}
}

// TestSingletonConcurrency verifies thread-safety with concurrent access
func TestSingletonConcurrency(t *testing.T) {
	var callCount atomic.Int32

	constructor := func(ctx context.Context) *testStruct {
		callCount.Add(1)
		return &testStruct{value: "concurrent", id: 999}
	}

	singleton := NewSingleton(constructor)

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	instances := make([]*testStruct, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			instances[index] = singleton.Instance()
		}(i)
	}

	wg.Wait()

	if callCount.Load() != 1 {
		t.Errorf("Constructor was called %d times, expected 1", callCount.Load())
	}

	firstInstance := instances[0]
	for i := 1; i < numGoroutines; i++ {
		if instances[i] != firstInstance {
			t.Errorf("Goroutine %d got a different instance", i)
		}
	}
}

// TestWithContext verifies that context is passed to the constructor
func TestWithContext(t *testing.T) {
	type ctxKey string
	const key ctxKey = "test-key"
	const expectedValue = "test-value"

	var receivedContext context.Context

	constructor := func(ctx context.Context) *testStruct {
		receivedContext = ctx
		return &testStruct{value: "test"}
	}

	ctx := context.WithValue(context.Background(), key, expectedValue)
	singleton := NewSingleton(constructor).WithContext(ctx)
	singleton.Instance()

	if receivedContext == nil {
		t.Fatal("Constructor did not receive context")
	}

	value := receivedContext.Value(key)
	if value != expectedValue {
		t.Errorf("Expected context value '%s', got '%v'", expectedValue, value)
	}
}

// TestWithContextNil verifies behavior when no context is set
func TestWithContextNil(t *testing.T) {
	var receivedContext context.Context

	constructor := func(ctx context.Context) *testStruct {
		receivedContext = ctx
		return &testStruct{value: "test"}
	}

	singleton := NewSingleton(constructor)
	singleton.Instance()

	if receivedContext != nil {
		t.Error("Expected nil context when WithContext is not called")
	}
}

// TestWithContextMethodChaining verifies that WithContext returns the singleton for chaining
func TestWithContextMethodChaining(t *testing.T) {
	constructor := func(ctx context.Context) *testStruct {
		return &testStruct{value: "test"}
	}

	singleton := NewSingleton(constructor)
	result := singleton.WithContext(context.Background())

	if result != singleton {
		t.Error("WithContext should return the same singleton instance for method chaining")
	}
}

// TestWithContextAfterInstance verifies that setting context after Instance() has no effect
func TestWithContextAfterInstance(t *testing.T) {
	type ctxKey string
	const key ctxKey = "test-key"

	var receivedContext context.Context

	constructor := func(ctx context.Context) *testStruct {
		receivedContext = ctx
		return &testStruct{value: "test"}
	}

	ctx1 := context.WithValue(context.Background(), key, "first")
	ctx2 := context.WithValue(context.Background(), key, "second")

	singleton := NewSingleton(constructor).WithContext(ctx1)
	singleton.Instance()

	singleton.WithContext(ctx2)

	value := receivedContext.Value(key)
	if value != "first" {
		t.Errorf("Expected context from first WithContext call, got '%v'", value)
	}
}

// TestMultipleTypes verifies that Singleton works with different types
func TestMultipleTypes(t *testing.T) {
	type stringContainer struct {
		data string
	}

	type intContainer struct {
		data int
	}

	stringConstructor := func(ctx context.Context) *stringContainer {
		return &stringContainer{data: "hello"}
	}

	intConstructor := func(ctx context.Context) *intContainer {
		return &intContainer{data: 123}
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
