package graph

import "encoding/json"

// matrixProxy is a helper structure used for JSON serialization and deserialization
// of AdjacencyMatrix instances. It mirrors the internal structure of AdjacencyMatrix
// but makes the fields accessible for JSON marshaling.
//
// Type parameters:
//   - V: The vertex type, must implement Vertex[K]
//   - E: The edge type, must implement Edge
//   - K: The key type for vertices, must be comparable
type matrixProxy[V Vertex[K], E Edge, K comparable] struct {
	// Directed indicates whether the graph is directed (true) or undirected (false)
	Directed bool `json:"Directed"`
	// Data stores the adjacency matrix representation
	Data map[K]map[K]E `json:"data"`
}

// matrixProxyFrom creates a proxy object from an AdjacencyMatrix for JSON serialization.
//
// Parameters:
//   - matrix: The AdjacencyMatrix to convert to a proxy
//
// Returns a new matrixProxy instance containing the data from the matrix.
func matrixProxyFrom[V Vertex[K], E Edge, K comparable](matrix *AdjacencyMatrix[V, E, K]) *matrixProxy[V, E, K] {
	return &matrixProxy[V, E, K]{
		Directed: matrix.directed,
		Data:     matrix.data,
	}
}

// MarshalJSON implements the json.Marshaler interface for AdjacencyMatrix.
// This allows adjacency matrices to be serialized to JSON format.
//
// Returns:
//   - A JSON byte array representing the matrix
//   - An error if marshaling fails
func (matrix *AdjacencyMatrix[V, E, K]) MarshalJSON() ([]byte, error) {
	return json.Marshal(matrixProxyFrom(matrix))
}

// UnmarshalJSON implements the json.Unmarshaler interface for AdjacencyMatrix.
// This allows adjacency matrices to be deserialized from JSON format.
//
// Parameters:
//   - data: The JSON byte array to deserialize
//
// Returns an error if unmarshaling fails.
func (matrix *AdjacencyMatrix[V, E, K]) UnmarshalJSON(data []byte) error {
	proxy := &matrixProxy[V, E, K]{}
	err := json.Unmarshal(data, proxy)

	if err != nil {
		return err
	}

	matrix.directed = proxy.Directed
	matrix.data = proxy.Data

	return nil
}

// graphProxy is a helper structure used for JSON serialization and deserialization
// of Graph instances. It mirrors the internal structure of Graph but makes the
// fields accessible for JSON marshaling.
//
// Type parameters:
//   - V: The vertex type, must implement Vertex[K]
//   - E: The edge type, must implement Edge
//   - K: The key type for vertices, must be comparable
type graphProxy[V Vertex[K], E Edge, K comparable] struct {
	// Feature stores the graph's feature flags
	Feature Feature `json:"feature"`
	// Vertices maps vertex keys to vertex objects
	Vertices map[K]V `json:"vertices"`
	// Matrix is the adjacency matrix representation of the graph
	Matrix *AdjacencyMatrix[V, E, K] `json:"matrix"`
}

// graphProxyFrom creates a proxy object from a Graph for JSON serialization.
//
// Parameters:
//   - graph: The Graph to convert to a proxy
//
// Returns a new graphProxy instance containing the data from the graph.
func graphProxyFrom[V Vertex[K], E Edge, K comparable](graph *Graph[V, E, K]) *graphProxy[V, E, K] {
	return &graphProxy[V, E, K]{
		Feature:  graph.feature,
		Vertices: graph.vertices,
		Matrix:   graph.matrix,
	}
}

// MarshalJSON implements the json.Marshaler interface for Graph.
// This allows graphs to be serialized to JSON format.
//
// Returns:
//   - A JSON byte array representing the graph
//   - An error if marshaling fails
func (graph *Graph[V, E, K]) MarshalJSON() ([]byte, error) {
	return json.Marshal(graphProxyFrom(graph))
}

// UnmarshalJSON implements the json.Unmarshaler interface for Graph.
// This allows graphs to be deserialized from JSON format.
//
// Parameters:
//   - data: The JSON byte array to deserialize
//
// Returns an error if unmarshaling fails.
func (graph *Graph[V, E, K]) UnmarshalJSON(data []byte) error {
	proxy := &graphProxy[V, E, K]{}
	err := json.Unmarshal(data, proxy)

	if err != nil {
		return err
	}

	graph.feature = proxy.Feature
	graph.vertices = proxy.Vertices
	graph.matrix = proxy.Matrix

	return nil
}
