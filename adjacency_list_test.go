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

func TestAddVertex(t *testing.T) {
	g := NewAdjacencyList()

	if g.AddVertex("foo") != true {
		t.Error("Fails to add string vertex correctly.")
	}

	if g.AddVertex(1) != true {
		t.Error("Fails to add int vertex correctly.")
	}

}
