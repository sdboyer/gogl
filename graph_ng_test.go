package gogl

import (
	"gopkg.in/fatih/set.v0"
	. "launchpad.net/gocheck"
	"math"
)

/////////////////////////////////////////////////////////////////////
//
// GRAPH FIXTURES
//
/////////////////////////////////////////////////////////////////////

var graphFixtures = make(map[string]Graph)

var baseWeightedEdgeSet = []WeightedEdge{
	BaseWeightedEdge{BaseEdge{1, 2}, 5},
	BaseWeightedEdge{BaseEdge{2, 3}, -5},
}

func init() {
	// TODO use hardcoded fixtures, like the NullGraph (...?)
	// TODO improve naming basis/patterns for these
	base := BMBD.Create()
	base.AddEdges(edgeSet...)
	graphFixtures["2e3v"] = BIBD.From(base).Create()

	base.AddEdges(BaseEdge{"foo", "qux"})
	base2 := BMBD.From(base).Create()
	graphFixtures["3e4v"] = BIBD.From(base).Create()

	base.EnsureVertex("isolate")
	graphFixtures["3e5v1i"] = BIBD.From(base).Create()

	base2.AddEdges(BaseEdge{"foo", "qux"}, BaseEdge{"qux", "bar"})
	graphFixtures["arctest"] = BIBD.From(base2).Create()

	pair := BMBD.Create()
	pair.AddEdges(BaseEdge{1, 2})
	graphFixtures["pair"] = BIBD.From(pair).Create()

	weightedbase := BMWD.Create()
	weightedbase.AddEdges(baseWeightedEdgeSet...)
	graphFixtures["w-2e3v"] = BMWD.From(weightedbase).Create()
}

var _ = SetUpTestsFromBuilder(BMBD)
var _ = SetUpTestsFromBuilder(BMBU)
var _ = SetUpTestsFromBuilder(BIBD)
var _ = SetUpTestsFromBuilder(BMWD)
var _ = SetUpTestsFromBuilder(BMWU)
var _ = SetUpTestsFromBuilder(BMLD)
var _ = SetUpTestsFromBuilder(BMLU)
var _ = SetUpTestsFromBuilder(BMPD)
var _ = SetUpTestsFromBuilder(BMPU)

func SetUpTestsFromBuilder(b GraphBuilder) bool {
	var directed bool

	g := b.Graph()

	if _, ok := g.(DirectedGraph); ok {
		directed = true
		Suite(&DirectedGraphSuite{Builder: b})
	}

	// Set up the basic Graph suite unconditionally
	Suite(&GraphSuiteNG{b, directed})

	if _, ok := g.(SimpleGraph); ok {
		Suite(&SimpleGraphSuite{b, directed})
	}

	if _, ok := g.(MutableGraph); ok {
		Suite(&MutableGraphSuite{b, directed})
	}

	if _, ok := g.(WeightedGraph); ok {
		Suite(&WeightedGraphSuite{b, directed})
	}

	if _, ok := g.(MutableWeightedGraph); ok {
		Suite(&MutableWeightedGraphSuite{b, directed})
	}

	return true
}

/* GraphSuiteNG - tests for non-mutable graph methods */

type GraphSuiteNG struct {
	Builder  GraphBuilder
	Directed bool
}

func (s *GraphSuiteNG) TestHasVertex(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()
	c.Assert(g.HasVertex("qux"), Equals, false)
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *GraphSuiteNG) TestHasEdge(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()
	c.Assert(g.HasEdge(edgeSet[0]), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{"qux", "quark"}), Equals, false)
}

func (s *GraphSuiteNG) TestEachVertex(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()

	vset := set.NewNonTS()
	var hit int
	g.EachVertex(func(v Vertex) (terminate bool) {
		hit++
		vset.Add(v)
		return
	})

	c.Assert(vset.Has("foo"), Equals, true)
	c.Assert(vset.Has("bar"), Equals, true)
	c.Assert(vset.Has("baz"), Equals, true)
	c.Assert(hit, Equals, 3)
}

