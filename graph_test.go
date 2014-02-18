package gogl

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	. "launchpad.net/gocheck"
	"testing"
)

var _ = fmt.Println

var edgeSet = []Edge{
	&BaseEdge{"foo", "bar"},
	&BaseEdge{"bar", "baz"},
}

type GraphFactory struct {
	CreateMutableGraph func() MutableGraph
	CreateGraph        func([]Edge) Graph
}

type GraphSuite struct {
	Graph   Graph
	Factory *GraphFactory
}

type MutableGraphSuite struct {
	Graph   MutableGraph
	Factory *GraphFactory
}

func (s *MutableGraphSuite) SetUpTest(c *C) {
	s.Graph = s.Factory.CreateMutableGraph()
}

func (s *MutableGraphSuite) TestEnsureVertex(c *C) {
	s.Graph.EnsureVertex("foo")
	c.Assert(s.Graph.HasVertex("foo"), Equals, true)
}

func (s *MutableGraphSuite) TestMultiEnsureVertex(c *C) {
	s.Graph.EnsureVertex("bar", "baz")
	c.Assert(s.Graph.HasVertex("bar"), Equals, true)
	c.Assert(s.Graph.HasVertex("baz"), Equals, true)
}

func (s *MutableGraphSuite) TestRemoveVertex(c *C) {
	s.Graph.EnsureVertex("bar", "baz")
	s.Graph.RemoveVertex("bar")
	c.Assert(s.Graph.HasVertex("bar"), Equals, false)
}

func (s *MutableGraphSuite) TestMultiRemoveVertex(c *C) {
	s.Graph.EnsureVertex("bar", "baz")
	s.Graph.RemoveVertex("bar", "baz")
	c.Assert(s.Graph.HasVertex("bar"), Equals, false)
	c.Assert(s.Graph.HasVertex("baz"), Equals, false)
}

func GraphTestVertexMembership(f *GraphFactory, t *testing.T) {
	g := f.CreateMutableGraph()

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

func GraphTestVertexMultiOps(f *GraphFactory, t *testing.T) {
	g := f.CreateMutableGraph()

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

func GraphTestRemoveVertexWithEdges(f *GraphFactory, t *testing.T) {
	g := f.CreateMutableGraph()

	g.AddEdge(edgeSet[0])
	g.AddEdge(edgeSet[1])

	Convey("Ensure outdegree is decremented when vertex is removed.", t, func() {
		g.RemoveVertex("bar")
		count, exists := g.OutDegree("foo")
		So(count, ShouldEqual, 0)
		So(exists, ShouldEqual, true)
	})
}

func GraphTestEachVertex(f *GraphFactory, t *testing.T) {
	g := f.CreateMutableGraph()

	var hit int
	it := func(v Vertex) {
		hit++
	}

	Convey("With no vertices, EachVertex does not call the injected closure at all.", t, func() {
		g.EachVertex(it)
		So(hit, ShouldEqual, 0)
	})

	// Ensure clean state, since goconvey failures do not stop the test
	hit = 0

	Convey("With two vertices, EachVertex calls the injected closure twice.", t, func() {
		g.EnsureVertex("foo")
		g.EnsureVertex("bar")
		g.EachVertex(it)
		So(hit, ShouldEqual, 2)
	})
}
