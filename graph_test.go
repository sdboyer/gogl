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
	gf := &GraphFactory{g}

	// Set up the basic Graph suite unconditionally
	_ = Suite(&GraphSuite{Graph: g, Factory: gf})

	mg, ok := g.(MutableGraph)
	if ok {
		_ = Suite(&MutableGraphSuite{Graph: mg, Factory: gf})
	}

	wg, ok := g.(WeightedGraph)
	if ok {
		_ = Suite(&WeightedGraphSuite{Graph: wg, Factory: gf})
	}

	mwg, ok := g.(MutableWeightedGraph)
	if ok {
		_ = Suite(&MutableWeightedGraphSuite{Graph: mwg, Factory: gf})
	}

	return true
}

type GraphFactory struct {
	sourceGraph Graph
}

func (gf *GraphFactory) create() interface{} {
	return reflect.New(reflect.Indirect(reflect.ValueOf(gf.sourceGraph)).Type()).Interface()
}

func (gf *GraphFactory) CreateEmptyGraph() Graph {
	return gf.create().(Graph)
}

func (gf *GraphFactory) CreateGraphFromEdges(edges ...Edge) Graph {
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

func (gf *GraphFactory) CreateMutableGraph() MutableGraph {
	return gf.create().(MutableGraph)
}

func (gf *GraphFactory) CreateWeightedGraphFromEdges(edges ...WeightedEdge) WeightedGraph {
	base := gf.create()
	if mwg, ok := base.(MutableWeightedGraph); ok {
		mwg.AddEdges(edges...)
		return mwg
	}

	panic("Until GraphInitializers are made to work properly, all graphs have to be mutable for this testing harness to work.")
}

func (gf *GraphFactory) CreateEmptyWeightedGraph() WeightedGraph {
	return gf.create().(WeightedGraph)
}

func (gf *GraphFactory) CreateMutableWeightedGraph() MutableWeightedGraph {
	return gf.create().(MutableWeightedGraph)
}

/* Factory interfaces for tests */

type GraphCreator interface {
	CreateEmptyGraph() Graph
	CreateGraphFromEdges(edges ...Edge) Graph
}

type MutableGraphCreator interface {
	CreateMutableGraph() MutableGraph
}

type WeightedGraphCreator interface {
	CreateEmptyWeightedGraph() WeightedGraph
	CreateWeightedGraphFromEdges(edges ...WeightedEdge) WeightedGraph
}

type MutableWeightedGraphCreator interface {
	CreateMutableWeightedGraph() MutableWeightedGraph
}

/* GraphSuite - tests for non-mutable graph methods */

type GraphSuite struct {
	Graph   Graph
	Factory GraphCreator
}

func (s *GraphSuite) SetUpTest(c *C) {
	s.Graph = s.Factory.CreateGraphFromEdges(edgeSet...)
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
	g := s.Factory.CreateGraphFromEdges(&BaseEdge{"foo", "bar"})

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
	g := s.Factory.CreateGraphFromEdges(&BaseEdge{"foo", "bar"})

	count, exists := g.InDegree("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.InDegree("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *GraphSuite) TestSize(c *C) {
	c.Assert(s.Graph.Size(), Equals, 2)

	g := s.Factory.CreateEmptyGraph()
	c.Assert(g.Size(), Equals, 0)
}

func (s *GraphSuite) TestOrder(c *C) {
	c.Assert(s.Graph.Size(), Equals, 2)

	g := s.Factory.CreateEmptyGraph()
	c.Assert(g.Size(), Equals, 0)
}

/* MutableGraphSuite - tests for mutable graph methods */

type MutableGraphSuite struct {
	Graph   MutableGraph
	Factory MutableGraphCreator
}

func (s *MutableGraphSuite) TestEnsureVertex(c *C) {
	g := s.Factory.CreateMutableGraph()

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Factory.CreateMutableGraph()

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableGraphSuite) TestRemoveVertex(c *C) {
	g := s.Factory.CreateMutableGraph()

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Factory.CreateMutableGraph()

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Factory.CreateMutableGraph()

	g.AddEdges(&BaseEdge{1, 2})

	edgeset := set.NewNonTS()

	g.EachEdge(func(e Edge) {
		edgeset.Add(e)
	})
	c.Assert(edgeset.Has(BaseEdge{1, 2}) != edgeset.Has(BaseEdge{2, 1}), Equals, true)

	// Now test removal
	g.RemoveEdges(&BaseEdge{1, 2})
	g.EachEdge(func(e Edge) {
		c.Error("Graph should not contain any edges after removal.")
	})
}

func (s *MutableGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory.CreateMutableGraph()

	g.AddEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3})

	set := set.NewNonTS()
	g.EachEdge(func(e Edge) {
		set.Add(e)
	})
	c.Assert(set.Has(BaseEdge{1, 2}), Equals, true)
	c.Assert(set.Has(BaseEdge{2, 3}), Equals, true)

	// Now test removal
	g.RemoveEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3})
	g.EachEdge(func(e Edge) {
		c.Error("Graph should not contain any edges after removal.")
	})
}

