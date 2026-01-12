// Package event provides a flexible event processing framework with pipeline and routing capabilities.
package event

// Event represents a generic event with a root value of type T.
// The root value typically contains the core data or context of the event.
type Event[T any] interface {
	// Root returns the root value of the event.
	Root() T
}

// Handler is a function that processes an event and returns whether to continue processing
// and an optional error. If the boolean return value is false or an error is returned,
// the pipeline processing will stop.
type Handler[Root any] func(Event[Root]) (bool, error)

// ErrorHandler is a function that handles errors that occur during event processing.
// It receives an error and can perform logging, monitoring, or other error handling tasks.
type ErrorHandler func(error)

// Receiver is a function that receives and processes an event without returning any value.
// It is typically used as a terminal handler in routing scenarios where events are
// dispatched to specific receivers based on routing logic.
type Receiver[Root any] func(Event[Root])

// Resolver is a function that examines an event and returns an ID used for routing decisions.
// The ID type must be comparable and is used to determine which receiver should handle the event.
type Resolver[Root any, ID comparable] func(Event[Root]) ID
