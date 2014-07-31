package gogl

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	. "github.com/sdboyer/gocheck"
	"gopkg.in/fatih/set.v0"
)

/////////////////////////////////////////////////////////////////////
//
// GRAPH FIXTURES
//
/////////////////////////////////////////////////////////////////////

var graphFixtures = make(map[string]Graph)

var edgeSet = []Edge{
	BaseEdge{"foo", "bar"},
	BaseEdge{"bar", "baz"},
}

var baseWeightedEdgeSet = []WeightedEdge{
	BaseWeightedEdge{BaseEdge{1, 2}, 5.23},
	BaseWeightedEdge{BaseEdge{2, 3}, 5.821},
}

var baseLabeledEdgeSet = []LabeledEdge{
	BaseLabeledEdge{BaseEdge{1, 2}, "foo"},
	BaseLabeledEdge{BaseEdge{2, 3}, "bar"},
}

var baseDataEdgeSet = []DataEdge{
	BaseDataEdge{BaseEdge{1, 2}, "foo"},
	BaseDataEdge{BaseEdge{2, 3}, struct{ a int }{a: 2}},
}

func init() {
	// TODO use hardcoded fixtures, like the NullGraph (...?)
	// TODO improve naming basis/patterns for these
	spec := BuildGraph().Mutable().BasicEdges().Directed()
	base := spec.Create(AdjacencyList).(MutableGraph)
	base.AddEdges(edgeSet...)

	ispec := BuildGraph().Immutable().BasicEdges().Directed()
	graphFixtures["2e3v"] = ispec.Using(base).Create(AdjacencyList)

	base.AddEdges(BaseEdge{"foo", "qux"})
	base2 := spec.Using(base).Create(AdjacencyList).(MutableGraph)
	graphFixtures["3e4v"] = ispec.Using(base).Create(AdjacencyList)

	base.EnsureVertex("isolate")
	graphFixtures["3e5v1i"] = ispec.Using(base).Create(AdjacencyList)

	base2.AddEdges(BaseEdge{"foo", "qux"}, BaseEdge{"qux", "bar"})
	graphFixtures["arctest"] = ispec.Using(base2).Create(AdjacencyList)

	pair := spec.Using(nil).Create(AdjacencyList).(MutableGraph)
	pair.AddEdges(BaseEdge{1, 2})
	graphFixtures["pair"] = ispec.Using(pair).Create(AdjacencyList)

	wb := BuildGraph().Directed().WeightedEdges()
	weightedbase := wb.Create(AdjacencyList).(MutableWeightedGraph)
	weightedbase.AddEdges(baseWeightedEdgeSet...)
	graphFixtures["w-2e3v"] = wb.Using(weightedbase).Create(AdjacencyList)

	lb := BuildGraph().Directed().LabeledEdges()
	labeledbase := lb.Create(AdjacencyList).(MutableLabeledGraph)
	labeledbase.AddEdges(baseLabeledEdgeSet...)
	graphFixtures["l-2e3v"] = lb.Using(labeledbase).Create(AdjacencyList)

	db := BuildGraph().Directed().DataEdges()
	data_base := db.Create(AdjacencyList).(MutableDataGraph)
	data_base.AddEdges(baseDataEdgeSet...)
	graphFixtures["p-2e3v"] = db.Using(data_base).Create(AdjacencyList)

	for gp, _ := range alCreators {
		SetUpTestsFromSpec(gp, AdjacencyList)
	}
}

/////////////////////////////////////////////////////////////////////
//
// HELPERS
//
/////////////////////////////////////////////////////////////////////

// Hook gocheck into the go test runner
func TestHookup(t *testing.T) { TestingT(t) }

// swap method is useful for some testing shorthand
func (e BaseEdge) swap() Edge {
	return BaseEdge{e.V, e.U}
}

func gdebug(g Graph, args ...interface{}) {
	fmt.Println("DEBUG: graph type", reflect.New(reflect.Indirect(reflect.ValueOf(g)).Type()))
	pretty.Print(args...)
}