func (s *GraphSuiteNG) TestEachVertexTermination(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()

	var hit int
	g.EachVertex(func(v Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuiteNG) TestEachEdge(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()

	var hit int
	g.EachEdge(func(e Edge) (terminate bool) {
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
}

func (s *GraphSuiteNG) TestEachEdgeTermination(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()

	var hit int
	g.EachEdge(func(e Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuiteNG) TestEachAdjacent(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()

	vset := set.NewNonTS()
	var hit int
	g.EachAdjacent("bar", func(adj Vertex) (terminate bool) {
		hit++
		vset.Add(adj)
		return
	})

	c.Assert(vset.Has("foo"), Equals, true)
	c.Assert(vset.Has("bar"), Equals, false)
	c.Assert(vset.Has("baz"), Equals, true)
	c.Assert(hit, Equals, 2)
}

func (s *GraphSuiteNG) TestEachAdjacentTermination(c *C) {
	g := s.Builder.Using(graphFixtures["3e4v"]).Graph()

	var hit int
	g.EachAdjacent("foo", func(adjacent Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuiteNG) TestEachEdgeIncidentTo(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()

	flipset := []Edge{
		edgeSet[0].(BaseEdge).swap(),
		edgeSet[1].(BaseEdge).swap(),
	}

	eset := set.NewNonTS()
	var hit int
	g.EachEdgeIncidentTo("foo", func(e Edge) (terminate bool) {
		hit++
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(BaseEdge{U: e.Source(), V: e.Target()})
		return
	})

	c.Assert(hit, Equals, 1)
	if s.Directed {
		c.Assert(eset.Has(edgeSet[0]), Equals, true)
		c.Assert(eset.Has(edgeSet[1]), Equals, false)
	} else {
		c.Assert(eset.Has(edgeSet[0]) != eset.Has(flipset[0]), Equals, true)
		c.Assert(eset.Has(edgeSet[1]) != eset.Has(flipset[1]), Equals, false)
		c.Assert(eset.Has(edgeSet[1]), Equals, false)
	}

	eset = set.NewNonTS()
	g.EachEdgeIncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(BaseEdge{U: e.Source(), V: e.Target()})
		return
	})

	c.Assert(hit, Equals, 3)
	if s.Directed {
		c.Assert(eset.Has(edgeSet[0]), Equals, true)
		c.Assert(eset.Has(edgeSet[1]), Equals, true)
	} else {
		c.Assert(eset.Has(edgeSet[0]) != eset.Has(flipset[0]), Equals, true)
		c.Assert(eset.Has(edgeSet[1]) != eset.Has(flipset[1]), Equals, true)
	}
}

func (s *GraphSuiteNG) TestEachEdgeIncidentToTermination(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph()

	var hit int
	g.EachEdgeIncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuiteNG) TestDegreeOf(c *C) {
	g := s.Builder.Using(graphFixtures["3e5v1i"]).Graph()

	count, exists := g.DegreeOf("foo")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 2)

	count, exists = g.DegreeOf("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 2)

	count, exists = g.DegreeOf("baz")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.DegreeOf("qux")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.DegreeOf("isolate")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.DegreeOf("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *GraphSuiteNG) TestSize(c *C) {
	c.Assert(s.Builder.Graph().Size(), Equals, 0)
	c.Assert(s.Builder.Using(graphFixtures["2e3v"]).Graph().Size(), Equals, 2)
}

func (s *GraphSuiteNG) TestOrder(c *C) {
	c.Assert(s.Builder.Graph().Order(), Equals, 0)
	c.Assert(s.Builder.Using(graphFixtures["2e3v"]).Graph().Order(), Equals, 3)
}

/* DirectedGraphSuite - tests for directed graph methods */

type DirectedGraphSuite struct {
	Builder GraphBuilder
}

func (s *DirectedGraphSuite) TestTranspose(c *C) {
	g := s.Builder.Using(graphFixtures["2e3v"]).Graph().(DirectedGraph)

	g2 := g.Transpose()

	c.Assert(g2.HasEdge(edgeSet[0].(BaseEdge).swap()), Equals, true)
	c.Assert(g2.HasEdge(edgeSet[1].(BaseEdge).swap()), Equals, true)

	c.Assert(g2.HasEdge(edgeSet[0].(BaseEdge)), Equals, false)
	c.Assert(g2.HasEdge(edgeSet[1].(BaseEdge)), Equals, false)
}

func (s *DirectedGraphSuite) TestOutDegreeOf(c *C) {
	g := s.Builder.Using(graphFixtures["3e5v1i"]).Graph().(DirectedGraph)

	count, exists := g.OutDegreeOf("foo")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 2)

	count, exists = g.OutDegreeOf("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.OutDegreeOf("baz")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.OutDegreeOf("qux")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.DegreeOf("isolate")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.OutDegreeOf("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *DirectedGraphSuite) TestInDegreeOf(c *C) {
	g := s.Builder.Using(graphFixtures["3e5v1i"]).Graph().(DirectedGraph)

	count, exists := g.InDegreeOf("foo")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.InDegreeOf("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.InDegreeOf("baz")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.InDegreeOf("qux")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.DegreeOf("isolate")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.InDegreeOf("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *DirectedGraphSuite) TestEachArcTo(c *C) {
	g := s.Builder.Using(graphFixtures["arctest"]).Graph().(DirectedGraph)

	eset := set.NewNonTS()
	var hit int
	g.EachArcTo("foo", func(e Edge) (terminate bool) {
		c.Error("Vertex 'foo' should have no in-edges")
		c.FailNow()
		return
	})

	g.EachArcTo("bar", func(e Edge) (terminate bool) {
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(BaseEdge{U: e.Source(), V: e.Target()})
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
	c.Assert(eset.Has(edgeSet[0]), Equals, true)
	c.Assert(eset.Has(edgeSet[1]), Equals, false)
	c.Assert(eset.Has(BaseEdge{"qux", "bar"}), Equals, true)
}

func (s *DirectedGraphSuite) TestEachArcToTermination(c *C) {
	g := s.Builder.Using(graphFixtures["arctest"]).Graph().(DirectedGraph)

	var hit int
	g.EachArcTo("baz", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DirectedGraphSuite) TestEachArcFrom(c *C) {
	g := s.Builder.Using(graphFixtures["arctest"]).Graph().(DirectedGraph)

	eset := set.NewNonTS()
	var hit int
	g.EachArcFrom("baz", func(e Edge) (terminate bool) {
		c.Error("Vertex 'baz' should have no out-edges")
		c.FailNow()
		return
	})

	g.EachArcFrom("foo", func(e Edge) (terminate bool) {
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(BaseEdge{U: e.Source(), V: e.Target()})
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
	c.Assert(eset.Has(edgeSet[0]), Equals, true)
	c.Assert(eset.Has(edgeSet[1]), Equals, false)
	c.Assert(eset.Has(BaseEdge{"foo", "qux"}), Equals, true)
}

func (s *DirectedGraphSuite) TestEachArcFromTermination(c *C) {
	g := s.Builder.Using(graphFixtures["arctest"]).Graph().(DirectedGraph)

	var hit int
	g.EachArcFrom("foo", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

/* SimpleGraphSuite - tests for simple graph methods */

type SimpleGraphSuite struct {
	Builder  GraphBuilder
	Directed bool
}

func (s *SimpleGraphSuite) TestDensity(c *C) {
	c.Assert(math.IsNaN(s.Builder.Graph().(SimpleGraph).Density()), DeepEquals, true)

	g := s.Builder.Using(graphFixtures["pair"]).Graph().(SimpleGraph)
	if s.Directed {
		c.Assert(g.Density(), Equals, float64(0.5))
	} else {
		c.Assert(g.Density(), Equals, float64(1))
	}

	g = s.Builder.Using(graphFixtures["2e3v"]).Graph().(SimpleGraph)
	if s.Directed {
		c.Assert(g.Density(), Equals, float64(2)/float64(6))
	} else {
		c.Assert(g.Density(), Equals, float64(2)/float64(3))
	}
}

/* MutableGraphSuite - tests for mutable graph methods */

type MutableGraphSuite struct {
	Builder  GraphBuilder
	Directed bool
}

func (s *MutableGraphSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Builder.Graph().(MutableGraph)

	g.EnsureVertex()
	c.Assert(g.Order(), Equals, 0)

	g.RemoveVertex()
	c.Assert(g.Order(), Equals, 0)

	g.AddEdges()
	c.Assert(g.Order(), Equals, 0)
	c.Assert(g.Size(), Equals, 0)

	g.RemoveEdges()
	c.Assert(g.Order(), Equals, 0)
	c.Assert(g.Size(), Equals, 0)
}

func (s *MutableGraphSuite) TestEnsureVertex(c *C) {
	g := s.Builder.Graph().(MutableGraph)

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Builder.Graph().(MutableGraph)

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableGraphSuite) TestRemoveVertex(c *C) {
	g := s.Builder.Graph().(MutableGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Builder.Graph().(MutableGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Builder.Graph().(MutableGraph)
	g.AddEdges(&BaseEdge{1, 2})

	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed)

	// Now test removal
	g.RemoveEdges(&BaseEdge{1, 2})
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, false)
}

func (s *MutableGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Builder.Graph().(MutableGraph)

	g.AddEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3})

	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed)
	c.Assert(g.HasEdge(BaseEdge{3, 2}), Equals, !s.Directed)

	// Now test removal
	g.RemoveEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3})
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, false)
}

// Checks to ensure that removal works for both in-edges and out-edges.
func (s *MutableGraphSuite) TestVertexRemovalAlsoRemovesConnectedEdges(c *C) {
	g := s.Builder.Graph().(MutableGraph)

	g.AddEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3}, &BaseEdge{4, 1})
	g.RemoveVertex(1)

	c.Assert(g.Size(), Equals, 1)
}

/* WeightedGraphSuite - tests for weighted graphs */

type WeightedGraphSuite struct {
	Builder  GraphBuilder
	Directed bool
}

func (s *WeightedGraphSuite) TestEachEdge(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachEdge() iterator actually do implement WeightedEdge.
	g := s.Builder.Using(graphFixtures["w-2e3v"]).Graph().(WeightedGraph)

	var we WeightedEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *WeightedGraphSuite) TestEachWeightedEdge(c *C) {
	g := s.Builder.Using(graphFixtures["w-2e3v"]).Graph().(WeightedGraph)

	edgeset := set.NewNonTS()
	g.EachWeightedEdge(func(e WeightedEdge) {
		edgeset.Add(e)
	})

	if s.Directed {
		c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{1, 2}, 5}), Equals, true)
		c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{2, 3}, -5}), Equals, true)
		c.Assert(edgeset.Has(BaseEdge{1, 2}), Equals, false)
		c.Assert(edgeset.Has(BaseEdge{2, 3}), Equals, false)
	} else {
		c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{1, 2}, 5}) != edgeset.Has(BaseWeightedEdge{BaseEdge{2, 1}, 5}), Equals, true)
		c.Assert(edgeset.Has(BaseWeightedEdge{BaseEdge{2, 3}, -5}) != edgeset.Has(BaseWeightedEdge{BaseEdge{3, 2}, -5}), Equals, true)
		c.Assert(edgeset.Has(BaseEdge{1, 2}) || edgeset.Has(BaseEdge{2, 1}), Equals, false)
		c.Assert(edgeset.Has(BaseEdge{2, 3}) || edgeset.Has(BaseEdge{3, 2}), Equals, false)
	}
}

