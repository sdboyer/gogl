package gogl

import (
	"fmt"
	"testing"

	"github.com/fatih/set"
	. "launchpad.net/gocheck"
)

var _ = fmt.Println

// Hook gocheck into the go test runner
func Test(t *testing.T) { TestingT(t) }

var edgeSet = []Edge{
	&BaseEdge{"foo", "bar"},
	&BaseEdge{"bar", "baz"},
}

type GraphFactory struct {
	CreateMutableGraph func() MutableGraph
	CreateGraph        func([]Edge) Graph
}

/* GraphSuite - tests for non-mutable graph methods */
type GraphSuite struct {
	Graph   Graph
	Factory *GraphFactory
}

func (s *GraphSuite) SetUpTest(c *C) {
	s.Graph = s.Factory.CreateGraph(edgeSet)
}

func (s *GraphSuite) TestVertexMembership(c *C) {
	c.Assert(s.Graph.HasVertex("foo"), Equals, true)
}

func (s *GraphSuite) TestEachVertex(c *C) {
	var hit int
	f := func(v Vertex) {
		hit++
		c.Log("EachVertex hit closure, hit count", hit)
	}

	s.Graph.EachVertex(f)
	if !c.Check(hit, Equals, 3) {
		c.Error("EachVertex should have called injected closure iterator 3 times, actual count was ", hit)
	}
}

func (s *GraphSuite) TestEachEdge(c *C) {
	var hit int
	f := func(e Edge) {
		hit++
		c.Log("EachAdjacent hit closure with edge pair ", e.Source(), " ", e.Target(), " at hit count ", hit)
	}

	s.Graph.EachEdge(f)
	if !c.Check(hit, Equals, 2) {
		c.Error("EachEdge should have called injected closure iterator 2 times, actual count was ", hit)
	}
}

func (s *GraphSuite) TestEachAdjacent(c *C) {
	var hit int
	f := func(adj Vertex) {
		hit++
		c.Log("EachAdjacent hit closure with vertex ", adj, " at hit count ", hit)
	}

	s.Graph.EachAdjacent("foo", f)
	if !c.Check(hit, Equals, 1) {
		c.Error("EachEdge should have called injected closure iterator 2 times, actual count was ", hit)
	}
}

// This test is carefully constructed to be fully correct for directed graphs,
// and incidentally correct for undirected graphs.
func (s *GraphSuite) TestOutDegree(c *C) {
	g := s.Factory.CreateGraph([]Edge{&BaseEdge{"foo", "bar"}})

	count, exists := g.OutDegree("foo")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.OutDegree("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

// This test is carefully constructed to be fully correct for directed graphs,
// and incidentally correct for undirected graphs.
func (s *GraphSuite) TestInDegree(c *C) {
	g := s.Factory.CreateGraph([]Edge{&BaseEdge{"foo", "bar"}})

	count, exists := g.InDegree("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.InDegree("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *GraphSuite) TestSize(c *C) {
	c.Assert(s.Graph.Size(), Equals, 2)

	g := s.Factory.CreateGraph([]Edge{})
	c.Assert(g.Size(), Equals, 0)
}

func (s *GraphSuite) TestOrder(c *C) {
	c.Assert(s.Graph.Size(), Equals, 2)

	g := s.Factory.CreateGraph([]Edge{})
	c.Assert(g.Size(), Equals, 0)
}

/* MutableGraphSuite - tests for mutable graph methods */
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

func (s *MutableGraphSuite) TestAddAndRemoveEdge(c *C) {
	s.Graph.AddEdges(&BaseEdge{1, 2})

	f := func(e Edge) {
		// Undirected graphs provide no guarantee of vertex output ordering,
		// and as such either {1,2} or {2,1} are valid outputs.
		// TODO this can be removed once a HasEdge() is implemented.
		if !c.Check(BaseEdge{e.Source(), e.Target()}, Equals, BaseEdge{1, 2}) {
			if c.Check(BaseEdge{e.Source(), e.Target()}, Equals, BaseEdge{2, 1}) {
				c.Succeed()
			} else {
				c.Log("Neither acceptable edge vertex ordering pair was provided.")
				c.FailNow()
			}
		}
	}

	s.Graph.EachEdge(f)

	// Now test removal
	f = func(e Edge) {
		c.Error("Graph should not contain any edges after removal.")
	}

	s.Graph.RemoveEdges(&BaseEdge{1, 2})
	s.Graph.EachEdge(f)
}

func (s *MutableGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	s.Graph.AddEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3})
	set := set.NewNonTS()

	f := func(e Edge) {
		set.Add(e)
	}

	s.Graph.EachEdge(f)
	c.Assert(set.Has(BaseEdge{1, 2}), Equals, true)
	c.Assert(set.Has(BaseEdge{2, 3}), Equals, true)

	// Now test removal
	f = func(e Edge) {
		c.Error("Graph should not contain any edges after removal.")
	}

	s.Graph.RemoveEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3})
	s.Graph.EachEdge(f)
}

// Checks to ensure that removal works for both in-edges and out-edges.
// TODO - make the edge membership check a little more robust.
func (s *MutableGraphSuite) TestVertexRemovalAlsoRemovesConnectedEdges(c *C) {
	s.Graph.AddEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3}, &BaseEdge{4, 1})
	s.Graph.RemoveVertex(1)

	c.Assert(s.Graph.Size(), Equals, 1)
}
