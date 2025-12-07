package graph

import "github.com/0x626f/go-kit/abstract"

// Route represents a path through a graph as an ordered sequence of vertices.
// A route can be used to store the result of path-finding algorithms or to
// define traversal patterns through the graph.
//
// Type parameters:
//   - V: The vertex type must implement abstract.Keyable[K]
//   - K: The key type for vertices must be comparable
type Route[V abstract.Keyable[K], K comparable] struct {
	// Vertices is an ordered slice containing all vertices in the path,
	// starting with the source vertex and ending with the destination vertex.
	Vertices []V
}
