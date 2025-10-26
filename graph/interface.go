// Package graph provides implementations for graph data structures and algorithms.
package graph

import (
	"github.com/0x626f/go-kit/abstract"
)

// BranchReceiver is a function type used for graph traversal algorithms.
// It's called for each edge during traversal and controls whether traversal should continue.
//
// Type parameters:
//   - V: The vertex type must implement Vertex[K]
//   - E: The edge type must implement Edge
//   - K: The key type for vertices, must be comparable
//
// Parameters:
//   - from: The source vertex of the current edge
//   - to: The target vertex of the current edge
//   - edge: The edge data connecting from and to
//
// Returns:
//   - true if traversal should stop after processing this edge
//   - false if traversal should continue
type BranchReceiver[V Vertex[K], E Edge, K comparable] func(from, to V, edge E) bool

// Listener defines an interface for objects that can react to graph events.
// This enables observer pattern implementations for graph modifications.
//
// Type parameters:
//   - V: The vertex type must implement abstract.Keyable[K]
//   - E: The edge type can be any type
//   - K: The key type for vertices must be comparable
type Listener[V abstract.Keyable[K], E any, K comparable] interface {
	// OnVertex is called when a vertex is added to or modified in the graph.
	// It provides a way to react to vertex-related events.
	//
	// Parameters:
	//   - v: The vertex that was added or modified
	OnVertex(V)

	// OnEdge is called when an edge is added to or modified in the graph.
	// It provides a way to react to edge-related events.
	//
	// Parameters:
	//   - e: The edge that was added or modified
	OnEdge(E)
}

// Vertex represents a node in a graph.
// It extends the abstract.Keyable interface to provide unique identification
// for each vertex in the graph.
//
// Type parameters:
//   - K: The key type, which must be comparable (support == and != operators)
type Vertex[K comparable] interface {
	// Embed the Keyable interface to provide unique identification
	abstract.Keyable[K]
}

// Edge represents a connection between two vertices in a graph.
// This is a marker interface that can be extended by specific edge implementations
// to include additional information like weights, directions, or other metadata.
type Edge interface{}
