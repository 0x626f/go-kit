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