/////////////////////////////////////////////////////////////////////
//
// SUITE SETUP
//
/////////////////////////////////////////////////////////////////////

func SetUpTestsFromSpec(gp GraphProperties, fn func(GraphSpec) Graph) bool {
	var directed bool

	g := fn(GraphSpec{Props: gp})

	fact := func(gs GraphSource) Graph {
		return fn(GraphSpec{Props: gp, Source: gs})
	}

	if _, ok := g.(DirectedGraph); ok {
		directed = true
		Suite(&DirectedGraphSuite{Factory: fact})
	}

	// Set up the basic Graph suite unconditionally
	Suite(&GraphSuite{fact, directed})

	if _, ok := g.(SimpleGraph); ok {
		Suite(&SimpleGraphSuite{fact, directed})
	}

	if _, ok := g.(MutableGraph); ok {
		Suite(&MutableGraphSuite{fact, directed})
	}

	if _, ok := g.(WeightedGraph); ok {
		Suite(&WeightedGraphSuite{fact, directed})
	}

	if _, ok := g.(MutableWeightedGraph); ok {
		Suite(&MutableWeightedGraphSuite{fact, directed})
	}

	if _, ok := g.(LabeledGraph); ok {
		Suite(&LabeledGraphSuite{fact, directed})
	}

	if _, ok := g.(MutableLabeledGraph); ok {
		Suite(&MutableLabeledGraphSuite{fact, directed})
	}

	if _, ok := g.(DataGraph); ok {
		Suite(&DataGraphSuite{fact, directed})
	}

	if _, ok := g.(MutableDataGraph); ok {
		Suite(&MutableDataGraphSuite{fact, directed})
	}

	return true
}

/////////////////////////////////////////////////////////////////////
//
// SUITES
//
/////////////////////////////////////////////////////////////////////

/* GraphSuite - tests for non-mutable graph methods */

type GraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *GraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *GraphSuite) TestHasVertex(c *C) {
	g := s.Factory(graphFixtures["2e3v"])
	c.Assert(g.HasVertex("qux"), Equals, false)
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *GraphSuite) TestHasEdge(c *C) {
	g := s.Factory(graphFixtures["2e3v"])
	c.Assert(g.HasEdge(edgeSet[0]), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{"qux", "quark"}), Equals, false)
}