func (s *WeightedGraphSuite) TestHasWeightedEdge(c *C) {
	g := s.Builder.Using(graphFixtures["w-2e3v"]).Graph().(WeightedGraph)

	// TODO figure out how to meaningfully test undirected graphs' logic here
	c.Assert(g.HasWeightedEdge(baseWeightedEdgeSet[0]), Equals, true)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 1}), Equals, false) // wrong weight
}

/* MutableWeightedGraphSuite - tests for mutable weighted graphs */

type MutableWeightedGraphSuite struct {
	Builder  GraphBuilder
	Directed bool
}

func (s *MutableWeightedGraphSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Builder.Graph().(MutableWeightedGraph)

	g.EnsureVertex()
	c.Assert(g.Order(), Equals, 0)

	g.RemoveVertex()
	c.Assert(g.Order(), Equals, 0)

	g.AddEdges()
	c.Assert(g.Order(), Equals, 0)
	c.Assert(g.Size(), Equals, 0)

	g.RemoveEdges()
	c.Assert(g.Order(), Equals, 0)
	c.Assert(g.Size(), Equals, 0)
}

func (s *MutableWeightedGraphSuite) TestEnsureVertex(c *C) {
	g := s.Builder.Graph().(MutableWeightedGraph)

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableWeightedGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Builder.Graph().(MutableWeightedGraph)

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableWeightedGraphSuite) TestRemoveVertex(c *C) {
	g := s.Builder.Graph().(MutableWeightedGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Builder.Graph().(MutableWeightedGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Builder.Graph().(MutableWeightedGraph)
	g.AddEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5})

	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed)

	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 5}), Equals, true)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 3}), Equals, false)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 1}, 5}), Equals, !s.Directed)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 1}, -3}), Equals, false)

	// Now test removal
	g.RemoveEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5})
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 5}), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Builder.Graph().(MutableWeightedGraph)
	g.AddEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5}, BaseWeightedEdge{BaseEdge{2, 3}, -5})

	// Basic edge tests first
	// We test both Has*Edge() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Edge() method.
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed) // only if undirected
	c.Assert(g.HasEdge(BaseEdge{3, 2}), Equals, !s.Directed) // only if undirected

	// Now weighted edge tests
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 5}), Equals, true)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 3}), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 1}, 5}), Equals, !s.Directed)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 1}, 3}), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 3}, -5}), Equals, true)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 3}, 1}), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{3, 2}, -5}), Equals, !s.Directed)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{3, 2}, 1}), Equals, false) // wrong weight

	// Now test removal
	g.RemoveEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5}, BaseWeightedEdge{BaseEdge{2, 3}, -5})
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 5}), Equals, false)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 3}, -5}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, false)
}
