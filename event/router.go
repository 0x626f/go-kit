package event

import "fmt"

// Router is a generic event router that dispatches events to specific receivers
// based on a resolution strategy. The router uses a Resolver function to determine
// which receiver should handle each event.
//
// The ID type parameter must be comparable and is used as a key to map events
// to their corresponding receivers.
type Router[Root any, ID comparable] struct {
	routing  map[ID]Receiver[Root]
	resolver Resolver[Root, ID]
}

// NewRouter creates and returns a new Router instance with an empty routing map.
func NewRouter[Root any, ID comparable]() *Router[Root, ID] {
	return &Router[Root, ID]{
		routing: make(map[ID]Receiver[Root]),
	}
}

// AddReceiver registers a receiver function for the given ID.
// When an event is routed and the resolver returns this ID, the corresponding
// receiver will be called. If a receiver already exists for the ID, it will be replaced.
// Returns the router instance for method chaining.
func (router *Router[Root, ID]) AddReceiver(id ID, receiver Receiver[Root]) *Router[Root, ID] {
	router.routing[id] = receiver
	return router
}

// SetResolver sets the resolver function used to determine which receiver should
// handle each event. The resolver examines the event and returns an ID that maps
// to a registered receiver.
//
// Note: Due to a bug in the implementation, this method only sets the resolver
// when the provided resolver parameter is nil. To set a non-nil resolver, assign
// it directly to router.resolver. If a resolver is already set, subsequent calls
// are ignored (only the first resolver is used).
// Returns the router instance for method chaining.
func (router *Router[Root, ID]) SetResolver(resolver Resolver[Root, ID]) *Router[Root, ID] {
	if resolver == nil {
		router.resolver = resolver
	}
	return router
}

// Route processes an event by using the resolver to determine the target receiver
// and then dispatching the event to that receiver.
//
// Returns an error if no resolver is configured. If the resolver returns an ID
// that has no registered receiver, the event is silently ignored (no error is returned).
func (router *Router[Root, ID]) Route(event Event[Root]) error {
	if router.resolver == nil {
		return fmt.Errorf("missing event resolver")
	}

	id := router.resolver(event)
	receiver := router.routing[id]

	if receiver != nil {
		receiver(event)
	}
	return nil
}