func (s *GraphSuite) TestEachVertex(c *C) {
	g := s.Factory(graphFixtures["2e3v"])

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

func (s *GraphSuite) TestEachVertexTermination(c *C) {
	g := s.Factory(graphFixtures["2e3v"])

	var hit int
	g.EachVertex(func(v Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachEdge(c *C) {
	g := s.Factory(graphFixtures["2e3v"])

	var hit int
	g.EachEdge(func(e Edge) (terminate bool) {
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
}

func (s *GraphSuite) TestEachEdgeTermination(c *C) {
	g := s.Factory(graphFixtures["2e3v"])

	var hit int
	g.EachEdge(func(e Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachAdjacentTo(c *C) {
	g := s.Factory(graphFixtures["2e3v"])

	vset := set.NewNonTS()
	var hit int
	g.EachAdjacentTo("bar", func(adj Vertex) (terminate bool) {
		hit++
		vset.Add(adj)
		return
	})

	c.Assert(vset.Has("foo"), Equals, true)
	c.Assert(vset.Has("bar"), Equals, false)
	c.Assert(vset.Has("baz"), Equals, true)
	c.Assert(hit, Equals, 2)
}

func (s *GraphSuite) TestEachAdjacentToTermination(c *C) {
	g := s.Factory(graphFixtures["3e4v"])

	var hit int
	g.EachAdjacentTo("foo", func(adjacent Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachEdgeIncidentTo(c *C) {
	g := s.Factory(graphFixtures["2e3v"])

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

func (s *GraphSuite) TestEachEdgeIncidentToTermination(c *C) {
	g := s.Factory(graphFixtures["2e3v"])

	var hit int
	g.EachEdgeIncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestDegreeOf(c *C) {
	g := s.Factory(graphFixtures["3e5v1i"])

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

func (s *GraphSuite) TestSize(c *C) {
	c.Assert(s.Factory(NullGraph).Size(), Equals, 0)
	c.Assert(s.Factory(graphFixtures["2e3v"]).Size(), Equals, 2)
}

func (s *GraphSuite) TestOrder(c *C) {
	c.Assert(s.Factory(NullGraph).Order(), Equals, 0)
	c.Assert(s.Factory(graphFixtures["2e3v"]).Order(), Equals, 3)
}

/* DirectedGraphSuite - tests for directed graph methods */

type DirectedGraphSuite struct {
	Factory func(GraphSource) Graph
}

func (s *DirectedGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DirectedGraphSuite) TestTranspose(c *C) {
	g := s.Factory(graphFixtures["2e3v"]).(DirectedGraph)

	g2 := g.Transpose()

	c.Assert(g2.HasEdge(edgeSet[0].(BaseEdge).swap()), Equals, true)
	c.Assert(g2.HasEdge(edgeSet[1].(BaseEdge).swap()), Equals, true)

	c.Assert(g2.HasEdge(edgeSet[0].(BaseEdge)), Equals, false)
	c.Assert(g2.HasEdge(edgeSet[1].(BaseEdge)), Equals, false)
}

func (s *DirectedGraphSuite) TestOutDegreeOf(c *C) {
	g := s.Factory(graphFixtures["3e5v1i"]).(DirectedGraph)

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
	g := s.Factory(graphFixtures["3e5v1i"]).(DirectedGraph)

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
	g := s.Factory(graphFixtures["arctest"]).(DirectedGraph)

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
	g := s.Factory(graphFixtures["arctest"]).(DirectedGraph)

	var hit int
	g.EachArcTo("baz", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DirectedGraphSuite) TestEachArcFrom(c *C) {
	g := s.Factory(graphFixtures["arctest"]).(DirectedGraph)

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
	g := s.Factory(graphFixtures["arctest"]).(DirectedGraph)

	var hit int
	g.EachArcFrom("foo", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

/* SimpleGraphSuite - tests for simple graph methods */

type SimpleGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *SimpleGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *SimpleGraphSuite) TestDensity(c *C) {
	c.Assert(math.IsNaN(s.Factory(NullGraph).(SimpleGraph).Density()), Equals, true)

	g := s.Factory(graphFixtures["pair"]).(SimpleGraph)
	if s.Directed {
		c.Assert(g.Density(), Equals, float64(0.5))
	} else {
		c.Assert(g.Density(), Equals, float64(1))
	}

	g = s.Factory(graphFixtures["2e3v"]).(SimpleGraph)
	if s.Directed {
		c.Assert(g.Density(), Equals, float64(2)/float64(6))
	} else {
		c.Assert(g.Density(), Equals, float64(2)/float64(3))
	}
}

/* MutableGraphSuite - tests for mutable graph methods */

type MutableGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *MutableGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *MutableGraphSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)

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
	g := s.Factory(NullGraph).(MutableGraph)

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableGraphSuite) TestRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)
	g.AddEdges(&BaseEdge{1, 2})

	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed)

	// Now test removal
	g.RemoveEdges(&BaseEdge{1, 2})
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, false)
}

func (s *MutableGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)

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
	g := s.Factory(NullGraph).(MutableGraph)

	g.AddEdges(&BaseEdge{1, 2}, &BaseEdge{2, 3}, &BaseEdge{4, 1})
	g.RemoveVertex(1)

	c.Assert(g.Size(), Equals, 1)
}

/* WeightedGraphSuite - tests for weighted graphs */

type WeightedGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *WeightedGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *WeightedGraphSuite) TestEachEdge(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachEdge() iterator actually do implement WeightedEdge.
	g := s.Factory(graphFixtures["w-2e3v"]).(WeightedGraph)

	var we WeightedEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *WeightedGraphSuite) TestHasWeightedEdge(c *C) {
	g := s.Factory(graphFixtures["w-2e3v"]).(WeightedGraph)

	// TODO figure out how to meaningfully test undirected graphs' logic here
	c.Assert(g.HasWeightedEdge(baseWeightedEdgeSet[0]), Equals, true)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 1}), Equals, false) // wrong weight
}

/* MutableWeightedGraphSuite - tests for mutable weighted graphs */

type MutableWeightedGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *MutableWeightedGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *MutableWeightedGraphSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph).(MutableWeightedGraph)

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
	g := s.Factory(NullGraph).(MutableWeightedGraph)

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableWeightedGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Factory(NullGraph).(MutableWeightedGraph)

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableWeightedGraphSuite) TestRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableWeightedGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableWeightedGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableWeightedGraph)
	g.AddEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5.23})

	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed)

	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 5.23}), Equals, true)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 3}), Equals, false)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 1}, 5.23}), Equals, !s.Directed)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{2, 1}, -3.22771}), Equals, false)

	// Now test removal
	g.RemoveEdges(BaseWeightedEdge{BaseEdge{1, 2}, 5.23})
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasWeightedEdge(BaseWeightedEdge{BaseEdge{1, 2}, 5.23}), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableWeightedGraph)
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

