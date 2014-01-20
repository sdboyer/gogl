package gogl

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
)

var fml = fmt.Println

var edgeSet = []Edge{
	&BaseEdge{"foo", "bar"},
	&BaseEdge{"bar", "baz"},
}

func TestEnsureGraphInterfaces(t *testing.T) {
	// What is Go's best practice for ensuring the implementation of an interface?
	_ = Graph(NewDirectedAdjacencyList())
	_ = SimpleGraph(NewDirectedAdjacencyList())
	t.Log("Implements Graph interface as expected.")
}

func TestVertexMembership(t *testing.T) {
	g := NewDirectedAdjacencyList()

	Convey("Test adding, removal, and membership of string literal vertex.", t, func() {
		So(g.HasVertex("foo"), ShouldEqual, false)
		g.EnsureVertex("foo")
		So(g.HasVertex("foo"), ShouldEqual, true)
		g.RemoveVertex("foo")
		So(g.HasVertex("foo"), ShouldEqual, false)
	})

	Convey("Test adding, removal, and membership of int literal vertex.", t, func() {
		So(g.HasVertex(1), ShouldEqual, false)
		g.EnsureVertex(1)
		So(g.HasVertex(1), ShouldEqual, true)
		g.RemoveVertex(1)
		So(g.HasVertex(1), ShouldEqual, false)
	})

	Convey("Test adding, removal, and membership of composite literal vertex.", t, func() {
		So(g.HasVertex(edgeSet[0]), ShouldEqual, false)
		g.EnsureVertex(edgeSet[0])
		So(g.HasVertex(edgeSet[0]), ShouldEqual, true)

		Convey("No membership match on new struct with same values or new pointer", func() {
			So(g.HasVertex(BaseEdge{"foo", "bar"}), ShouldEqual, false)
			So(g.HasVertex(&BaseEdge{"foo", "bar"}), ShouldEqual, false)
		})

		g.RemoveVertex(edgeSet[0])
		So(g.HasVertex(edgeSet[0]), ShouldEqual, false)
	})
}

func TestNonSingleAddRemoveVertex(t *testing.T) {
	g := NewDirectedAdjacencyList()

	Convey("Add and remove multiple vertices at once.", t, func() {
		g.EnsureVertex("foo", 1, edgeSet[0])
		So(g.HasVertex("foo"), ShouldEqual, true)
		So(g.HasVertex(1), ShouldEqual, true)
		So(g.HasVertex(edgeSet[0]), ShouldEqual, true)

		g.RemoveVertex("foo", 1, edgeSet[0])
		So(g.HasVertex("foo"), ShouldEqual, false)
		So(g.HasVertex(1), ShouldEqual, false)
		So(g.HasVertex(edgeSet[0]), ShouldEqual, false)
	})

	Convey("Ensure zero-length param to add/remove work correctly as no-ops.", t, func() {
		g.EnsureVertex()
		So(g.Order(), ShouldEqual, 0)
		g.RemoveVertex()
		So(g.Order(), ShouldEqual, 0)
	})
}

func TestRemoveVertexWithEdges(t *testing.T) {
	g := NewDirectedAdjacencyListFromEdgeSet(edgeSet)

	Convey("Ensure outdegree is decremented when vertex is removed.", t, func() {
		g.RemoveVertex("bar")
		count, exists := g.OutDegree("foo")
		So(count, ShouldEqual, 0)
		So(exists, ShouldEqual, true)
	})
}

func TestEachVertex(t *testing.T) {
	g := NewDirectedAdjacencyList()

	var hit int
	f := func(v Vertex) {
		hit++
	}

	g.EachVertex(f)

	if hit != 0 {
		t.Error("EachVertex did not call provided closure expected number of times.")
	}

	g.EnsureVertex("foo")
	g.EnsureVertex("bar")
	g.EachVertex(f)

	if hit != 2 {
		t.Error("EachVertex did not call provided closure expected number of times.")
	}
}

func TestAddEdge(t *testing.T) {
	g := NewDirectedAdjacencyList()

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
	g := NewDirectedAdjacencyList()

	g.AddEdge(&BaseEdge{"foo", "bar"})

	var count int
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
	g := NewDirectedAdjacencyList()

	g.AddEdge(&BaseEdge{"bar", "foo"})

	var count int
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

func TestEachEdge(t *testing.T) {
	g := NewDirectedAdjacencyListFromEdgeSet(edgeSet)

	var hit int
	f := func(e Edge) {
		hit++
	}

	g.EachEdge(f)

	if hit != 2 {
		t.Error("Edge iterator should have been called twice, was called", hit, "times.")
	}
}

func TestSize(t *testing.T) {
	g := NewDirectedAdjacencyList()

	if g.Size() != 0 {
		t.Error("Graph initializes with non-zero size.")
	}

	g.AddEdge(&BaseEdge{"foo", "bar"})

	if g.Size() != 1 {
		t.Error("Graph does not increment size properly on edge addition.")
	}

	g.RemoveEdge(&BaseEdge{"foo", "bar"})

	if g.Size() != 0 {
		t.Error("Graph does not decrement size properly on edge removal.")
	}
}

func TestOrder(t *testing.T) {
	g := NewDirectedAdjacencyList()

	if g.Order() != 0 {
		t.Error("Graph initializes with non-zero order.")
	}

	g.EnsureVertex("foo")

	if g.Order() != 1 {
		t.Error("Adding a vertex does not increment order properly.")
	}

	g.RemoveVertex("foo")

	if g.Order() != 0 {
		t.Error("Removing a vertex does not decrement order properly.")
	}
}

func TestDensity(t *testing.T) {
	g := NewDirectedAdjacencyList()
	var density float64

	if !math.IsNaN(g.Density()) {
		t.Error("On graph initialize, Density should be NaN (divides zero by zero)).")
	}

	g.AddEdge(&BaseEdge{"foo", "bar"})

	density = g.Density()
	if density != 1 {
		t.Error("In undirected graph of V = 2 and E = 1, density should be 1; was", density)
	}

	g.AddEdge(&BaseEdge{"baz", "qux"})

	density = g.Density()
	if density != float64(1)/float64(3) {
		t.Error("In undirected graph of V = 4 and E = 2, density should be 0.3333; was", density)
	}
}
