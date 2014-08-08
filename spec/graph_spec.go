package spec

import (
	"fmt"
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

type loopArc struct {
	v Vertex
}

func (e loopArc) Both() (Vertex, Vertex) {
	return e.v, e.v
}

func (e loopArc) Source() Vertex {
	return e.v
}

func (e loopArc) Target() Vertex {
	return e.v
}

type loopArcList []Arc

func (el loopArcList) EachVertex(fn VertexStep) {
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

func (el loopArcList) EachArc(fn ArcStep) {
	for _, e := range el {
		if _, ok := e.(loopArc); !ok {
			if fn(e) {
				return
			}
		}
	}
}

func (el loopArcList) EachEdge(fn EdgeStep) {
	for _, e := range el {
		if _, ok := e.(loopArc); !ok {
			if fn(e) {
				return
			}
		}
	}
}

var GraphFixtures = map[string]GraphSource{
	// TODO improve naming basis/patterns for these
	"arctest": ArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
		NewArc("foo", "qux"),
		NewArc("qux", "bar"),
	},
	"pair": ArcList{
		NewArc(1, 2),
	},
	"2e3v": ArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
	},
	"3e4v": ArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
		NewArc("foo", "qux"),
	},
	"3e5v1i": loopArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
		NewArc("foo", "qux"),
		loopArc{"isolate"},
	},
	"w-2e3v": WeightedArcList{
		NewWeightedArc(1, 2, 5.23),
		NewWeightedArc(2, 3, 5.821),
	},
	"l-2e3v": LabeledArcList{
		NewLabeledArc(1, 2, "foo"),
		NewLabeledArc(2, 3, "bar"),
	},
	"d-2e3v": DataArcList{
		NewDataArc(1, 2, "foo"),
		NewDataArc(2, 3, struct{ a int }{a: 2}),
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
func Swap(a Arc) Arc {
	return NewArc(a.Target(), a.Source())
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

	if _, ok := g.(VertexSetMutator); ok {
		Suite(&VertexSetMutatorSuite{fact})
	}

	if _, ok := g.(EdgeSetMutator); ok {
		Suite(&EdgeSetMutatorSuite{fact})
	}

	if _, ok := g.(ArcSetMutator); ok {
		Suite(&ArcSetMutatorSuite{fact})
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
	c.Assert(g.HasWeightedEdge(GraphFixtures["w-2e3v"].(WeightedArcList)[0].(WeightedArc)), Equals, true)
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
	c.Assert(g.HasLabeledEdge(GraphFixtures["l-2e3v"].(LabeledArcList)[0].(LabeledArc)), Equals, true)
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
	c.Assert(g.HasDataEdge(GraphFixtures["d-2e3v"].(DataArcList)[1].(DataArc)), Equals, true)
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

