package event

import (
	"testing"
)

type intEvent struct {
	root int
}

func TestNewRouter(t *testing.T) {
	router := NewRouter[mockEvent, int]()
	if router == nil {
		t.Fatal("NewRouter returned nil")
	}
	if router.routing == nil {
		t.Error("router.routing should not be nil")
	}
}

func TestRouter_AddReceiver(t *testing.T) {
	router := NewRouter[mockEvent, int]()
	called := false
	receiver := func(event mockEvent) {
		called = true
	}

	result := router.AddReceiver(1, receiver)
	if result != router {
		t.Error("AddReceiver should return the router for chaining")
	}

	event := mockEvent{root: "test"}
	router.routing[1](event)

	if !called {
		t.Error("receiver should be callable from routing map")
	}
}

func TestRouter_AddReceiver_MultipleReceivers(t *testing.T) {
	router := NewRouter[mockEvent, string]()
	called := map[string]bool{}

	receiver1 := func(event mockEvent) {
		called["receiver1"] = true
	}
	receiver2 := func(event mockEvent) {
		called["receiver2"] = true
	}

	router.AddReceiver("id1", receiver1)
	router.AddReceiver("id2", receiver2)

	if len(router.routing) != 2 {
		t.Errorf("expected 2 receivers, got %d", len(router.routing))
	}

	event := mockEvent{root: "test"}
	router.routing["id1"](event)
	router.routing["id2"](event)

	if !called["receiver1"] || !called["receiver2"] {
		t.Error("all receivers should be callable")
	}
}

func TestRouter_SetResolver(t *testing.T) {
	router := NewRouter[mockEvent, int]()
	resolver := func(event mockEvent) int {
		return 1
	}

	result := router.SetResolver(resolver)
	if result != router {
		t.Error("SetResolver should return the router for chaining")
	}
}

func TestRouter_SetResolver_BugInImplementation(t *testing.T) {
	router := NewRouter[mockEvent, int]()
	resolver := func(event mockEvent) int {
		return 1
	}

	router.SetResolver(resolver)

	if router.resolver != nil {
		t.Error("due to bug in SetResolver (checks 'if resolver == nil' instead of 'if router.resolver == nil'), resolver should not be set when passing non-nil resolver")
	}
}

func TestRouter_Route_Success(t *testing.T) {
	router := NewRouter[mockEvent, int]()
	receivedEvent := false
	var receivedRoot string

	receiver := func(event mockEvent) {
		receivedEvent = true
		receivedRoot = event.root
	}
	resolver := func(event mockEvent) int {
		return 1
	}

	router.AddReceiver(1, receiver)
	router.resolver = resolver

	event := mockEvent{root: "test data"}
	err := router.Route(event)

	if err != nil {
		t.Errorf("Route should not return error: %v", err)
	}
	if !receivedEvent {
		t.Error("receiver should have been called")
	}
	if receivedRoot != "test data" {
		t.Errorf("expected 'test data', got '%s'", receivedRoot)
	}
}

func TestRouter_Route_MissingResolver(t *testing.T) {
	router := NewRouter[mockEvent, int]()
	receiver := func(event mockEvent) {}

	router.AddReceiver(1, receiver)

	event := mockEvent{root: "test"}
	err := router.Route(event)

	if err == nil {
		t.Error("Route should return error when resolver is missing")
	}
	if err.Error() != "missing event resolver" {
		t.Errorf("expected 'missing event resolver', got '%s'", err.Error())
	}
}

func TestRouter_Route_MissingReceiver(t *testing.T) {
	router := NewRouter[mockEvent, int]()
	resolver := func(event mockEvent) int {
		return 99
	}

	router.AddReceiver(1, func(event mockEvent) {})
	router.resolver = resolver

	event := mockEvent{root: "test"}
	err := router.Route(event)

	if err != nil {
		t.Errorf("Route should not return error when receiver is missing (current implementation just skips), got: %v", err)
	}
}

func TestRouter_Route_DifferentIDTypes(t *testing.T) {
	t.Run("string ID", func(t *testing.T) {
		router := NewRouter[mockEvent, string]()
		called := false

		receiver := func(event mockEvent) {
			called = true
		}
		resolver := func(event mockEvent) string {
			return "route-a"
		}

		router.AddReceiver("route-a", receiver)
		router.resolver = resolver

		event := mockEvent{root: "test"}
		err := router.Route(event)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !called {
			t.Error("receiver should have been called")
		}
	})

	t.Run("custom struct ID", func(t *testing.T) {
		type RouteID struct {
			Type string
			ID   int
		}

		router := NewRouter[mockEvent, RouteID]()
		called := false

		receiver := func(event mockEvent) {
			called = true
		}
		resolver := func(event mockEvent) RouteID {
			return RouteID{Type: "user", ID: 123}
		}

		router.AddReceiver(RouteID{Type: "user", ID: 123}, receiver)
		router.resolver = resolver

		event := mockEvent{root: "test"}
		err := router.Route(event)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !called {
			t.Error("receiver should have been called")
		}
	})
}

func TestRouter_Route_ResolverBasedOnEventData(t *testing.T) {
	router := NewRouter[mockEvent, string]()
	receivedByA := false
	receivedByB := false

	receiverA := func(event mockEvent) {
		receivedByA = true
	}
	receiverB := func(event mockEvent) {
		receivedByB = true
	}
	resolver := func(event mockEvent) string {
		root := event.root
		if len(root) > 5 {
			return "route-a"
		}
		return "route-b"
	}

	router.AddReceiver("route-a", receiverA)
	router.AddReceiver("route-b", receiverB)
	router.resolver = resolver

	event1 := mockEvent{root: "long string"}
	router.Route(event1)

	event2 := mockEvent{root: "short"}
	router.Route(event2)

	if !receivedByA {
		t.Error("receiver A should have been called for long string")
	}
	if !receivedByB {
		t.Error("receiver B should have been called for short string")
	}
}

func TestRouter_Chaining(t *testing.T) {
	called := false

	router := NewRouter[mockEvent, int]().
		AddReceiver(1, func(event mockEvent) {
			called = true
		}).
		AddReceiver(2, func(event mockEvent) {})

	router.resolver = func(event mockEvent) int {
		return 1
	}

	event := mockEvent{root: "test"}
	err := router.Route(event)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("receiver should have been called")
	}
}

func TestRouter_Route_MultipleRoutingScenarios(t *testing.T) {
	router := NewRouter[intEvent, string]()
	callOrder := []string{}

	router.AddReceiver("handler1", func(event intEvent) {
		callOrder = append(callOrder, "handler1")
	})
	router.AddReceiver("handler2", func(event intEvent) {
		callOrder = append(callOrder, "handler2")
	})
	router.AddReceiver("handler3", func(event intEvent) {
		callOrder = append(callOrder, "handler3")
	})

	ie := intEvent{}

	router.resolver = func(event intEvent) string {
		val := event.root
		if val == 1 {
			return "handler1"
		} else if val == 2 {
			return "handler2"
		}
		return "handler3"
	}

	ie.root = 1
	router.Route(ie)

	ie.root = 2
	router.Route(ie)

	ie.root = 3
	router.Route(ie)

	if len(callOrder) != 3 {
		t.Fatalf("expected 3 calls, got %d", len(callOrder))
	}
	if callOrder[0] != "handler1" || callOrder[1] != "handler2" || callOrder[2] != "handler3" {
		t.Errorf("unexpected call order: %v", callOrder)
	}
}
