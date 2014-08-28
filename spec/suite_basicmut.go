package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

/* Suites for mutable graph methods */

type VertexSetMutatorSuite struct {
	Factory func(GraphSource) Graph
}

func (s *VertexSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *VertexSetMutatorSuite) TestEnsureVertex(c *C) {
	g := s.Factory(NullGraph)
	m := g.(VertexSetMutator)

	m.EnsureVertex("foo")
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *VertexSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph)
	m := g.(VertexSetMutator)

	m.EnsureVertex()
	c.Assert(Order(g), Equals, 0)

	m.RemoveVertex()
	c.Assert(Order(g), Equals, 0)
}

func (s *VertexSetMutatorSuite) TestMultiEnsureVertex(c *C) {
	g := s.Factory(NullGraph)
	m := g.(VertexSetMutator)

	m.EnsureVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, true)
	c.Assert(g.HasVertex("baz"), Equals, true)
}

func (s *VertexSetMutatorSuite) TestRemoveVertex(c *C) {
	g := s.Factory(NullGraph)
	m := g.(VertexSetMutator)

	m.EnsureVertex("bar", "baz")
	m.RemoveVertex("bar")
	c.Assert(g.HasVertex("bar"), Equals, false)
}

func (s *VertexSetMutatorSuite) TestMultiRemoveVertex(c *C) {
	g := s.Factory(NullGraph)
	m := g.(VertexSetMutator)

	m.EnsureVertex("bar", "baz")
	m.RemoveVertex("bar", "baz")
	c.Assert(g.HasVertex("bar"), Equals, false)
	c.Assert(g.HasVertex("baz"), Equals, false)
}

type EdgeSetMutatorSuite struct {
	Factory func(GraphSource) Graph
}

func (s *EdgeSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *EdgeSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph)
	m := g.(EdgeSetMutator)

	m.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *EdgeSetMutatorSuite) TestAddRemoveHasEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(EdgeSetMutator)

	m.AddEdges(NewEdge(1, 2))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, true)

	// Now test removal
	m.RemoveEdges(NewEdge(1, 2))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, false)
}

func (s *EdgeSetMutatorSuite) TestMultiAddRemoveHasEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(EdgeSetMutator)

	m.AddEdges(NewEdge(1, 2), NewEdge(2, 3))

	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(2, 1)), Equals, true)
	c.Assert(g.HasEdge(NewEdge(3, 2)), Equals, true)

	// Now test removal
	m.RemoveEdges(NewEdge(1, 2), NewEdge(2, 3))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
	c.Assert(g.HasEdge(NewEdge(2, 3)), Equals, false)
}

// Checks to ensure that removal works for both in-edges and out-edges.
func (s *EdgeSetMutatorSuite) TestVertexRemovalAlsoRemovesConnectedEdges(c *C) {
	g := s.Factory(NullGraph)
	m := g.(EdgeSetMutator)

	if v, ok := g.(VertexSetMutator); ok { // Almost always gonna be the case that we have both
		m.AddEdges(NewEdge(1, 2), NewEdge(2, 3), NewEdge(4, 1))
		v.RemoveVertex(1)

		c.Assert(Size(g), Equals, 1)
	}
}

func (s *ArcSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

type ArcSetMutatorSuite struct {
	Factory func(GraphSource) Graph
}

func (s *ArcSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph)
	m := g.(ArcSetMutator)

	m.AddArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *ArcSetMutatorSuite) TestAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(Digraph)
	m := g.(ArcSetMutator)

	m.AddArcs(NewArc(1, 2))
	c.Assert(g.HasArc(NewArc(1, 2)), Equals, true)
	c.Assert(g.HasArc(NewArc(2, 1)), Equals, false)

	// Now test removal
	m.RemoveArcs(NewArc(1, 2))
	c.Assert(g.HasArc(NewArc(1, 2)), Equals, false)
	c.Assert(g.HasArc(NewArc(2, 1)), Equals, false)
}

func (s *ArcSetMutatorSuite) TestMultiAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(Digraph)
	m := g.(ArcSetMutator)

	m.AddArcs(NewArc(1, 2), NewArc(2, 3))

	c.Assert(g.HasArc(NewArc(1, 2)), Equals, true)
	c.Assert(g.HasArc(NewArc(2, 3)), Equals, true)
	c.Assert(g.HasArc(NewArc(2, 1)), Equals, false)
	c.Assert(g.HasArc(NewArc(3, 2)), Equals, false)

	// Now test removal
	m.RemoveArcs(NewArc(1, 2), NewArc(2, 3))
	c.Assert(g.HasArc(NewArc(1, 2)), Equals, false)
	c.Assert(g.HasArc(NewArc(1, 2)), Equals, false)
	c.Assert(g.HasArc(NewArc(2, 3)), Equals, false)
	c.Assert(g.HasArc(NewArc(2, 3)), Equals, false)
}

// Checks to ensure that removal works for both in-edges and out-edges.
func (s *ArcSetMutatorSuite) TestVertexRemovalAlsoRemovesConnectedArcs(c *C) {
	g := s.Factory(NullGraph)
	m := g.(ArcSetMutator)

	if v, ok := g.(VertexSetMutator); ok {
		// Almost always gonona be the case that we have both
		m.AddArcs(NewArc(1, 2), NewArc(2, 3), NewArc(4, 1))
		v.RemoveVertex(1)

		c.Assert(Size(g), Equals, 1)
	}
}
