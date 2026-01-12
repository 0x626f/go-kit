package event

import (
	"errors"
	"testing"
)

type mockEvent struct {
	root string
}

func (m *mockEvent) Root() string {
	return m.root
}

func TestNewPipeline(t *testing.T) {
	pipeline := NewPipeline[string]()
	if pipeline == nil {
		t.Fatal("NewPipeline returned nil")
	}
	if pipeline.steps == nil {
		t.Error("pipeline.steps should not be nil")
	}
}

func TestPipeline_AddHandler(t *testing.T) {
	pipeline := NewPipeline[string]()
	handler := func(event Event[string]) (bool, error) {
		return true, nil
	}

	result := pipeline.AddHandler(handler)
	if result != pipeline {
		t.Error("AddHandler should return the pipeline for chaining")
	}
	if pipeline.steps.Size() != 1 {
		t.Errorf("expected 1 handler, got %d", pipeline.steps.Size())
	}
}

func TestPipeline_OnSuccess(t *testing.T) {
	pipeline := NewPipeline[string]()
	successHandler := func(event Event[string]) (bool, error) {
		return true, nil
	}

	result := pipeline.OnSuccess(successHandler)
	if result != pipeline {
		t.Error("OnSuccess should return the pipeline for chaining")
	}
	if pipeline.onSuccess == nil {
		t.Error("onSuccess handler should be set")
	}
}

func TestPipeline_OnSuccess_OnlyFirstSet(t *testing.T) {
	pipeline := NewPipeline[string]()
	called := 0
	firstHandler := func(event Event[string]) (bool, error) {
		called = 1
		return true, nil
	}
	secondHandler := func(event Event[string]) (bool, error) {
		called = 2
		return true, nil
	}

	pipeline.OnSuccess(firstHandler)
	pipeline.OnSuccess(secondHandler)

	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if called != 1 {
		t.Errorf("expected first handler to be called (called=%d), second handler should be ignored", called)
	}
}

func TestPipeline_OnError(t *testing.T) {
	pipeline := NewPipeline[string]()
	errorHandler := func(err error) {}

	result := pipeline.OnError(errorHandler)
	if result != pipeline {
		t.Error("OnError should return the pipeline for chaining")
	}
	if pipeline.onError == nil {
		t.Error("onError handler should be set")
	}
}

func TestPipeline_OnError_OnlyFirstSet(t *testing.T) {
	pipeline := NewPipeline[string]()
	called := 0
	firstHandler := func(err error) {
		called = 1
	}
	secondHandler := func(err error) {
		called = 2
	}

	pipeline.OnError(firstHandler)
	pipeline.OnError(secondHandler)

	handler := func(event Event[string]) (bool, error) {
		return false, errors.New("test error")
	}
	pipeline.AddHandler(handler)

	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if called != 1 {
		t.Errorf("expected first error handler to be called (called=%d), second handler should be ignored", called)
	}
}

func TestPipeline_Process_SingleHandler(t *testing.T) {
	pipeline := NewPipeline[string]()
	called := false
	handler := func(event Event[string]) (bool, error) {
		called = true
		return true, nil
	}

	pipeline.AddHandler(handler)
	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if !called {
		t.Error("handler should have been called")
	}
}

func TestPipeline_Process_MultipleHandlers(t *testing.T) {
	pipeline := NewPipeline[string]()
	order := []int{}

	handler1 := func(event Event[string]) (bool, error) {
		order = append(order, 1)
		return true, nil
	}
	handler2 := func(event Event[string]) (bool, error) {
		order = append(order, 2)
		return true, nil
	}
	handler3 := func(event Event[string]) (bool, error) {
		order = append(order, 3)
		return true, nil
	}

	pipeline.AddHandler(handler1).AddHandler(handler2).AddHandler(handler3)
	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if len(order) != 3 {
		t.Fatalf("expected 3 handlers to be called, got %d", len(order))
	}
	if order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("handlers called in wrong order: %v", order)
	}
}

func TestPipeline_Process_HandlerReturnsFalse(t *testing.T) {
	pipeline := NewPipeline[string]()
	order := []int{}

	handler1 := func(event Event[string]) (bool, error) {
		order = append(order, 1)
		return true, nil
	}
	handler2 := func(event Event[string]) (bool, error) {
		order = append(order, 2)
		return false, nil
	}
	handler3 := func(event Event[string]) (bool, error) {
		order = append(order, 3)
		return true, nil
	}

	pipeline.AddHandler(handler1).AddHandler(handler2).AddHandler(handler3)
	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if len(order) != 2 {
		t.Fatalf("expected 2 handlers to be called before stopping, got %d", len(order))
	}
	if order[0] != 1 || order[1] != 2 {
		t.Errorf("wrong handlers called: %v", order)
	}
}

