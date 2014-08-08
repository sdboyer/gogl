package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

/* WeightedGraphSuite - tests for weighted graphs */

type WeightedGraphSuite struct {
	Factory  func(GraphSource) WeightedGraph
}

func (s *WeightedGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *WeightedGraphSuite) TestEachEdge(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachEdge() iterator actually do implement WeightedEdge.
	g := s.Factory(GraphFixtures["w-2e3v"])

	var we WeightedEdge
	g.EachEdge(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *WeightedGraphSuite) TestHasWeightedEdge(c *C) {
	g := s.Factory(GraphFixtures["w-2e3v"])

	c.Assert(g.HasWeightedEdge(GraphFixtures["w-2e3v"].(WeightedArcList)[0].(WeightedArc)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 1)), Equals, false) // wrong weight
}

type WeightedDigraphSuite struct {
	Factory  func(GraphSource) WeightedGraph
}

func (s *WeightedDigraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *WeightedDigraphSuite) TestArcSubtypeImplementation(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachArc() iterator actually do implement WeightedArc.
	g := s.Factory(GraphFixtures["w-2e3v"]).(WeightedDigraph)

	var hit int // just internal safety check to ensure the fixture is good and hits
	var wa WeightedArc
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

/* WeightedEdgeSetMutatorSuite - tests for mutable weighted graphs */

type WeightedEdgeSetMutatorSuite struct {
	Factory  func(GraphSource) WeightedGraph
}

func (s *WeightedEdgeSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *WeightedEdgeSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph)
	m := g.(WeightedEdgeSetMutator)

	m.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *WeightedEdgeSetMutatorSuite) TestAddRemoveHasEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(WeightedEdgeSetMutator)
	m.AddEdges(NewWeightedEdge(1, 2, 5.23))

	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, true)

	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5.23)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 3)), Equals, false)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, 5.23)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, -3.22771)), Equals, false)

	// Now test removal
	m.RemoveEdges(NewWeightedEdge(1, 2, 5.23))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5.23)), Equals, false)
}

func (s *WeightedEdgeSetMutatorSuite) TestMultiAddRemoveHasEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(WeightedEdgeSetMutator)
	m.AddEdges(NewWeightedEdge(1, 2, 5), NewWeightedEdge(2, 3, -5))

	// Basic edge tests first
	// We test both Has*Edge() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Edge() method.
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, true) // only if undirected
	c.Assert(g.HasEdge(NewEdge(3, 2)), Equals, true) // only if undirected

	// Now weighted edge tests
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 3)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, 5)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 1, 3)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 3, -5)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 3, 1)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(3, 2, -5)), Equals, true)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(3, 2, 1)), Equals, false) // wrong weight

	// Now test removal
	m.RemoveEdges(NewWeightedEdge(1, 2, 5), NewWeightedEdge(2, 3, -5))
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(1, 2, 5)), Equals, false)
	c.Assert(g.HasWeightedEdge(NewWeightedEdge(2, 3, -5)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
}

/* WeightedArcSetMutatorSuite - tests for mutable weighted graphs */

type WeightedArcSetMutatorSuite struct {
	Factory  func(GraphSource) WeightedGraph
}

func (s *WeightedArcSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *WeightedArcSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph).(WeightedDigraph)
	m := g.(WeightedArcSetMutator)

	m.AddArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *WeightedArcSetMutatorSuite) TestAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(WeightedDigraph)
	m := g.(WeightedArcSetMutator)
	m.AddArcs(NewWeightedArc(1, 2, 5.23))

	c.Assert(g.HasArc(NewArc(1, 2)), Equals, true)
	c.Assert(g.HasArc(NewArc(2, 1)), Equals, false) // wrong direction

	c.Assert(g.HasWeightedArc(NewWeightedArc(1, 2, 5.23)), Equals, true)
	c.Assert(g.HasWeightedArc(NewWeightedArc(1, 2, 3)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedArc(NewWeightedArc(2, 1, 5.23)), Equals, false) // wrong direction
	c.Assert(g.HasWeightedArc(NewWeightedArc(2, 1, 3)), Equals, false) // wrong direction & weight

	// Now test removal
	m.RemoveArcs(NewWeightedArc(1, 2, 5.23))
	c.Assert(g.HasArc(NewArc(1, 2)), Equals, false)
	c.Assert(g.HasWeightedArc(NewWeightedArc(1, 2, 5.23)), Equals, false)
}

func (s *WeightedArcSetMutatorSuite) TestMultiAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(WeightedDigraph)
	m := g.(WeightedArcSetMutator)
	m.AddArcs(NewWeightedArc(1, 2, 5), NewWeightedArc(2, 3, -5))

	// Basic edge tests first
	// We test both Has*Arc() methods to ensure that adding our known edge fixture type results in the expected behavior.
	// Thus, this is not just duplicate testing of the Has*Arc() method.
	c.Assert(g.HasArc(NewArc(1, 2)), Equals, true)
	c.Assert(g.HasArc(NewArc(2, 3)), Equals, true)
	c.Assert(g.HasArc(NewArc(2, 1)), Equals, false)
	c.Assert(g.HasArc(NewArc(3, 2)), Equals, false)

	// Now weighted edge tests
	c.Assert(g.HasWeightedArc(NewWeightedArc(1, 2, 5)), Equals, true)
	c.Assert(g.HasWeightedArc(NewWeightedArc(1, 2, 3)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedArc(NewWeightedArc(2, 1, 5)), Equals, false) // wrong direction
	c.Assert(g.HasWeightedArc(NewWeightedArc(2, 1, 3)), Equals, false) // wrong direction & weight
	c.Assert(g.HasWeightedArc(NewWeightedArc(2, 3, -5)), Equals, true)
	c.Assert(g.HasWeightedArc(NewWeightedArc(2, 3, 1)), Equals, false) // wrong weight
	c.Assert(g.HasWeightedArc(NewWeightedArc(3, 2, -5)), Equals, false) // wrong direction
	c.Assert(g.HasWeightedArc(NewWeightedArc(3, 2, 1)), Equals, false) // wrong direction & weight

	// Now test removal
	m.RemoveArcs(NewWeightedArc(1, 2, 5), NewWeightedArc(2, 3, -5))
	c.Assert(g.HasWeightedArc(NewWeightedArc(1, 2, 5)), Equals, false)
	c.Assert(g.HasWeightedArc(NewWeightedArc(2, 3, -5)), Equals, false)
	c.Assert(g.HasArc(NewArc(1, 2)), Equals, false)
	c.Assert(g.HasArc(NewArc(2, 3)), Equals, false)
}
