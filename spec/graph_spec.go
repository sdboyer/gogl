package spec

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

/////////////////////////////////////////////////////////////////////
//
// GRAPH FIXTURES
//
/////////////////////////////////////////////////////////////////////

type loopEdge struct {
	v Vertex
}

func (e loopEdge) Both() (Vertex, Vertex) {
	return e.v, e.v
}

func (e loopEdge) Source() Vertex {
	return e.v
}

func (e loopEdge) Target() Vertex {
	return e.v
}

type loopEdgeList []Edge

func (el loopEdgeList) EachVertex(fn VertexStep) {
	set := set.NewNonTS()

	for _, e := range el {
		set.Add(e.Both())
	}

	for _, v := range set.List() {
		if fn(v) {
			return
		}
	}
}

func (el loopEdgeList) EachEdge(fn EdgeStep) {
	for _, e := range el {
		if _, ok := e.(loopEdge); !ok {
			if fn(e) {
				return
			}
		}
	}
}

var GraphFixtures = map[string]GraphSource{
	// TODO improve naming basis/patterns for these
	"arctest": EdgeList{
		NewEdge("foo", "bar"),
		NewEdge("bar", "baz"),
		NewEdge("foo", "qux"),
		NewEdge("qux", "bar"),
	},
	"pair": EdgeList{
		NewEdge(1, 2),
	},
	"2e3v": EdgeList{
		NewEdge("foo", "bar"),
		NewEdge("bar", "baz"),
	},
	"3e4v": EdgeList{
		NewEdge("foo", "bar"),
		NewEdge("bar", "baz"),
		NewEdge("foo", "qux"),
	},
	"3e5v1i": loopEdgeList{
		NewEdge("foo", "bar"),
		NewEdge("bar", "baz"),
		NewEdge("foo", "qux"),
		loopEdge{"isolate"},
	},
	"w-2e3v": WeightedEdgeList{
		NewWeightedEdge(1, 2, 5.23),
		NewWeightedEdge(2, 3, 5.821),
	},
	"l-2e3v": LabeledEdgeList{
		NewLabeledEdge(1, 2, "foo"),
		NewLabeledEdge(2, 3, "bar"),
	},
	"d-2e3v": DataEdgeList{
		NewDataEdge(1, 2, "foo"),
		NewDataEdge(2, 3, struct{ a int }{a: 2}),
	},
}

/////////////////////////////////////////////////////////////////////
//
// HELPERS
//
/////////////////////////////////////////////////////////////////////

// Hook gocheck into the go test runner
func TestHookup(t *testing.T) { TestingT(t) }

// Returns an arc with the directionality swapped.
func Swap(e Edge) Edge {
	u, v := e.Both()
	return NewEdge(v, u)
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

	if _, ok := g.(Digraph); ok {
		directed = true
		Suite(&DigraphSuite{Factory: fact})
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
	g := s.Factory(GraphFixtures["2e3v"])
	c.Assert(g.HasVertex("qux"), Equals, false)
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *GraphSuite) TestHasEdge(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])
	c.Assert(g.HasEdge(GraphFixtures["2e3v"].(EdgeList)[0]), Equals, true)
	c.Assert(g.HasEdge(NewEdge("qux", "quark")), Equals, false)
}

func (s *GraphSuite) TestEachVertex(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

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
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.EachVertex(func(v Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachEdge(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.EachEdge(func(e Edge) (terminate bool) {
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
}

func (s *GraphSuite) TestEachEdgeTermination(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.EachEdge(func(e Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachAdjacentTo(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

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
	g := s.Factory(GraphFixtures["3e4v"])

	var hit int
	g.EachAdjacentTo("foo", func(adjacent Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachEdgeIncidentTo(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	flipset := []Edge{
		Swap(GraphFixtures["2e3v"].(EdgeList)[0]),
		Swap(GraphFixtures["2e3v"].(EdgeList)[1]),
	}

	eset := set.NewNonTS()
	var hit int
	g.EachEdgeIncidentTo("foo", func(e Edge) (terminate bool) {
		hit++
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewEdge(e.Both()))
		return
	})

	c.Assert(hit, Equals, 1)
	if s.Directed {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[1]), Equals, false)
	} else {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[0]) != eset.Has(flipset[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[1]) != eset.Has(flipset[1]), Equals, false)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[1]), Equals, false)
	}

	eset = set.NewNonTS()
	g.EachEdgeIncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewEdge(e.Both()))
		return
	})

	c.Assert(hit, Equals, 3)
	if s.Directed {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[1]), Equals, true)
	} else {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[0]) != eset.Has(flipset[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[1]) != eset.Has(flipset[1]), Equals, true)
	}
}

