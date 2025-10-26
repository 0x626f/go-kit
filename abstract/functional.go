package abstract

// Functor represents a generic function with no parameters and no return value.
// It can be used for simple callback operations.
type Functor func()

// Receiver is a generic function type that receives a value of type T and returns a boolean.
// The boolean return value typically indicates whether processing should continue (true)
// or stop (false) when used in iteration contexts.
type Receiver[T any] func(T) bool

// IndexedReceiver is a generic function type that receives an index of type I and a value of type T,
// and returns a boolean. Similar to Receiver, the boolean return value typically indicates
// whether processing should continue (true) or stop (false).
// This is commonly used in collection iteration where both index and value are needed.
type IndexedReceiver[I any, T any] func(I, T) bool

// Predicate is a generic function type that evaluates a value of type T and returns a boolean.
// It's typically used for filtering, where the boolean indicates whether the value
// satisfies a condition (true) or not (false).
type Predicate[T any] func(T) bool

// Transformer is a generic function type that transforms a pointer to a value of type T
// into a value of type R. This allows for both reading and potentially modifying the
// original value during the transformation process.
type Transformer[T, R any] func(*T) R
