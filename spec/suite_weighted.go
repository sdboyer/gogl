package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

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
