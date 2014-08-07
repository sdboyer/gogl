package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

/* Suites for mutable graph methods */

type VertexSetMutatorSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
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

type EdgeSetMutatorSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
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

func (s *ArcSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

type ArcSetMutatorSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
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

type MutableGraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *MutableGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
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

