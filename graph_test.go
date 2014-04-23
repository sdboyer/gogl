package gogl

import (
	"fmt"
	"reflect"

	"gopkg.in/fatih/set.v0"
	. "launchpad.net/gocheck"
)

var _ = fmt.Println

// This function automatically sets up suites of black box unit tests for
// graphs by determining which gogl interfaces they implement.
//
// Passing a graph to this method for testing is the most official way to
// determine whether or not it complies with not just the interfaces, but also
// the graph semantics defined by gogl.
func SetUpSimpleGraphTests(g Graph) bool {
	gf := &GraphTestFactory{g}

	_, directed := g.(DirectedGraph)

	// Set up the basic Graph suite unconditionally
	Suite(&GraphSuite{Graph: g, Factory: gf, Directed: directed})

	return true
}

// Set up suites for all of gogl's graphs.
var _ = SetUpSimpleGraphTests(NewDirected())
var _ = SetUpSimpleGraphTests(NewUndirected())
var _ = SetUpSimpleGraphTests(NewWeightedDirected())
var _ = SetUpSimpleGraphTests(NewWeightedUndirected())
var _ = SetUpSimpleGraphTests(NewLabeledDirected())
var _ = SetUpSimpleGraphTests(NewLabeledUndirected())
var _ = SetUpSimpleGraphTests(NewPropertyDirected())
var _ = SetUpSimpleGraphTests(NewPropertyUndirected())

/* The GraphTestFactory - this generates graph instances for the tests. */

type GraphTestFactory struct {
	sourceGraph Graph
}

func (gf *GraphTestFactory) create() interface{} {
	return reflect.New(reflect.Indirect(reflect.ValueOf(gf.sourceGraph)).Type()).Interface()
}

func (gf *GraphTestFactory) graphFromEdges(edges ...Edge) Graph {
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
	} else if mlg, ok := base.(MutableLabeledGraph); ok {
		labeled_edges := make([]LabeledEdge, 0, len(edges))
		for _, edge := range edges {
			labeled_edges = append(labeled_edges, BaseLabeledEdge{BaseEdge{edge.Source(), edge.Target()}, ""})
		}
		mlg.AddEdges(labeled_edges...)
	} else if mpg, ok := base.(MutablePropertyGraph); ok {
		property_edges := make([]PropertyEdge, 0, len(edges))
		for _, edge := range edges {
			property_edges = append(property_edges, BasePropertyEdge{BaseEdge{edge.Source(), edge.Target()}, ""})
		}
		mpg.AddEdges(property_edges...)
	} else {
		panic("Until GraphInitializers are made to work properly, all graphs have to be mutable for this testing harness to work.")
	}

	return base.(Graph)

}

func (gf *GraphTestFactory) CreateEmptyGraph() Graph {
	return gf.create().(Graph)
}

func (gf *GraphTestFactory) CreateGraphFromEdges(edges ...Edge) Graph {
	return gf.graphFromEdges(edges...)
}

func (gf *GraphTestFactory) CreateDirectedGraphFromEdges(edges ...Edge) DirectedGraph {
	return gf.graphFromEdges(edges...).(DirectedGraph)
}

func (gf *GraphTestFactory) CreateEmptySimpleGraph() SimpleGraph {
	return gf.create().(SimpleGraph)
}

func (gf *GraphTestFactory) CreateSimpleGraphFromEdges(edges ...Edge) SimpleGraph {
	return gf.graphFromEdges(edges...).(SimpleGraph)
}

func (gf *GraphTestFactory) CreateMutableGraph() MutableGraph {
	return gf.create().(MutableGraph)
}

func (gf *GraphTestFactory) CreateWeightedGraphFromEdges(edges ...WeightedEdge) WeightedGraph {
	base := gf.create()
	if mwg, ok := base.(MutableWeightedGraph); ok {
		mwg.AddEdges(edges...)
		return mwg
	}

	panic("Until GraphInitializers are made to work properly, all graphs have to be mutable for this testing harness to work.")
}

func (gf *GraphTestFactory) CreateEmptyWeightedGraph() WeightedGraph {
	return gf.create().(WeightedGraph)
}

func (gf *GraphTestFactory) CreateMutableWeightedGraph() MutableWeightedGraph {
	return gf.create().(MutableWeightedGraph)
}

func (gf *GraphTestFactory) CreateLabeledGraphFromEdges(edges ...LabeledEdge) LabeledGraph {
	base := gf.create()
	if mlg, ok := base.(MutableLabeledGraph); ok {
		mlg.AddEdges(edges...)
		return mlg
	}

	panic("Until GraphInitializers are made to work properly, all graphs have to be mutable for this testing harness to work.")
}

