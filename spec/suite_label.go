package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

/* LabeledGraphSuite - tests for labeled graphs */

type LabeledGraphSuite struct {
	Factory  func(GraphSource) LabeledGraph
}

func (s *LabeledGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *LabeledGraphSuite) TestEachEdge(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachEdge() iterator actually do implement LabeledEdge.
	g := s.Factory(GraphFixtures["l-2e3v"])

	var we LabeledEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *LabeledGraphSuite) TestHasLabeledEdge(c *C) {
	g := s.Factory(GraphFixtures["l-2e3v"])

	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 1, "foo")), Equals, true) // both directions work
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "qux")), Equals, false) // wrong label
}

type LabeledDigraphSuite struct {
	Factory  func(GraphSource) LabeledGraph
}

func (s *LabeledDigraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *LabeledDigraphSuite) TestArcSubtypeImplementation(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachArc() iterator actually do implement LabeledArc.
	g := s.Factory(GraphFixtures["l-2e3v"]).(LabeledDigraph)

	var hit int // just internal safety check to ensure the fixture is good and hits
	var wa LabeledArc
	g.EachArc(func(e Arc) (terminate bool) {
		hit++
		c.Assert(e, Implements, &wa)
		return
	})

	g.EachArcFrom(2, func(e Arc) (terminate bool) {
		hit++
		c.Assert(e, Implements, &wa)
		return
	})

	g.EachArcFrom(2, func(e Arc) (terminate bool) {
		hit++
		c.Assert(e, Implements, &wa)
		return
	})

	c.Assert(hit, Equals, 4)
}

/* LabeledEdgeSetMutatorSuite - tests for mutable labeled graphs */

type LabeledEdgeSetMutatorSuite struct {
	Factory  func(GraphSource) LabeledGraph
}

func (s *LabeledEdgeSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *LabeledEdgeSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph)
	m := g.(LabeledEdgeSetMutator)

	m.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *LabeledEdgeSetMutatorSuite) TestAddRemoveEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(LabeledEdgeSetMutator)

	m.AddEdges(NewLabeledEdge(1, 2, "foo"))
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, true)

	// Now test removal
	m.RemoveEdges(NewLabeledEdge(1, 2, "foo"))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, false)
}

func (s *LabeledEdgeSetMutatorSuite) TestMultiAddRemoveEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(LabeledEdgeSetMutator)

	m.AddEdges(NewLabeledEdge(1, 2, "foo"), NewLabeledEdge(2, 3, "bar"))
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 3, "bar")), Equals, true)

	// Now test removal
	m.RemoveEdges(NewLabeledEdge(1, 2, "foo"), NewLabeledEdge(2, 3, "bar"))
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(1, 2, "foo")), Equals, false)
	c.Assert(g.HasLabeledEdge(NewLabeledEdge(2, 3, "bar")), Equals, false)
}

/* LabeledArcSetMutatorSuite - tests for mutable labeled graphs */

type LabeledArcSetMutatorSuite struct {
	Factory  func(GraphSource) LabeledGraph
}

func (s *LabeledArcSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *LabeledArcSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph).(LabeledDigraph)
	m := g.(LabeledArcSetMutator)

	m.AddArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *LabeledArcSetMutatorSuite) TestAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(LabeledDigraph)
	m := g.(LabeledArcSetMutator)

	m.AddArcs(NewLabeledArc(1, 2, "foo"))
	c.Assert(g.HasLabeledArc(NewLabeledArc(1, 2, "foo")), Equals, true)
	c.Assert(g.HasLabeledArc(NewLabeledArc(1, 2, "qux")), Equals, false) // wrong label

	// Now test removal
	m.RemoveArcs(NewLabeledArc(1, 2, "foo"))
	c.Assert(g.HasLabeledArc(NewLabeledArc(1, 2, "foo")), Equals, false)
}

func (s *LabeledArcSetMutatorSuite) TestMultiAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(LabeledDigraph)
	m := g.(LabeledArcSetMutator)

	m.AddArcs(NewLabeledArc(1, 2, "foo"), NewLabeledArc(2, 3, "bar"))
	c.Assert(g.HasLabeledArc(NewLabeledArc(1, 2, "foo")), Equals, true)
	c.Assert(g.HasLabeledArc(NewLabeledArc(2, 3, "bar")), Equals, true)

	m.RemoveArcs(NewLabeledArc(1, 2, "foo"), NewLabeledArc(2, 3, "bar"))
	c.Assert(g.HasLabeledArc(NewLabeledArc(1, 2, "foo")), Equals, false)
	c.Assert(g.HasLabeledArc(NewLabeledArc(2, 3, "bar")), Equals, false)
}
