package gogl

import (
	"gopkg.in/fatih/set.v0"
	. "launchpad.net/gocheck"
)

/////////////////////////////////////////////////////////////////////
//
// GRAPH FIXTURES
//
/////////////////////////////////////////////////////////////////////

var graphFixtures = make(map[string]Graph)

func init() {
	base := BMBD.Create()
	base.AddEdges(edgeSet...)
	graphFixtures["2e3v"] = BIBD.From(base).Create()

	base.AddEdges(BaseEdge{"foo", "qux"})
	graphFixtures["3e4v"] = BIBD.From(base).Create()
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
	_, directed := b.Graph().(DirectedGraph)

	// Set up the basic Graph suite unconditionally
	Suite(&GraphSuiteNG{Builder: b, Directed: directed})

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
	g := s.Builder.Using(graphFixtures["3e4v"]).Graph()

	// TODO test vertex isolates...can't make them in current testing harness
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