func (gf *GraphTestFactory) CreateEmptyLabeledGraph() LabeledGraph {
	return gf.create().(LabeledGraph)
}

func (gf *GraphTestFactory) CreateMutableLabeledGraph() MutableLabeledGraph {
	return gf.create().(MutableLabeledGraph)
}

func (gf *GraphTestFactory) CreatePropertyGraphFromEdges(edges ...PropertyEdge) PropertyGraph {
	base := gf.create()
	if mpg, ok := base.(MutablePropertyGraph); ok {
		mpg.AddEdges(edges...)
		return mpg
	}

	panic("Until GraphInitializers are made to work properly, all graphs have to be mutable for this testing harness to work.")
}

func (gf *GraphTestFactory) CreateEmptyPropertyGraph() PropertyGraph {
	return gf.create().(PropertyGraph)
}

func (gf *GraphTestFactory) CreateMutablePropertyGraph() MutablePropertyGraph {
	return gf.create().(MutablePropertyGraph)
}

/* Factory interfaces for tests */

type GraphCreator interface {
	CreateEmptyGraph() Graph
	CreateGraphFromEdges(edges ...Edge) Graph
}

type SimpleGraphCreator interface {
	CreateEmptySimpleGraph() SimpleGraph
	CreateSimpleGraphFromEdges(edges ...Edge) SimpleGraph
}

type DirectedGraphCreator interface {
	CreateDirectedGraphFromEdges(edges ...Edge) DirectedGraph
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

type LabeledGraphCreator interface {
	CreateEmptyLabeledGraph() LabeledGraph
	CreateLabeledGraphFromEdges(edges ...LabeledEdge) LabeledGraph
}

type MutableLabeledGraphCreator interface {
	CreateMutableLabeledGraph() MutableLabeledGraph
}

type PropertyGraphCreator interface {
	CreateEmptyPropertyGraph() PropertyGraph
	CreatePropertyGraphFromEdges(edges ...PropertyEdge) PropertyGraph
}

type MutablePropertyGraphCreator interface {
	CreateMutablePropertyGraph() MutablePropertyGraph
}

/* GraphSuite - tests for non-mutable graph methods */

type GraphSuite struct {
	Graph    Graph
	Factory  GraphCreator
	Directed bool
}

func (s *GraphSuite) SetUpTest(c *C) {
	s.Graph = s.Factory.CreateGraphFromEdges(edgeSet...)
}

func (s *GraphSuite) TestHasVertex(c *C) {
	c.Assert(s.Graph.HasVertex("qux"), Equals, false)
	c.Assert(s.Graph.HasVertex("foo"), Equals, true)
}

func (s *GraphSuite) TestHasEdge(c *C) {
	c.Assert(s.Graph.HasEdge(edgeSet[0]), Equals, true)
	c.Assert(s.Graph.HasEdge(BaseEdge{"qux", "quark"}), Equals, false)
}

func (s *GraphSuite) TestEachVertex(c *C) {
	g := s.Factory.CreateGraphFromEdges(edgeSet...)

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
	g := s.Factory.CreateGraphFromEdges(edgeSet...)

	var hit int
	g.EachVertex(func(v Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachEdge(c *C) {
	g := s.Factory.CreateGraphFromEdges(edgeSet...)

	var hit int
	g.EachEdge(func(e Edge) (terminate bool) {
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
}

func (s *GraphSuite) TestEachEdgeTermination(c *C) {
	g := s.Factory.CreateGraphFromEdges(edgeSet...)

	var hit int
	g.EachEdge(func(e Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachAdjacent(c *C) {
	g := s.Factory.CreateGraphFromEdges(edgeSet...)

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

func (s *GraphSuite) TestEachAdjacentTermination(c *C) {
	g := s.Factory.CreateGraphFromEdges(append(edgeSet, BaseEdge{"foo", "qux"})...)

	var hit int
	g.EachAdjacent("foo", func(adjacent Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachEdgeIncidentTo(c *C) {
	g := s.Factory.CreateGraphFromEdges(edgeSet...)
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
	g := s.Factory.CreateGraphFromEdges(edgeSet...)

	var hit int
	g.EachEdgeIncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestDegreeOf(c *C) {
	g := s.Factory.CreateGraphFromEdges(append(edgeSet, BaseEdge{"foo", "qux"})...)

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
