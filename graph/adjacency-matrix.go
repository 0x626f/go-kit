package graph

import (
	"fmt"
)

// AdjacencyMatrix represents a graph using an adjacency matrix representation.
// This implementation uses a map of maps to store connections between vertices.
//
// Type parameters:
//   - V: The vertex type, must implement Vertex[K]
//   - E: The edge type, must implement Edge
//   - K: The key type for vertices, must be comparable
type AdjacencyMatrix[V Vertex[K], E Edge, K comparable] struct {
	// directed indicates whether the graph is directed (true) or undirected (false)
	directed bool
	// data stores the adjacency matrix as a map of maps
	// The outer map keys are source vertex keys, the inner map keys are target vertex keys,
	// and the values are the edges connecting them
	data map[K]map[K]E
}

// NewAdjacencyMatrix creates a new adjacency matrix instance.
//
// Parameters:
//   - directed: If true, creates a directed graph; if false, creates an undirected graph
//
// Returns a new, empty adjacency matrix.
func NewAdjacencyMatrix[V Vertex[K], E Edge, K comparable](directed bool) *AdjacencyMatrix[V, E, K] {
	return &AdjacencyMatrix[V, E, K]{
		directed: directed,
		data:     make(map[K]map[K]E),
	}
}

// assertRow ensures that the row for the specified vertex key exists in the adjacency matrix.
// If the row doesn't exist, it creates a new empty row.
//
// Parameters:
//   - i: The key of the vertex for which to ensure a row exists
func (matrix *AdjacencyMatrix[V, E, K]) assertRow(i K) {
	_, exists := matrix.data[i]

	if !exists {
		matrix.data[i] = make(map[K]E)
	}
}

// Add creates an edge between vertices with keys i and j, with the specified edge data.
// In an undirected graph, this adds connections in both directions.
// In a directed graph, this only adds a connection from i to j.
//
// Parameters:
//   - i: The key of the source vertex
//   - j: The key of the target vertex
//   - data: The edge data to associate with this connection
//
// Returns an error if:
//   - In a directed graph, an edge from j to i already exists (to prevent cycles)
func (matrix *AdjacencyMatrix[V, E, K]) Add(i, j K, data E) error {
	if matrix.directed && matrix.Has(j, i) {
		return fmt.Errorf("couldn't add connection (%v,%v) as connection (%v,%v) is present", i, j, j, i)
	}

	matrix.assertRow(i)
	matrix.data[i][j] = data

	if !matrix.directed {
		matrix.assertRow(j)
		matrix.data[j][i] = data
	}

	return nil
}

// Remove deletes the edge between vertices with keys from and to.
// In an undirected graph, this removes connections in both directions.
// In a directed graph, this only removes the connection from the source to the target.
//
// Parameters:
//   - from: The key of the source vertex
//   - to: The key of the target vertex
//
// Returns an error if:
//   - The specified connection doesn't exist
func (matrix *AdjacencyMatrix[V, E, K]) Remove(from, to K) error {
	if !matrix.Has(from, to) {
		return fmt.Errorf("connection (%v,%v) doesn't exist", from, to)
	}

	row := matrix.data[from]
	delete(row, to)

	if len(row) == 0 {
		delete(matrix.data, from)
	}

	if !matrix.directed && matrix.Has(to, from) {
		row = matrix.data[to]
		delete(row, from)

		if len(row) == 0 {
			delete(matrix.data, to)
		}
	}

	return nil
}

// Has checks if there is an edge between vertices with keys i and j.
//
// Parameters:
//   - i: The key of the source vertex
//   - j: The key of the target vertex
//
// Returns:
//   - true if an edge exists between the vertices
//   - false if no edge exists or if either vertex doesn't exist
func (matrix *AdjacencyMatrix[V, E, K]) Has(i, j K) bool {
	row, ok := matrix.data[i]
	if !ok {
		return false
	}
	_, ok = row[j]
	return ok
}

// HasCycles checks if the graph contains any cycles.
// A cycle is a path that starts and ends at the same vertex.
//
// The implementation uses two approaches:
// 1. First checks for self-loops (edges from a vertex to itself)
// 2. Then performs a depth-first search to find other cycles
//
// Returns:
//   - true if the graph contains at least one cycle
//   - false if the graph is acyclic
func (matrix *AdjacencyMatrix[V, E, K]) HasCycles() bool {
	// Quick check for self-loops first
	for vertex, links := range matrix.data {
		if _, ok := links[vertex]; ok {
			return true
		}
	}

	var visited, records map[K]struct{}
	var zero K

	visited = make(map[K]struct{})
	if matrix.directed {
		records = make(map[K]struct{})
	}

	// Check each vertex as a potential start of a cycle
	for vertex := range matrix.data {
		if _, beenVisited := visited[vertex]; !beenVisited {
			if matrix.cycled(vertex, zero, visited, records, false) {
				return true
			}
		}
	}
	return false
}

// cycled is a helper method for HasCycles that performs a depth-first search
// to detect cycles in the graph.
//
// Parameters:
//   - vertex: The current vertex being explored
//   - parent: The parent vertex in the DFS tree (for undirected graphs)
//   - visited: A map of vertices that have been visited in the current DFS
//   - records: A map of vertices in the current DFS path (for directed graphs)
//   - hasParent: Whether the current vertex has a parent (for undirected graphs)
//
// Returns:
//   - true if a cycle is detected starting from the current vertex
//   - false if no cycle is detected
func (matrix *AdjacencyMatrix[V, E, K]) cycled(vertex, parent K, visited, records map[K]struct{}, hasParent bool) bool {
	visited[vertex] = struct{}{}

	if matrix.directed {
		records[vertex] = struct{}{}
	}

	for subVertex := range matrix.data[vertex] {
		_, beenVisited := visited[subVertex]

		if !beenVisited {
			// Recurse to unvisited neighbors
			if matrix.cycled(subVertex, vertex, visited, records, true) {
				return true
			}
		} else if matrix.directed {
			// For directed graphs, check if the neighbor is in the current path
			if _, beenRecorded := records[subVertex]; beenRecorded {
				return true
			}
		} else if hasParent && subVertex != parent {
			// For undirected graphs, a cycle exists if we find an edge to a visited vertex
			// that is not the parent (except for the first vertex which has no parent)
			return true
		}
	}

	if matrix.directed {
		delete(records, vertex)
	}

	return false
}