/* LabeledGraphSuite - tests for labeled graphs */

type LabeledGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *LabeledGraphSuite) TestEachEdge(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachEdge() iterator actually do implement LabeledEdge.
	g := s.Factory(graphFixtures["l-2e3v"]).(LabeledGraph)

	var we LabeledEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *LabeledGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *LabeledGraphSuite) TestHasLabeledEdge(c *C) {
	g := s.Factory(graphFixtures["l-2e3v"]).(LabeledGraph)

	// TODO figure out how to meaningfully test undirected graphs' logic here
	c.Assert(g.HasLabeledEdge(baseLabeledEdgeSet[0]), Equals, true)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{1, 2}, "qux"}), Equals, false) // wrong label
}

/* MutableLabeledGraphSuite - tests for mutable labeled graphs */

type MutableLabeledGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *MutableLabeledGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *MutableLabeledGraphSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)

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

func (s *MutableLabeledGraphSuite) TestEnsureVertex(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableLabeledGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableLabeledGraphSuite) TestRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableLabeledGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableLabeledGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)
	g.AddEdges(BaseLabeledEdge{BaseEdge{1, 2}, "foo"})

	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed)

	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{1, 2}, "foo"}), Equals, true)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{1, 2}, "baz"}), Equals, false)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{2, 1}, "foo"}), Equals, !s.Directed)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{2, 1}, "quark"}), Equals, false)

	// Now test removal
	g.RemoveEdges(BaseLabeledEdge{BaseEdge{1, 2}, "foo"})
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{1, 2}, "foo"}), Equals, false)
}

func (s *MutableLabeledGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)
	g.AddEdges(BaseLabeledEdge{BaseEdge{1, 2}, "foo"}, BaseLabeledEdge{BaseEdge{2, 3}, "bar"})

	// Basic edge tests first
	// We test both Has*Edge() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Edge() method.
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed) // only if undirected
	c.Assert(g.HasEdge(BaseEdge{3, 2}), Equals, !s.Directed) // only if undirected

	// Now labeled edge tests
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{1, 2}, "foo"}), Equals, true)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{1, 2}, "baz"}), Equals, false) // wrong label
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{2, 1}, "foo"}), Equals, !s.Directed)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{2, 1}, "baz"}), Equals, false) // wrong label
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{2, 3}, "bar"}), Equals, true)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{2, 3}, "qux"}), Equals, false) // wrong label
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{3, 2}, "bar"}), Equals, !s.Directed)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{3, 2}, "qux"}), Equals, false) // wrong label

	// Now test removal
	g.RemoveEdges(BaseLabeledEdge{BaseEdge{1, 2}, "foo"}, BaseLabeledEdge{BaseEdge{2, 3}, "bar"})
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{1, 2}, "foo"}), Equals, false)
	c.Assert(g.HasLabeledEdge(BaseLabeledEdge{BaseEdge{2, 3}, "bar"}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, false)
}

