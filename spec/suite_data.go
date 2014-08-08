package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

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

