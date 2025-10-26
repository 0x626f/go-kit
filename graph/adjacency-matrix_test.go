package graph

import (
	"testing"
)

type User struct {
	Id int `json:"id"`
}

func (user *User) Key() int {
	return user.Id
}

func createSampleUsers(count, startFrom int) []*User {
	var users []*User
	for index := range count {
		users = append(users, &User{Id: startFrom + index})
	}
	return users
}

func TestCyclesForDirected(t *testing.T) {

	var err error

	matrix := NewAdjacencyMatrix[*User, int, int](true)
	users := createSampleUsers(10, 0)

	for index := 0; index < len(users)-1; index++ {
		err = matrix.Add(users[index].Key(), users[index+1].Key(), 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	if matrix.Has(0, 1) && matrix.Has(1, 0) {
		t.Fatal("undirected graph")
	}

	cycled := matrix.HasCycles()

	if cycled {
		t.Fatal("cycled")
	}

	err = matrix.Add(users[9].Key(), users[0].Key(), 0)
	err = matrix.Add(users[3].Key(), users[0].Key(), 0)
	err = matrix.Add(users[3].Key(), users[0].Key(), 0)
	err = matrix.Add(users[8].Key(), users[1].Key(), 0)

	cycled = matrix.HasCycles()

	if !cycled {
		t.Fatal("not cycled")
	}
}

func TestCyclesForUnDirected(t *testing.T) {

	var err error

	matrix := NewAdjacencyMatrix[*User, int, int](false)
	users := createSampleUsers(10, 0)

	for index := 0; index < len(users)-1; index++ {
		err = matrix.Add(users[index].Key(), users[index+1].Key(), 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	if matrix.Has(0, 1) && !matrix.Has(1, 0) {
		t.Fatal("directed graph")
	}

	cycled := matrix.HasCycles()

	if cycled {
		t.Fatal("cycled")
	}

	err = matrix.Add(users[9].Key(), users[0].Key(), 0)
	err = matrix.Add(users[3].Key(), users[0].Key(), 0)
	err = matrix.Add(users[3].Key(), users[0].Key(), 0)
	err = matrix.Add(users[3].Key(), users[2].Key(), 0)
	err = matrix.Add(users[8].Key(), users[1].Key(), 0)

	cycled = matrix.HasCycles()

	if !cycled {
		t.Fatal("not cycled")
	}
}
