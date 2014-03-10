package gogl

import (
	"fmt"
	"reflect"
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

// This function automatically sets up suites of tests for graphs according to
// they implement.
func SetUpSimpleGraphTests(g Graph) bool {
	gf := &GraphFactory2{g}

	// Set up the basic Graph suite unconditionally
	_ = Suite(&GraphSuite{Graph: g, F2: gf})

	mg, ok := g.(MutableGraph)
	if ok {
		_ = Suite(&MutableGraphSuite{Graph: mg, F2: gf})
	}

	return true
}

type GraphFactory struct {
	CreateMutableGraph func() MutableGraph
	CreateGraph        func([]Edge) Graph
}

type GraphFactory2 struct {
	sourceGraph Graph
}

func (gf *GraphFactory2) create() interface{} {
	return reflect.New(reflect.Indirect(reflect.ValueOf(gf.sourceGraph)).Type()).Interface()
}

func (gf *GraphFactory2) CreateEmptyGraph() Graph {
	return gf.create().(Graph)
}

func (gf *GraphFactory2) CreateGraphFromEdges(edges ...Edge) Graph {
	// For now just cheat and work through a Mutable interface
	base := gf.create()

	if mg, ok := base.(MutableGraph); ok {
		mg.AddEdges(edges...)
	} else if mwg, ok := base.(MutableWeightedGraph); ok {
		weighted_edges := make([]WeightedEdge, 0, len(edges))
		for _, edge := range edges {
			weighted_edges = append(weighted_edges, BaseWeightedEdge{BaseEdge{edge.Source(), edge.Target()}, 0})
		}
		mwg.AddEdges(weighted_edges...)
	} else {
		panic("Until GraphInitializers are made to work properly, all graphs have to be mutable for this testing harness to work.")
	}

	return base.(Graph)
}

func (gf *GraphFactory2) CreateMutableGraph() MutableGraph {
	return gf.create().(MutableGraph)
}

/* Factory interfaces for tests */

type GraphCreator interface {
	CreateEmptyGraph() Graph
	CreateGraphFromEdges(edges ...Edge) Graph
}

type MutableGraphCreator interface {
	CreateMutableGraph() MutableGraph
}

/* GraphSuite - tests for non-mutable graph methods */

type GraphSuite struct {
	Graph   Graph
	Factory *GraphFactory
	F2      GraphCreator
}

func (s *GraphSuite) SetUpTest(c *C) {
	s.Graph = s.F2.CreateGraphFromEdges(edgeSet...)
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
	g := s.F2.CreateGraphFromEdges(&BaseEdge{"foo", "bar"})

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
	g := s.F2.CreateGraphFromEdges(&BaseEdge{"foo", "bar"})

	count, exists := g.InDegree("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.InDegree("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *GraphSuite) TestSize(c *C) {
	c.Assert(s.Graph.Size(), Equals, 2)

	g := s.F2.CreateEmptyGraph()
	c.Assert(g.Size(), Equals, 0)
}

func (s *GraphSuite) TestOrder(c *C) {
	c.Assert(s.Graph.Size(), Equals, 2)

	g := s.F2.CreateEmptyGraph()
	c.Assert(g.Size(), Equals, 0)
}

/* MutableGraphSuite - tests for mutable graph methods */
type MutableGraphSuite struct {
	Graph   MutableGraph
	Factory *GraphFactory
	F2      MutableGraphCreator
}

func (s *MutableGraphSuite) SetUpTest(c *C) {
	s.Graph = s.F2.CreateMutableGraph()
	//s.Graph = s.Factory.CreateMutableGraph()
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