/* DataGraphSuite - tests for labeled graphs */

type DataGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *DataGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DataGraphSuite) TestEachEdge(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachEdge() iterator actually do implement DataEdge.
	g := s.Factory(graphFixtures["p-2e3v"]).(DataGraph)

	var we DataEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *DataGraphSuite) TestHasDataEdge(c *C) {
	g := s.Factory(graphFixtures["p-2e3v"]).(DataGraph)

	// TODO figure out how to meaningfully test undirected graphs' logic here
	c.Assert(g.HasDataEdge(baseDataEdgeSet[1]), Equals, true)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{1, 2}, "qux"}), Equals, false) // wrong label
}

/* MutableDataGraphSuite - tests for mutable labeled graphs */

type MutableDataGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *MutableDataGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *MutableDataGraphSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)

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

func (s *MutableDataGraphSuite) TestEnsureVertex(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)

	g.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *MutableDataGraphSuite) TestMultiEnsureVertex(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)

	g.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *MutableDataGraphSuite) TestRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *MutableDataGraphSuite) TestMultiRemoveVertex(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)

	g.EnsureVertex("bar", "baz")
	g.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

func (s *MutableDataGraphSuite) TestAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)
	g.AddEdges(BaseDataEdge{BaseEdge{1, 2}, "foo"})

	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed)

	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{1, 2}, "foo"}), Equals, true)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{1, 2}, "baz"}), Equals, false)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{2, 1}, "foo"}), Equals, !s.Directed)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{2, 1}, "quark"}), Equals, false)

	// Now test removal
	g.RemoveEdges(BaseDataEdge{BaseEdge{1, 2}, "foo"})
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{1, 2}, "foo"}), Equals, false)
}

func (s *MutableDataGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)
	g.AddEdges(BaseDataEdge{BaseEdge{1, 2}, "foo"}, BaseDataEdge{BaseEdge{2, 3}, struct{ a int }{a: 2}})

	// Basic edge tests first
	// We test both Has*Edge() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Edge() method.
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, true)
	c.Assert(g.HasEdge(BaseEdge{2, 1}), Equals, !s.Directed) // only if undirected
	c.Assert(g.HasEdge(BaseEdge{3, 2}), Equals, !s.Directed) // only if undirected

	// Now labeled edge tests
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{1, 2}, "foo"}), Equals, true)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{1, 2}, "baz"}), Equals, false) // wrong label
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{2, 1}, "foo"}), Equals, !s.Directed)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{2, 1}, "baz"}), Equals, false) // wrong label
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{2, 3}, struct{ a int }{a: 2}}), Equals, true)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{2, 3}, "qux"}), Equals, false) // wrong label
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{3, 2}, struct{ a int }{a: 2}}), Equals, !s.Directed)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{3, 2}, "qux"}), Equals, false) // wrong label

	// Now test removal
	g.RemoveEdges(BaseDataEdge{BaseEdge{1, 2}, "foo"}, BaseDataEdge{BaseEdge{2, 3}, struct{ a int }{a: 2}})
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{1, 2}, "foo"}), Equals, false)
	c.Assert(g.HasDataEdge(BaseDataEdge{BaseEdge{2, 3}, struct{ a int }{a: 2}}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{1, 2}), Equals, false)
	c.Assert(g.HasEdge(BaseEdge{2, 3}), Equals, false)
}