// Checks to ensure that removal works for both in-edges and out-edges.
// TODO - make the edge membership check a little more robust.
func (s *MutableGraphSuite) TestVertexRemovalAlsoRemovesConnectedEdges(c *C) {
	g := s.Factory.CreateMutableGraph()

	g.AddEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3}, &BaseEdge{4, 1})
	g.RemoveVertex(1)

	c.Assert(g.Size(), Equals, 1)
}

/* WeightedGraphSuite - tests for weighted graphs */

type WeightedGraphSuite struct {
	Graph   WeightedGraph
	Factory WeightedGraphCreator
}

func (s *WeightedGraphSuite) TestEachWeightedEdge(c *C) {
	g := s.Factory.CreateWeightedGraphFromEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5}, BaseWeightedEdge{BaseEdge{2, 3}, -5})

	edgeset := set.NewNonTS()
	g.EachWeightedEdge(func(e WeightedEdge) {
		edgeset.Add(e)
	})

	c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{1, 2}, 5}) != edgeset.Has(BaseWeightedEdge{BaseEdge{2, 1}, 5}), Equals, true)
	c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{2, 3}, -5}) != edgeset.Has(BaseWeightedEdge{BaseEdge{3, 2}, -5}), Equals, true)
	c.Assert(edgeset.Has(BaseEdge{1, 2}) || edgeset.Has(BaseEdge{2, 1}), Equals, false)
	c.Assert(edgeset.Has(BaseEdge{2, 3}) || edgeset.Has(BaseEdge{3, 2}), Equals, false)
}

/* MutableWeightedGraphSuite - tests for mutable weighted graphs */

type MutableWeightedGraphSuite struct {
	Graph   MutableWeightedGraph
	Factory MutableWeightedGraphCreator
}

func (s *MutableWeightedGraphSuite) TestEnsureVertex(c *C) {
	g := s.Factory.CreateMutableWeightedGraph()

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableWeightedGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Factory.CreateMutableWeightedGraph()

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableWeightedGraphSuite) TestRemoveVertex(c *C) {
	g := s.Factory.CreateMutableWeightedGraph()

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Factory.CreateMutableWeightedGraph()

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Factory.CreateMutableWeightedGraph()
	g.AddEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5})

	edgeset := set.NewNonTS()
	g.EachEdge(func(e Edge) {
		edgeset.Add(e)
	})
	c.Assert(edgeset.Has(BaseEdge{1, 2}) != edgeset.Has(BaseEdge{2, 1}), Equals, true)

	edgeset2 := set.NewNonTS()
	g.EachWeightedEdge(func(e WeightedEdge) {
		edgeset2.Add(e)
	})
	c.Assert(edgeset2.Has(BaseWeightedEdge{BaseEdge{1, 2}, 5}) || edgeset2.Has(BaseWeightedEdge{BaseEdge{2, 1}, 5}), Equals, true)

	// Now test removal
	g.RemoveEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5})
	g.EachEdge(func(e Edge) {
		c.Error("Graph should not contain any edges after removal.")
	})
}

func (s *MutableWeightedGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory.CreateMutableWeightedGraph()
	g.AddEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5}, BaseWeightedEdge{BaseEdge{2, 3}, -5})
	edgeset := set.NewNonTS()

	// First test with generic edge iteration method
	g.EachEdge(func(e Edge) {
		edgeset.Add(e)
	})

	// Check both directed and undirected edge permutations
	c.Assert(edgeset.Has(BaseEdge{1, 2}) != edgeset.Has(BaseEdge{2, 1}), Equals, true)
	c.Assert(edgeset.Has(BaseEdge{2, 3}) != edgeset.Has(BaseEdge{3, 2}), Equals, true)
	c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{1, 2}, 5}) || edgeset.Has(BaseWeightedEdge{BaseEdge{2, 1}, 5}), Equals, false)
	c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{2, 3}, -5}) || edgeset.Has(BaseWeightedEdge{BaseEdge{3, 2}, -5}), Equals, false)

	edgeset2 := set.NewNonTS()
	// Now with weighted edge iteration method
	g.EachWeightedEdge(func(e WeightedEdge) {
		edgeset2.Add(e)
	})
	c.Assert(edgeset2.Has(BaseWeightedEdge{BaseEdge{1, 2}, 5}) != edgeset2.Has(BaseWeightedEdge{BaseEdge{2, 1}, 5}), Equals, true)
	c.Assert(edgeset2.Has(BaseWeightedEdge{BaseEdge{2, 3}, -5}) != edgeset2.Has(BaseWeightedEdge{BaseEdge{3, 2}, -5}), Equals, true)
	c.Assert(edgeset2.Has(BaseEdge{1, 2}) || edgeset2.Has(BaseEdge{2, 1}), Equals, false)
	c.Assert(edgeset2.Has(BaseEdge{2, 3}) || edgeset2.Has(BaseEdge{3, 2}), Equals, false)

	// Now test removal
	g.RemoveEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5}, BaseWeightedEdge{BaseEdge{2, 3}, -5})
	g.EachEdge(func(e Edge) {
		c.Error("Graph should not contain any edges after removal.")
	})
	g.EachWeightedEdge(func(e WeightedEdge) {
		c.Error("Graph should not contain any edges after removal.")
	})
}