func (s *GraphSuite) TestEachEdgeIncidentToTermination(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.EachEdgeIncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestDegreeOf(c *C) {
	g := s.Factory(GraphFixtures["3e5v1i"])

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

/* DigraphSuite - tests for directed graph methods */

type DigraphSuite struct {
	Factory func(GraphSource) Graph
}

func (s *DigraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DigraphSuite) TestTranspose(c *C) {
	g := s.Factory(GraphFixtures["2e3v"]).(Digraph)

	g2 := g.Transpose()

	c.Assert(g2.HasEdge(Swap(GraphFixtures["2e3v"].(EdgeList)[0])), Equals, true)
	c.Assert(g2.HasEdge(Swap(GraphFixtures["2e3v"].(EdgeList)[1])), Equals, true)

	c.Assert(g2.HasEdge(GraphFixtures["2e3v"].(EdgeList)[0]), Equals, false)
	c.Assert(g2.HasEdge(GraphFixtures["2e3v"].(EdgeList)[1]), Equals, false)
}

func (s *DigraphSuite) TestOutDegreeOf(c *C) {
	g := s.Factory(GraphFixtures["3e5v1i"]).(Digraph)

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

func (s *DigraphSuite) TestInDegreeOf(c *C) {
	g := s.Factory(GraphFixtures["3e5v1i"]).(Digraph)

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

func (s *DigraphSuite) TestEachArcTo(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	var hit int
	g.EachArcTo("foo", func(e Arc) (terminate bool) {
		c.Error("Vertex 'foo' should have no in-edges")
		c.FailNow()
		return
	})

	g.EachArcTo("bar", func(e Arc) (terminate bool) {
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewArc(e.Both()))
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[0]), Equals, true)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[1]), Equals, false)
	c.Assert(eset.Has(NewEdge("qux", "bar")), Equals, true)
}

func (s *DigraphSuite) TestEachArcToTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachArcTo("baz", func(e Arc) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DigraphSuite) TestEachPredecessorOf(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	g.EachPredecessorOf("foo", func(v Vertex) (terminate bool) {
		c.Error("Vertex 'foo' should have no predecessors")
		c.FailNow()
		return
	})

	g.EachPredecessorOf("bar", func(v Vertex) (terminate bool) {
		eset.Add(v)
		return
	})

	c.Assert(eset.Size(), Equals, 2)
	c.Assert(eset.Has("foo"), Equals, true)
	c.Assert(eset.Has("qux"), Equals, true)

}

func (s *DigraphSuite) TestEachPredecessorOfTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachPredecessorOf("baz", func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DigraphSuite) TestEachArcFrom(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	var hit int
	g.EachArcFrom("baz", func(e Arc) (terminate bool) {
		c.Error("Vertex 'baz' should have no out-edges")
		c.FailNow()
		return
	})

	g.EachArcFrom("foo", func(e Arc) (terminate bool) {
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewArc(e.Both()))
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[0]), Equals, true)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(EdgeList)[1]), Equals, false)
	c.Assert(eset.Has(NewEdge("foo", "qux")), Equals, true)
}

func (s *DigraphSuite) TestEachArcFromTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachArcFrom("foo", func(e Arc) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DigraphSuite) TestEachSuccessorOf(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	g.EachSuccessorOf("baz", func(v Vertex) (terminate bool) {
		c.Error("Vertex 'foo' should have no successors")
		c.FailNow()
		return
	})

	g.EachSuccessorOf("foo", func(v Vertex) (terminate bool) {
		eset.Add(v)
		return
	})

	c.Assert(eset.Size(), Equals, 2)
	c.Assert(eset.Has("qux"), Equals, true)
	c.Assert(eset.Has("bar"), Equals, true)

}

