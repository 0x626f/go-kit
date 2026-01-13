// Package event provides a flexible event processing framework with pipeline and routing capabilities.
package event

// Handler is a function that processes an event and returns whether to continue processing
// and an optional error. If the boolean return value is false or an error is returned,
// the pipeline processing will stop.
type Handler[Event any] func(Event) (bool, error)

// ErrorHandler is a function that handles errors that occur during event processing.
// It receives an error and can perform logging, monitoring, or other error handling tasks.
type ErrorHandler func(error)

// Receiver is a function that receives and processes an event without returning any value.
// It is typically used as a terminal handler in routing scenarios where events are
// dispatched to specific receivers based on routing logic.
type Receiver[Event any] func(Event)

// Resolver is a function that examines an event and returns an ID used for routing decisions.
// The ID type must be comparable and is used to determine which receiver should handle the event.
type Resolver[Event any, ID comparable] func(Event) ID
