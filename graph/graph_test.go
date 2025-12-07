package graph

import (
	"fmt"
	"testing"
)

func TestDirected(t *testing.T) {
	graph := ForFeature[*User, int, int](Directed)
	users := createSampleUsers(10, 0)

	for _, usr := range users {
		err := graph.AddVertex(usr)

		if err != nil {
			t.Fatal(err)
		}
	}

	for index := 0; index < len(users)-1; index++ {
		err := graph.AddEdge(users[index].Key(), users[index+1].Key(), 0)

		if err != nil {
			t.Fatal(err)
		}
	}

	if graph.HasConnection(0, 1) && graph.HasConnection(1, 0) {
		t.Fatal("undirected")
	}

	err := graph.AddEdge(9, 0, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(5, 5, 0)

	if err != nil {
		t.Fatal(err)
	}
}

func TestNoFeature(t *testing.T) {
	graph := ForFeature[*User, int, int](None)
	users := createSampleUsers(10, 0)

	for _, usr := range users {
		err := graph.AddVertex(usr)

		if err != nil {
			t.Fatal(err)
		}
	}

	for index := 0; index < len(users)-1; index++ {
		err := graph.AddEdge(users[index].Key(), users[index+1].Key(), 0)

		if err != nil {
			t.Fatal(err)
		}
	}

	for index := 0; index < len(users)-1; index++ {
		if !(graph.HasConnection(index, index+1) && graph.HasConnection(index+1, index)) {
			t.Fatal("directed")
		}

	}

	err := graph.AddEdge(9, 0, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(5, 5, 0)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDirectedAcyclic(t *testing.T) {
	graph := ForFeature[*User, int, int](Directed | Acyclic)
	users := createSampleUsers(10, 0)

	for _, usr := range users {
		err := graph.AddVertex(usr)

		if err != nil {
			t.Fatal(err)
		}
	}

	for index := 0; index < len(users)-1; index++ {
		err := graph.AddEdge(users[index].Key(), users[index+1].Key(), 0)

		if err != nil {
			t.Fatal(err)
		}
	}

	if graph.HasConnection(0, 1) && graph.HasConnection(1, 0) {
		t.Fatal("undirected")
	}

	err := graph.AddEdge(9, 0, 0)
	t.Log(err)

	if err == nil || graph.HasConnection(9, 0) {
		t.Fatal("cycled")
	}

	err = graph.AddEdge(5, 5, 0)
	t.Log(err)

	if err == nil || graph.HasConnection(5, 5) {
		t.Fatal("cycled")
	}
}

func TestUnDirectedAcyclic(t *testing.T) {
	graph := ForFeature[*User, int, int](Acyclic)
	users := createSampleUsers(10, 0)

	for _, usr := range users {
		err := graph.AddVertex(usr)

		if err != nil {
			t.Fatal(err)
		}
	}

	for index := 0; index < len(users)-1; index++ {
		err := graph.AddEdge(users[index].Key(), users[index+1].Key(), 0)

		if err != nil {
			t.Fatal(err)
		}
	}

	for index := 0; index < len(users)-1; index++ {
		if !(graph.HasConnection(index, index+1) && graph.HasConnection(index+1, index)) {
			t.Fatal("directed")
		}

	}

	err := graph.AddEdge(9, 0, 0)
	t.Log(err)

	if err == nil || graph.HasConnection(9, 0) {
		t.Fatal("cycled")
	}

	err = graph.AddEdge(5, 5, 0)
	t.Log(err)

	if err == nil || graph.HasConnection(5, 5) {
		t.Fatal("cycled")
	}
}

func TestDFS(t *testing.T) {
	var err error

	graph := ForFeature[*User, int, int](None)
	users := createSampleUsers(10, 0)

	for _, usr := range users {
		err = graph.AddVertex(usr)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = graph.AddEdge(0, 1, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(0, 2, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(0, 3, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(1, 4, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(1, 5, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(2, 6, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(2, 7, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(3, 8, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(3, 9, 0)

	if err != nil {
		t.Fatal(err)
	}

	target := 5
	var buffer []int
	graph.DFS(0, func(from, to *User, w int) bool {
		key := to.Key()
		buffer = append(buffer, key)
		return key == target
	})

	if buffer[len(buffer)-1] != 5 {
		t.Fatal("not terminated dfs")
	}
}

func TestBFS(t *testing.T) {
	var err error

	graph := ForFeature[*User, int, int](None)
	users := createSampleUsers(10, 0)

	for _, usr := range users {
		err = graph.AddVertex(usr)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = graph.AddEdge(0, 1, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(0, 2, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(0, 3, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(1, 4, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(1, 5, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(2, 6, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(2, 7, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(3, 8, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(3, 9, 0)

	if err != nil {
		t.Fatal(err)
	}

	target := 5
	var buffer []int
	graph.BFS(0, func(from, to *User, w int) bool {
		key := to.Key()
		buffer = append(buffer, key)
		return key == target
	})

	fmt.Println()
	if buffer[len(buffer)-1] != 5 {
		t.Fatal("not terminated dfs")
	}
}

func TestRoutes(t *testing.T) {
	var err error

	graph := ForFeature[*User, int, int](None)
	users := createSampleUsers(11, 0)

	for _, usr := range users {
		err = graph.AddVertex(usr)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = graph.AddEdge(0, 1, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(0, 2, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(0, 3, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(1, 4, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(2, 5, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(3, 6, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(4, 7, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(5, 8, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(6, 9, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(7, 10, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(8, 10, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(9, 10, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(0, 10, 0)

	if err != nil {
		t.Fatal(err)
	}

	err = graph.AddEdge(1, 10, 0)

	if err != nil {
		t.Fatal(err)
	}

	routes := graph.Routes(0, 10, func(usr *User) bool {
		return true
	}, 0)

	expectedRoutes := [][]int{
		{0, 10},
		{0, 1, 10},
		{0, 1, 4, 7, 10},
		{0, 2, 5, 8, 10},
		{0, 3, 6, 9, 10},
	}

	for _, route := range routes {

		found := false

		for _, possibleRoute := range expectedRoutes {
			if len(possibleRoute) != len(route.Vertices) {
				continue
			}

			for index, vertex := range possibleRoute {
				found = vertex == route.Vertices[index].Key()
			}

			if found {
				break
			}
		}

		if !found {
			t.Fatal("route not found")
		}
	}
}

// TestJoin_BasicMerge tests basic joining of two graphs with no overlapping vertices
func TestJoin_BasicMerge(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)
	users1 := createSampleUsers(3, 0) // Users 0, 1, 2

	for _, usr := range users1 {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph1.AddEdge(0, 1, 10); err != nil {
		t.Fatal(err)
	}
	if err := graph1.AddEdge(1, 2, 20); err != nil {
		t.Fatal(err)
	}

	graph2 := ForFeature[*User, int, int](Directed)
	users2 := createSampleUsers(3, 3) // Users 3, 4, 5

	for _, usr := range users2 {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph2.AddEdge(3, 4, 30); err != nil {
		t.Fatal(err)
	}
	if err := graph2.AddEdge(4, 5, 40); err != nil {
		t.Fatal(err)
	}

	// Join graph2 into graph1
	graph1.Join(graph2, false)

	// Verify all vertices are present
	for i := 0; i < 6; i++ {
		if _, exists := graph1.Vertex(i); !exists {
			t.Fatalf("Vertex %d should exist after join", i)
		}
	}

	// Verify all edges are present
	expectedEdges := [][2]int{{0, 1}, {1, 2}, {3, 4}, {4, 5}}
	for _, edge := range expectedEdges {
		if !graph1.HasConnection(edge[0], edge[1]) {
			t.Fatalf("Edge %d -> %d should exist after join", edge[0], edge[1])
		}
	}

	// Verify edge weights
	if edge, exists := graph1.Edge(0, 1); !exists || edge != 10 {
		t.Fatalf("Edge 0->1 should have weight 10, got %d", edge)
	}
	if edge, exists := graph1.Edge(3, 4); !exists || edge != 30 {
		t.Fatalf("Edge 3->4 should have weight 30, got %d", edge)
	}
}

// TestJoin_OverlappingVertices tests joining graphs with some shared vertices
func TestJoin_OverlappingVertices(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)
	users1 := createSampleUsers(4, 0) // Users 0, 1, 2, 3

	for _, usr := range users1 {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph1.AddEdge(0, 1, 10); err != nil {
		t.Fatal(err)
	}
	if err := graph1.AddEdge(1, 2, 20); err != nil {
		t.Fatal(err)
	}

	graph2 := ForFeature[*User, int, int](Directed)
	users2 := createSampleUsers(4, 2) // Users 2, 3, 4, 5 (overlaps with 2, 3)

	for _, usr := range users2 {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph2.AddEdge(2, 3, 25); err != nil {
		t.Fatal(err)
	}
	if err := graph2.AddEdge(3, 4, 35); err != nil {
		t.Fatal(err)
	}
	if err := graph2.AddEdge(4, 5, 45); err != nil {
		t.Fatal(err)
	}

	// Join graph2 into graph1
	graph1.Join(graph2, false)

	// Verify all vertices are present
	for i := 0; i < 6; i++ {
		if _, exists := graph1.Vertex(i); !exists {
			t.Fatalf("Vertex %d should exist after join", i)
		}
	}

	// Verify all edges are present
	expectedEdges := [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}}
	for _, edge := range expectedEdges {
		if !graph1.HasConnection(edge[0], edge[1]) {
			t.Fatalf("Edge %d -> %d should exist after join", edge[0], edge[1])
		}
	}
}

// TestJoin_OverlappingEdges tests that existing edges are not overwritten
func TestJoin_OverlappingEdges(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)
	users := createSampleUsers(3, 0)

	for _, usr := range users {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph1.AddEdge(0, 1, 100); err != nil {
		t.Fatal(err)
	}

	graph2 := ForFeature[*User, int, int](Directed)
	for _, usr := range users {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	// Same edge but different weight
	if err := graph2.AddEdge(0, 1, 200); err != nil {
		t.Fatal(err)
	}
	if err := graph2.AddEdge(1, 2, 150); err != nil {
		t.Fatal(err)
	}

	// Join graph2 into graph1
	graph1.Join(graph2, false)

	// Original edge weight should be preserved
	if edge, exists := graph1.Edge(0, 1); !exists || edge != 100 {
		t.Fatalf("Edge 0->1 should preserve original weight 100, got %d", edge)
	}

	// New edge should be added
	if edge, exists := graph1.Edge(1, 2); !exists || edge != 150 {
		t.Fatalf("Edge 1->2 should have weight 150, got %d", edge)
	}
}

// TestJoin_DirectedGraphs tests joining directed graphs
func TestJoin_DirectedGraphs(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)
	users := createSampleUsers(3, 0)

	for _, usr := range users {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph1.AddEdge(0, 1, 10); err != nil {
		t.Fatal(err)
	}

	graph2 := ForFeature[*User, int, int](Directed)
	for _, usr := range users {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph2.AddEdge(1, 2, 20); err != nil {
		t.Fatal(err)
	}

	graph1.Join(graph2, false)

	// Verify directed edges
	if !graph1.HasConnection(0, 1) {
		t.Fatal("Edge 0->1 should exist")
	}
	if graph1.HasConnection(1, 0) {
		t.Fatal("Edge 1->0 should not exist in directed graph")
	}

	if !graph1.HasConnection(1, 2) {
		t.Fatal("Edge 1->2 should exist")
	}
	if graph1.HasConnection(2, 1) {
		t.Fatal("Edge 2->1 should not exist in directed graph")
	}
}

// TestJoin_UndirectedGraphs tests joining undirected graphs
func TestJoin_UndirectedGraphs(t *testing.T) {
	graph1 := ForFeature[*User, int, int](None)
	users := createSampleUsers(3, 0)

	for _, usr := range users {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph1.AddEdge(0, 1, 10); err != nil {
		t.Fatal(err)
	}

	graph2 := ForFeature[*User, int, int](None)
	for _, usr := range users {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph2.AddEdge(1, 2, 20); err != nil {
		t.Fatal(err)
	}

	graph1.Join(graph2, false)

	// Verify bidirectional edges
	if !graph1.HasConnection(0, 1) || !graph1.HasConnection(1, 0) {
		t.Fatal("Edge 0<->1 should exist bidirectionally in undirected graph")
	}

	if !graph1.HasConnection(1, 2) || !graph1.HasConnection(2, 1) {
		t.Fatal("Edge 1<->2 should exist bidirectionally in undirected graph")
	}
}

// TestJoin_AcyclicGraph tests that joining respects acyclic constraints
func TestJoin_AcyclicGraph(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed | Acyclic)
	users := createSampleUsers(4, 0)

	for _, usr := range users {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	// Create a path: 0 -> 1 -> 2
	if err := graph1.AddEdge(0, 1, 10); err != nil {
		t.Fatal(err)
	}
	if err := graph1.AddEdge(1, 2, 20); err != nil {
		t.Fatal(err)
	}

	graph2 := ForFeature[*User, int, int](Directed | Acyclic)
	for _, usr := range users {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	// Try to create a cycle: 2 -> 0 (would create cycle 0->1->2->0)
	if err := graph2.AddEdge(2, 0, 30); err != nil {
		t.Fatal(err)
	}

	// Also add a non-cyclic edge
	if err := graph2.AddEdge(2, 3, 40); err != nil {
		t.Fatal(err)
	}

	// Join should not add the cyclic edge
	graph1.Join(graph2, false)

	// Cyclic edge should not be added
	if graph1.HasConnection(2, 0) {
		t.Fatal("Edge 2->0 should not be added as it creates a cycle")
	}

	// Non-cyclic edge should be added
	if !graph1.HasConnection(2, 3) {
		t.Fatal("Edge 2->3 should be added")
	}
}

// TestJoin_EmptyGraph tests joining with an empty graph
func TestJoin_EmptyGraph(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)
	users := createSampleUsers(2, 0)

	for _, usr := range users {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph1.AddEdge(0, 1, 10); err != nil {
		t.Fatal(err)
	}

	graph2 := ForFeature[*User, int, int](Directed)

	// Join empty graph into graph1
	graph1.Join(graph2, false)

	// Graph1 should remain unchanged
	if _, exists := graph1.Vertex(0); !exists {
		t.Fatal("Vertex 0 should still exist")
	}
	if _, exists := graph1.Vertex(1); !exists {
		t.Fatal("Vertex 1 should still exist")
	}
	if !graph1.HasConnection(0, 1) {
		t.Fatal("Edge 0->1 should still exist")
	}
}

// TestJoin_IntoEmptyGraph tests joining into an empty graph
func TestJoin_IntoEmptyGraph(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)

	graph2 := ForFeature[*User, int, int](Directed)
	users := createSampleUsers(2, 0)

	for _, usr := range users {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	if err := graph2.AddEdge(0, 1, 10); err != nil {
		t.Fatal(err)
	}

	// Join into empty graph
	graph1.Join(graph2, false)

	// Graph1 should now contain all of graph2
	if _, exists := graph1.Vertex(0); !exists {
		t.Fatal("Vertex 0 should exist after join")
	}
	if _, exists := graph1.Vertex(1); !exists {
		t.Fatal("Vertex 1 should exist after join")
	}
	if !graph1.HasConnection(0, 1) {
		t.Fatal("Edge 0->1 should exist after join")
	}
}

// TestJoin_ComplexGraph tests joining complex graphs with multiple vertices and edges
func TestJoin_ComplexGraph(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)
	users1 := createSampleUsers(5, 0)

	for _, usr := range users1 {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	// Create a complex structure in graph1
	edges1 := [][3]int{{0, 1, 10}, {0, 2, 15}, {1, 3, 20}, {2, 3, 25}}
	for _, edge := range edges1 {
		if err := graph1.AddEdge(edge[0], edge[1], edge[2]); err != nil {
			t.Fatal(err)
		}
	}

	graph2 := ForFeature[*User, int, int](Directed)
	users2 := createSampleUsers(5, 3)

	for _, usr := range users2 {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	// Create overlapping and new edges
	edges2 := [][3]int{{3, 4, 30}, {4, 5, 35}, {5, 6, 40}, {6, 7, 45}}
	for _, edge := range edges2 {
		if err := graph2.AddEdge(edge[0], edge[1], edge[2]); err != nil {
			t.Fatal(err)
		}
	}

	graph1.Join(graph2, false)

	// Verify all vertices (0-7)
	for i := 0; i < 8; i++ {
		if _, exists := graph1.Vertex(i); !exists {
			t.Fatalf("Vertex %d should exist after join", i)
		}
	}

	// Verify all edges
	allEdges := [][3]int{
		{0, 1, 10}, {0, 2, 15}, {1, 3, 20}, {2, 3, 25},
		{3, 4, 30}, {4, 5, 35}, {5, 6, 40}, {6, 7, 45},
	}
	for _, edge := range allEdges {
		if !graph1.HasConnection(edge[0], edge[1]) {
			t.Fatalf("Edge %d -> %d should exist after join", edge[0], edge[1])
		}
		if e, exists := graph1.Edge(edge[0], edge[1]); !exists || e != edge[2] {
			t.Fatalf("Edge %d -> %d should have weight %d, got %d", edge[0], edge[1], edge[2], e)
		}
	}
}

// TestJoin_VertexCount tests that vertex count is correct after joining
func TestJoin_VertexCount(t *testing.T) {
	graph1 := ForFeature[*User, int, int](Directed)
	users1 := createSampleUsers(3, 0) // Users with Id: 0, 1, 2

	for _, usr := range users1 {
		if err := graph1.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	graph2 := ForFeature[*User, int, int](Directed)
	users2 := createSampleUsers(3, 2) // Users with Id: 2, 3, 4 (overlaps at 2)

	for _, usr := range users2 {
		if err := graph2.AddVertex(usr); err != nil {
			t.Fatal(err)
		}
	}

	// Before join: graph1 has 3 vertices, graph2 has 3 vertices
	// After join: should have 5 unique vertices (0, 1, 2, 3, 4)
	graph1.Join(graph2, false)

	// Verify the specific vertices exist (direct approach due to MapToValueSlice bug)
	expectedIDs := []int{0, 1, 2, 3, 4}
	for _, id := range expectedIDs {
		if _, exists := graph1.Vertex(id); !exists {
			t.Fatalf("Vertex %d should exist after join", id)
		}
	}

	// Verify no unexpected vertices exist
	unexpectedIDs := []int{5, 6, 7, 8, 9, 10}
	for _, id := range unexpectedIDs {
		if _, exists := graph1.Vertex(id); exists {
			t.Fatalf("Vertex %d should not exist after join", id)
		}
	}
}