func (s *DigraphSuite) TestEachSuccessorOfTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachSuccessorOf("foo", func(v Vertex) (terminate bool) {
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

	g := s.Factory(GraphFixtures["pair"]).(SimpleGraph)
	if s.Directed {
		c.Assert(g.Density(), Equals, float64(0.5))
	} else {
		c.Assert(g.Density(), Equals, float64(1))
	}

	g = s.Factory(GraphFixtures["2e3v"]).(SimpleGraph)
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
	c.Assert(Order(g), Equals, 0)

	g.RemoveVertex()
	c.Assert(Order(g), Equals, 0)

	g.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	g.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
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
	g.AddEdges(NewEdge(1, 2))

	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed)

	// Now test removal
	g.RemoveEdges(NewEdge(1, 2))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, false)
}

func (s *MutableGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)

	g.AddEdges(NewEdge(1, 2), NewEdge(2, 3))

	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed)
	c.Assert(g.HasEdge(NewEdge(3, 2)), Equals, !s.Directed)

	// Now test removal
	g.RemoveEdges(NewEdge(1, 2), NewEdge(2, 3))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
}

// Checks to ensure that removal works for both in-edges and out-edges.
func (s *MutableGraphSuite) TestVertexRemovalAlsoRemovesConnectedEdges(c *C) {
	g := s.Factory(NullGraph).(MutableGraph)

	g.AddEdges(NewEdge(1, 2), NewEdge(2, 3), NewEdge(4, 1))
	g.RemoveVertex(1)

	c.Assert(Size(g), Equals, 1)
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
	g := s.Factory(GraphFixtures["w-2e3v"]).(WeightedGraph)

	var we WeightedEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *WeightedGraphSuite) TestHasWeightedEdge(c *C) {
	g := s.Factory(GraphFixtures["w-2e3v"]).(WeightedGraph)

	// TODO figure out how to meaningfully test undirected graphs' logic here
	c.Assert(g.HasWeightedEdge(GraphFixtures["w-2e3v"].(WeightedEdgeList)[0]), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 1)), Equals, false) // wrong weight
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
	c.Assert(Order(g), Equals, 0)

	g.RemoveVertex()
	c.Assert(Order(g), Equals, 0)

	g.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	g.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
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
	g.AddEdges(NewWeightedEdge(1, 2, 5.23))

	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed)

	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5.23)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 3)), Equals, false)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, 5.23)), Equals, !s.Directed)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, -3.22771)), Equals, false)

	// Now test removal
	g.RemoveEdges(NewWeightedEdge(1, 2, 5.23))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5.23)), Equals, false)
}

func (s *MutableWeightedGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableWeightedGraph)
	g.AddEdges(NewWeightedEdge(1, 2, 5), NewWeightedEdge(2, 3, -5))

	// Basic edge tests first
	// We test both Has*Edge() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Edge() method.
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed) // only if undirected
	c.Assert(g.HasEdge(NewEdge(3, 2)), Equals, !s.Directed) // only if undirected

	// Now weighted edge tests
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 3)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, 5)), Equals, !s.Directed)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, 3)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 3, -5)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 3, 1)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(3, 2, -5)), Equals, !s.Directed)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(3, 2, 1)), Equals, false) // wrong weight

	// Now test removal
	g.RemoveEdges(NewWeightedEdge(1, 2, 5), NewWeightedEdge(2, 3, -5))
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5)), Equals, false)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 3, -5)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
}

/* LabeledGraphSuite - tests for labeled graphs */

type LabeledGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *LabeledGraphSuite) TestEachEdge(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachEdge() iterator actually do implement LabeledEdge.
	g := s.Factory(GraphFixtures["l-2e3v"]).(LabeledGraph)

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
	g := s.Factory(GraphFixtures["l-2e3v"]).(LabeledGraph)

	// TODO figure out how to meaningfully test undirected graphs' logic here
	c.Assert(g.HasLabeledEdge(GraphFixtures["l-2e3v"].(LabeledEdgeList)[0]), Equals, true)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "qux")), Equals, false) // wrong label
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
	c.Assert(Order(g), Equals, 0)

	g.RemoveVertex()
	c.Assert(Order(g), Equals, 0)

	g.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	g.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
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
	g.AddEdges(NewLabeledEdge(1, 2, "foo"))

	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed)

	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "baz")), Equals, false)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 1, "foo")), Equals, !s.Directed)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 1, "quark")), Equals, false)

	// Now test removal
	g.RemoveEdges(NewLabeledEdge(1, 2, "foo"))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, false)
}

