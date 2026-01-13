package event

import "github.com/0x626f/go-kit/linkedlist"

// Pipeline represents a sequential event processing pipeline with support for
// success and error handling. Events are processed through a chain of handlers
// in the order they were added.
//
// The pipeline stops processing if any handler returns false or an error.
// If all handlers complete successfully, the optional success handler is called.
// Any errors encountered are passed to the optional error handler.
type Pipeline[Event any] struct {
	steps     *linkedlist.LinkedList[Handler[Event]]
	onSuccess Handler[Event]
	onError   ErrorHandler
}

// NewPipeline creates and returns a new Pipeline instance.
// The pipeline is initialized with an empty list of handlers.
func NewPipeline[Event any]() *Pipeline[Event] {
	return &Pipeline[Event]{
		steps: linkedlist.NewLinkedList[Handler[Event]](),
	}
}

// AddHandler adds a handler to the end of the processing pipeline.
// Handlers are executed in the order they are added.
// Returns the pipeline instance for method chaining.
func (pipeline *Pipeline[Event]) AddHandler(handler Handler[Event]) *Pipeline[Event] {
	pipeline.steps.Push(handler)
	return pipeline
}

// OnSuccess sets the success handler to be called after all pipeline handlers
// complete successfully. If a success handler is already set, subsequent calls
// are ignored (only the first success handler is used).
// Returns the pipeline instance for method chaining.
func (pipeline *Pipeline[Event]) OnSuccess(handler Handler[Event]) *Pipeline[Event] {
	if pipeline.onSuccess == nil {
		pipeline.onSuccess = handler
	}
	return pipeline
}

// OnError sets the error handler to be called when any handler in the pipeline
// returns an error. If an error handler is already set, subsequent calls are
// ignored (only the first error handler is used).
// Returns the pipeline instance for method chaining.
func (pipeline *Pipeline[Event]) OnError(handler ErrorHandler) *Pipeline[Event] {
	if pipeline.onError == nil {
		pipeline.onError = handler
	}
	return pipeline
}

// Process executes the event through the pipeline's handlers in sequence.
// Processing stops if any handler returns false or an error.
// If all handlers complete successfully, the success handler (if set) is called.
// Any errors encountered during processing are passed to the error handler (if set).
func (pipeline *Pipeline[Event]) Process(event Event) {
	for index := 0; index < pipeline.steps.Size(); index++ {
		handler := pipeline.steps.At(index)

		if handler == nil {
			return
		}

		if doNext, err := handler(event); err != nil || !doNext {
			if err != nil && pipeline.onError != nil {
				pipeline.onError(err)
			}
			return
		}
	}

	if pipeline.onSuccess != nil {
		if _, err := pipeline.onSuccess(event); err != nil && pipeline.onError != nil {
			pipeline.onError(err)

		}
	}
}
