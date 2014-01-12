package gogl

import (
	"testing"
)

func TestEnsureIsGraph(t *testing.T) {
	_ = Graph(NewAdjacencyList())
	t.Log("Implements Graph interface as expected.")
}

func TestGraphZeroValues(t *testing.T) {
	g := NewAdjacencyList()

	if g.Size() != 0 {
		t.Error("Initializes with non-zero size.")
	}

	if g.Order() != 0 {
		t.Error("Initializes with non-zero order.")
	}
}

func TestVertexMembership(t *testing.T) {
	g := NewAdjacencyList()

	if g.HasVertex("foo") != false {
		t.Error("Incorrectly reports nonexistent vertex as present.")
	}

	if g.AddVertex("foo") != true {
		t.Error("Fails to add string vertex correctly.")
	}

	if g.HasVertex("foo") != true {
		t.Error("Fails to locate existing string vertex.")
	}

	if g.AddVertex(1) != true {
		t.Error("Fails to add int vertex correctly.")
	}

	if g.RemoveVertex("foo") != true {
		t.Error("Reports incorrect failure on removing existing vertex.")
	}

	if g.HasVertex("foo") != false {
		t.Error("Reports vertex still present after removal.")
	}
}

func TestEachVertex(t *testing.T) {
	g := NewAdjacencyList()

	var hit uint
	f := func(v Vertex) {
		hit++
	}

	g.EachVertex(f)

	if hit != 0 {
		t.Error("EachVertex did not call provided closure expected number of times.")
	}

	g.AddVertex("foo")
	g.AddVertex("bar")
	g.EachVertex(f)

	if hit != 2 {
		t.Error("EachVertex did not call provided closure expected number of times.")
	}
}

func TestAddEdge(t *testing.T) {
	g := NewAdjacencyList()

	g.AddEdge(&BaseEdge{"foo", "bar"})

	if g.HasVertex("foo") != true {
		t.Error("AddEdge did not create vertices as expected.")
	}

	f := func(adj Vertex) {
		if adj != "bar" {
			t.Error("Adjacency relationship from foo to bar not reported correctly, ", adj, " was passed.")
		}
	}

	g.EachAdjacent("foo", f)
}

func TestOutDegree(t *testing.T) {
	g := NewAdjacencyList()

	g.AddEdge(&BaseEdge{"foo", "bar"})

	var count uint
	var exists bool
	count, exists = g.OutDegree("foo")

	if count != 1 {
		t.Error("Vertex should have outdegree of one, but", count, "was reported.")
	}

	if exists != true {
		t.Error("Vertex should exist.")
	}

	count, exists = g.OutDegree("bar")

	if count != 0 {
		t.Error("Vertex should have outdegree of zero, but", count, "was reported.")
	}

	if exists != true {
		t.Error("Vertex should exist.")
	}

	count, exists = g.OutDegree("baz")

	if count != 0 {
		t.Error("Zero outdegree count is reported when vertex does not exist.")
	}

	if exists != false {
		t.Error("Vertex should not exist.")
	}
}

func TestInDegree(t *testing.T) {
	g := NewAdjacencyList()

	g.AddEdge(&BaseEdge{"bar", "foo"})

	var count uint
	var exists bool
	count, exists = g.InDegree("foo")

	if count != 1 {
		t.Error("Vertex should have indegree of one, but", count, "was reported.")
	}

	if exists != true {
		t.Error("Vertex should exist.")
	}

	count, exists = g.InDegree("bar")

	if count != 0 {
		t.Error("Vertex should have indegree of zero, but", count, "was reported.")
	}

	if exists != true {
		t.Error("Vertex should exist.")
	}

	count, exists = g.InDegree("baz")

	if count != 0 {
		t.Error("Zero indegree count is reported when vertex does not exist.")
	}

	if exists != false {
		t.Error("Vertex should not exist.")
	}
}
