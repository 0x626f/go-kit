// Package graph provides implementations for graph data structures and algorithms.
package graph

import (
	"errors"
	"fmt"
	"github.com/0x626f/go-kit/pkg/abstract"
	"github.com/0x626f/go-kit/pkg/utils"
)

// Feature is a bitwise flag type used to configure graph properties.
// Multiple features can be combined using bitwise OR operations.
type Feature uint8

const (
	// None represents a graph with no special features.
	None Feature = 0

	// Directed indicates a graph where edges have a direction (from one vertex to another).
	// In a directed graph, an edge from A to B does not imply an edge from B to A.
	Directed Feature = 1 << iota

	// Acyclic indicates a graph that doesn't allow cycles.
	// Operations that would create a cycle will fail with an error.
	Acyclic
)

// Graph implements a generic graph data structure.
// It supports both directed and undirected graphs, with optional acyclic constraints.
//
// Type parameters:
//   - V: The vertex type must implement Vertex[K]
//   - E: The edge type must implement Edge
//   - K: The key type for vertices must be comparable
type Graph[V Vertex[K], E Edge, K comparable] struct {
	// feature stores the configured features for this graph instance
	feature Feature
	// vertices maps vertex keys to vertex objects
	vertices map[K]V
	// matrix stores the graph structure using an adjacency matrix
	matrix *AdjacencyMatrix[V, E, K]
}

// ForFeature creates a new Graph instance with the specified features.
// This is the recommended way to create a Graph with specific behavior.
//
// Parameters:
//   - feature: Bitwise combination of Feature flags that define graph behavior
//
// Returns a new, empty Graph configured with the specified features.
func ForFeature[V Vertex[K], E Edge, K comparable](feature Feature) *Graph[V, E, K] {
	return &Graph[V, E, K]{
		feature:  feature,
		vertices: make(map[K]V),
		matrix:   NewAdjacencyMatrix[V, E, K](feature&Directed != 0),
	}
}

// HasFeature checks if the graph has a specific feature enabled.
//
// Parameters:
//   - feature: The feature flag to check for
//
// Returns true if the graph has the specified feature, false otherwise.
func (graph *Graph[V, E, K]) HasFeature(feature Feature) bool {
	return graph.feature&feature != 0
}

// AddVertex adds a new vertex to the graph.
//
// Parameters:
//   - data: The vertex to add
//
// Returns an error if:
//   - A vertex with the same key already exists in the graph
func (graph *Graph[V, E, K]) AddVertex(data V) error {
	key := data.Key()

	if _, exist := graph.vertices[key]; exist {
		return errors.New("vertex already exists")
	}

	graph.vertices[key] = data
	return nil
}

// AddEdge creates an edge between the vertices with keys from and to.
// The behavior depends on the graph's features:
//   - In a directed graph, creates a one-way connection from -> to
//   - In an undirected graph, creates a two-way connection between from and to
//   - In an acyclic graph, prevents additions that would create cycles
//
// Parameters:
//   - from: The key of the source vertex
//   - to: The key of the target vertex
//   - data: The edge data
//
// Returns an error if:
//   - Either vertex doesn't exist
//   - The graph has the Acyclic feature and the edge would create a cycle
//   - The underlying adjacency matrix rejects the edge (e.g., due to directedness constraints)
func (graph *Graph[V, E, K]) AddEdge(from, to K, data E) error {
	_, existFrom := graph.vertices[from]
	_, existTo := graph.vertices[to]

	if !existFrom {
		return fmt.Errorf("vertex %v is not present", from)
	}

	if !existTo {
		return fmt.Errorf("vertex %v is not present", to)
	}

	// Prevent self-loops in acyclic graphs
	if graph.HasFeature(Acyclic) && from == to {
		return fmt.Errorf("connection (%v,%v) makes loop", from, to)
	}

	// Do nothing if the edge already exists
	if graph.matrix.Has(from, to) {
		return nil
	}

	// Try to add the edge
	err := graph.matrix.Add(from, to, data)
	if err != nil {
		return err
	}

	// If the graph is acyclic, check for cycles after adding the edge
	if graph.HasFeature(Acyclic) && graph.matrix.HasCycles() {
		// Remove the edge that created the cycle
		err = graph.matrix.Remove(from, to)
		if err != nil {
			return err
		}
		return fmt.Errorf("connection (%v,%v) makes cycles", from, to)
	}

	return nil
}

// Vertex retrieves a vertex from the graph by its key.
//
// Parameters:
//   - key: The key of the vertex to retrieve
//
// Returns:
//   - The vertex object if found
//   - A boolean indicating whether the vertex was found
func (graph *Graph[V, E, K]) Vertex(key K) (V, bool) {
	if vertex, exists := graph.vertices[key]; exists {
		return vertex, exists
	}

	return utils.Zero[V](), false
}

// Edge retrieves an edge between two vertices by their keys.
//
// Parameters:
//   - from: The key of the source vertex
//   - to: The key of the target vertex
//
// Returns:
//   - The edge object if found
//   - A boolean indicating whether the edge was found
func (graph *Graph[V, E, K]) Edge(from, to K) (E, bool) {
	if !graph.matrix.Has(from, to) {
		return utils.Zero[E](), false
	}

	return graph.matrix.data[from][to], true
}