func TestPipeline_Process_HandlerReturnsError(t *testing.T) {
	pipeline := NewPipeline[string]()
	order := []int{}
	var capturedError error

	handler1 := func(event Event[string]) (bool, error) {
		order = append(order, 1)
		return true, nil
	}
	handler2 := func(event Event[string]) (bool, error) {
		order = append(order, 2)
		return false, errors.New("test error")
	}
	handler3 := func(event Event[string]) (bool, error) {
		order = append(order, 3)
		return true, nil
	}

	pipeline.AddHandler(handler1).AddHandler(handler2).AddHandler(handler3)
	pipeline.OnError(func(err error) {
		capturedError = err
	})

	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if len(order) != 2 {
		t.Fatalf("expected 2 handlers to be called before error, got %d", len(order))
	}
	if capturedError == nil {
		t.Error("error handler should have been called")
	}
	if capturedError.Error() != "test error" {
		t.Errorf("expected 'test error', got '%s'", capturedError.Error())
	}
}

func TestPipeline_Process_OnSuccessCalled(t *testing.T) {
	pipeline := NewPipeline[string]()
	successCalled := false

	handler := func(event Event[string]) (bool, error) {
		return true, nil
	}
	successHandler := func(event Event[string]) (bool, error) {
		successCalled = true
		return true, nil
	}

	pipeline.AddHandler(handler).OnSuccess(successHandler)
	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if !successCalled {
		t.Error("success handler should have been called")
	}
}

func TestPipeline_Process_OnSuccessNotCalledWhenHandlerFails(t *testing.T) {
	pipeline := NewPipeline[string]()
	successCalled := false

	handler := func(event Event[string]) (bool, error) {
		return false, errors.New("test error")
	}
	successHandler := func(event Event[string]) (bool, error) {
		successCalled = true
		return true, nil
	}

	pipeline.AddHandler(handler).OnSuccess(successHandler)
	pipeline.OnError(func(err error) {})

	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if successCalled {
		t.Error("success handler should not have been called when handler fails")
	}
}

func TestPipeline_Process_OnSuccessReturnsError(t *testing.T) {
	pipeline := NewPipeline[string]()
	var capturedError error

	handler := func(event Event[string]) (bool, error) {
		return true, nil
	}
	successHandler := func(event Event[string]) (bool, error) {
		return true, errors.New("success error")
	}

	pipeline.AddHandler(handler).OnSuccess(successHandler)
	pipeline.OnError(func(err error) {
		capturedError = err
	})

	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if capturedError == nil {
		t.Error("error handler should have been called for success handler error")
	}
	if capturedError.Error() != "success error" {
		t.Errorf("expected 'success error', got '%s'", capturedError.Error())
	}
}

func TestPipeline_Process_NoOnErrorHandler(t *testing.T) {
	pipeline := NewPipeline[string]()

	handler := func(event Event[string]) (bool, error) {
		return false, errors.New("test error")
	}

	pipeline.AddHandler(handler)
	event := &mockEvent{root: "test"}
	pipeline.Process(event)
}

func TestPipeline_Process_EmptyPipeline(t *testing.T) {
	pipeline := NewPipeline[string]()
	successCalled := false

	successHandler := func(event Event[string]) (bool, error) {
		successCalled = true
		return true, nil
	}

	pipeline.OnSuccess(successHandler)
	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if !successCalled {
		t.Error("success handler should be called even with no handlers")
	}
}

func TestPipeline_Chaining(t *testing.T) {
	successCalled := false
	errorHandled := false

	pipeline := NewPipeline[string]().
		AddHandler(func(event Event[string]) (bool, error) {
			return true, nil
		}).
		AddHandler(func(event Event[string]) (bool, error) {
			return true, nil
		}).
		OnSuccess(func(event Event[string]) (bool, error) {
			successCalled = true
			return true, nil
		}).
		OnError(func(err error) {
			errorHandled = true
		})

	event := &mockEvent{root: "test"}
	pipeline.Process(event)

	if !successCalled {
		t.Error("success handler should have been called")
	}
	if errorHandled {
		t.Error("error handler should not have been called")
	}
}