func (s *MutableLabeledGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableLabeledGraph)
	g.AddEdges(NewLabeledEdge(1, 2, "foo"), NewLabeledEdge(2, 3, "bar"))

	// Basic edge tests first
	// We test both Has*Edge() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Edge() method.
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed) // only if undirected
	c.Assert(g.HasEdge(NewEdge(3, 2)), Equals, !s.Directed) // only if undirected

	// Now labeled edge tests
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "baz")), Equals, false) // wrong label
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 1, "foo")), Equals, !s.Directed)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 1, "baz")), Equals, false) // wrong label
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 3, "bar")), Equals, true)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 3, "qux")), Equals, false) // wrong label
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(3, 2, "bar")), Equals, !s.Directed)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(3, 2, "qux")), Equals, false) // wrong label

	// Now test removal
	g.RemoveEdges(NewLabeledEdge(1, 2, "foo"), NewLabeledEdge(2, 3, "bar"))
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, false)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 3, "bar")), Equals, false)
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
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
	g := s.Factory(GraphFixtures["d-2e3v"]).(DataGraph)

	var we DataEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *DataGraphSuite) TestHasDataEdge(c *C) {
	g := s.Factory(GraphFixtures["d-2e3v"]).(DataGraph)

	// TODO figure out how to meaningfully test undirected graphs' logic here
	c.Assert(g.HasDataEdge(GraphFixtures["d-2e3v"].(DataEdgeList)[1]), Equals, true)
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "qux")), Equals, false) // wrong label
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
	c.Assert(Order(g), Equals, 0)

	g.RemoveVertex()
	c.Assert(Order(g), Equals, 0)

	g.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	g.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
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
	g.AddEdges(NewDataEdge(1, 2, "foo"))

	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed)

	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "baz")), Equals, false)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 1, "foo")), Equals, !s.Directed)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 1, "quark")), Equals, false)

	// Now test removal
	g.RemoveEdges(NewDataEdge(1, 2, "foo"))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, false)
}

func (s *MutableDataGraphSuite) TestMultiAddAndRemoveEdge(c *C) {
	g := s.Factory(NullGraph).(MutableDataGraph)
	g.AddEdges(NewDataEdge(1, 2, "foo"), NewDataEdge(2, 3, struct{ a int }{a: 2}))

	// Basic edge tests first
	// We test both Has*Edge() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Edge() method.
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, !s.Directed) // only if undirected
	c.Assert(g.HasEdge(NewEdge(3, 2)), Equals, !s.Directed) // only if undirected

	// Now labeled edge tests
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "baz")), Equals, false) // wrong label
	c.Assert(g.HasDataEdge(NewDataEdge(2, 1, "foo")), Equals, !s.Directed)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 1, "baz")), Equals, false) // wrong label
	c.Assert(g.HasDataEdge(NewDataEdge(2, 3, struct{ a int }{a: 2})), Equals, true)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 3, "qux")), Equals, false) // wrong label
	c.Assert(g.HasDataEdge(NewDataEdge(3, 2, struct{ a int }{a: 2})), Equals, !s.Directed)
	c.Assert(g.HasDataEdge(NewDataEdge(3, 2, "qux")), Equals, false) // wrong label

	// Now test removal
	g.RemoveEdges(NewDataEdge(1, 2, "foo"), NewDataEdge(2, 3, struct{ a int }{a: 2}))
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, false)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 3, struct{ a int }{a: 2})), Equals, false)
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
}

/* Counting suites - tests for Size() and Order() */

type OrderSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *OrderSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *OrderSuite) TestOrder(c *C) {
	c.Assert(s.Factory(NullGraph).(VertexCounter).Order(), Equals, 0)
	c.Assert(s.Factory(GraphFixtures["2e3v"]).(VertexCounter).Order(), Equals, 3)
}

type SizeSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *SizeSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *SizeSuite) TestSize(c *C) {
	c.Assert(s.Factory(NullGraph).(EdgeCounter).Size(), Equals, 0)
	c.Assert(s.Factory(GraphFixtures["2e3v"]).(EdgeCounter).Size(), Equals, 2)
}
