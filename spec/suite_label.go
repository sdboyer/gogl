package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

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