// Vertices return a slice containing all vertices in the graph.
//
// Returns a new slice containing all vertex objects.
func (graph *Graph[V, E, K]) Vertices() []V {
	return utils.MapToValueSlice(graph.vertices)
}

// HasConnection checks if there is an edge between two vertices.
//
// Parameters:
//   - from: The key of the source vertex
//   - to: The key of the target vertex
//
// Returns true if an edge exists from the source to the target vertex, false otherwise.
func (graph *Graph[V, E, K]) HasConnection(from, to K) bool {
	return graph.matrix.Has(from, to)
}

// DFS performs a depth-first search traversal of the graph starting from a specified vertex.
// For each edge visited, it calls the provided receiver function.
//
// Parameters:
//   - from: The key of the starting vertex for traversal
//   - receiver: A function called for each edge during traversal. If this function returns true,
//     the traversal stops early.
func (graph *Graph[V, E, K]) DFS(from K, receiver BranchReceiver[V, E, K]) {
	vertex, exists := graph.vertices[from]
	done := false
	if !exists {
		return
	}

	visited := make(map[K]struct{})
	visited[from] = struct{}{}

	// Inner recursive DFS function
	var dfs func(from K, receiver BranchReceiver[V, E, K])
	dfs = func(from K, receiver BranchReceiver[V, E, K]) {
		visited[from] = struct{}{}
		vertex = graph.vertices[from]

		for subVertex, edge := range graph.matrix.data[from] {
			done = receiver(graph.vertices[from], graph.vertices[subVertex], edge)

			if done {
				goto dfsEnd
			}

			if _, beenVisited := visited[subVertex]; !beenVisited {
				dfs(subVertex, receiver)
				if done {
					goto dfsEnd
				}
			}
		}
	dfsEnd:
		return
	}

	// Start DFS from the initial vertex
	for subVertex, edge := range graph.matrix.data[from] {
		done = receiver(vertex, graph.vertices[subVertex], edge)

		if done {
			goto End
		}

		dfs(subVertex, receiver)

		if done {
			goto End
		}
	}
End:
	return
}

// BFS performs a breadth-first search traversal of the graph starting from a specified vertex.
// For each edge visited, it calls the provided receiver function.
//
// Parameters:
//   - from: The key of the starting vertex for traversal
//   - receiver: A function called for each edge during traversal. If this function returns true,
//     the traversal stops early.
func (graph *Graph[V, E, K]) BFS(from K, receiver BranchReceiver[V, E, K]) {
	vertex, exists := graph.vertices[from]

	if !exists {
		return
	}

	visited := make(map[K]struct{})
	queue := []K{from}
	visited[vertex.Key()] = struct{}{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		vertex = graph.vertices[node]

		for subVertex, edge := range graph.matrix.data[node] {
			if receiver(vertex, graph.vertices[subVertex], edge) {
				goto End
			}

			if _, beenVisited := visited[subVertex]; !beenVisited {
				visited[subVertex] = struct{}{}
				queue = append(queue, subVertex)
			}
		}
	}
End:
	return
}

// Routes finds all paths between two vertices that satisfy certain conditions.
//
// Parameters:
//   - from: The key of the starting vertex
//   - to: The key of the destination vertex
//   - predicate: An optional function that filters vertices to include in paths.
//     If nil, all vertices are considered valid.
//   - length: The maximum path length to consider. If 0, all path lengths are considered.
//
// Returns a slice of Route objects, each representing a valid path from the start to the destination.
// Returns nil if no valid paths exist.
func (graph *Graph[V, E, K]) Routes(from, to K, predicate abstract.Predicate[V], length int) []*Route[V, K] {
	// No path from a vertex to itself
	if from == to {
		return nil
	}

	fromVertex := graph.vertices[from]
	toVertex := graph.vertices[to]

	// Check if start and end vertices satisfy the predicate
	if predicate != nil && (!predicate(fromVertex) || !predicate(toVertex)) {
		return nil
	}

	visited := make(map[K]struct{})
	var routes []*Route[V, K]
	var route []V

	// Inner recursive DFS function for path finding
	var dfs func(key K, depth int)
	dfs = func(key K, depth int) {
		vertex := graph.vertices[key]
		if predicate != nil && !predicate(vertex) {
			return
		}

		visited[key] = struct{}{}
		route = append(route, vertex)

		// Found a path to the destination
		if key == to && (length == 0 || depth <= length) {
			routes = append(routes, &Route[V, K]{
				Vertices: append([]V{}, route...),
			})
		} else if length == 0 || depth < length {
			// Continue searching for paths
			for subKey, _ := range graph.matrix.data[key] {
				subVertex := graph.vertices[subKey]
				_, beenVisited := visited[subKey]
				if !beenVisited && (predicate == nil || predicate(subVertex)) {
					dfs(subKey, depth+1)
				}
			}
		}
		// Backtrack
		route = route[:len(route)-1]
		delete(visited, key)
	}

	dfs(from, 1)
	return routes
}
