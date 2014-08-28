package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
)

/* DataGraphSuite - tests for data graphs */

type DataGraphSuite struct {
	Factory func(GraphSource) DataGraph
}

func (s *DataGraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DataGraphSuite) TestEdges(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the Edges() iterator actually do implement DataEdge.
	g := s.Factory(GraphFixtures["d-2e3v"])

	var we DataEdge
	g.Edges(func(e Edge) (terminate bool) {
		c.Assert(e, Implements, &we)
		return
	})
}

func (s *DataGraphSuite) TestHasDataEdge(c *C) {
	g := s.Factory(GraphFixtures["d-2e3v"])

	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 1, "foo")), Equals, true)  // both directions work
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "qux")), Equals, false) // wrong data
}

type DataDigraphSuite struct {
	Factory func(GraphSource) DataGraph
}

func (s *DataDigraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DataDigraphSuite) TestArcSubtypeImplementation(c *C) {
	// This method is not redundant with the base Graph suite as it ensures that the edges
	// provided by the EachArc() iterator actually do implement DataArc.
	g := s.Factory(GraphFixtures["d-2e3v"]).(DataDigraph)

	var hit int // just internal safety check to ensure the fixture is good and hits
	var wa DataArc
	g.EachArc(func(e Arc) (terminate bool) {
		hit++
		c.Assert(e, Implements, &wa)
		return
	})

	g.ArcsFrom(2, func(e Arc) (terminate bool) {
		hit++
		c.Assert(e, Implements, &wa)
		return
	})

	g.ArcsFrom(2, func(e Arc) (terminate bool) {
		hit++
		c.Assert(e, Implements, &wa)
		return
	})

	c.Assert(hit, Equals, 4)
}

/* DataEdgeSetMutatorSuite - tests for mutable data graphs */

type DataEdgeSetMutatorSuite struct {
	Factory func(GraphSource) DataGraph
}

func (s *DataEdgeSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DataEdgeSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph)
	m := g.(DataEdgeSetMutator)

	m.AddEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveEdges()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *DataEdgeSetMutatorSuite) TestAddRemoveEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(DataEdgeSetMutator)

	m.AddEdges(NewDataEdge(1, 2, "foo"))
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, true)

	// Now test removal
	m.RemoveEdges(NewDataEdge(1, 2, "foo"))
	c.Assert(g.HasEdge(NewEdge(1, 2)), Equals, false)
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, false)
}

func (s *DataEdgeSetMutatorSuite) TestMultiAddRemoveEdge(c *C) {
	g := s.Factory(NullGraph)
	m := g.(DataEdgeSetMutator)

	m.AddEdges(NewDataEdge(1, 2, "foo"), NewDataEdge(2, 3, "bar"))
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, true)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 3, "bar")), Equals, true)

	// Now test removal
	m.RemoveEdges(NewDataEdge(1, 2, "foo"), NewDataEdge(2, 3, "bar"))
	c.Assert(g.HasDataEdge(NewDataEdge(1, 2, "foo")), Equals, false)
	c.Assert(g.HasDataEdge(NewDataEdge(2, 3, "bar")), Equals, false)
}

/* DataArcSetMutatorSuite - tests for mutable data graphs */

type DataArcSetMutatorSuite struct {
	Factory func(GraphSource) DataGraph
}

func (s *DataArcSetMutatorSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DataArcSetMutatorSuite) TestGracefulEmptyVariadics(c *C) {
	g := s.Factory(NullGraph).(DataDigraph)
	m := g.(DataArcSetMutator)

	m.AddArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)

	m.RemoveArcs()
	c.Assert(Order(g), Equals, 0)
	c.Assert(Size(g), Equals, 0)
}

func (s *DataArcSetMutatorSuite) TestAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(DataDigraph)
	m := g.(DataArcSetMutator)

	m.AddArcs(NewDataArc(1, 2, "foo"))
	c.Assert(g.HasDataArc(NewDataArc(1, 2, "foo")), Equals, true)
	c.Assert(g.HasDataArc(NewDataArc(1, 2, "bar")), Equals, false) // wrong data

	// Now test removal
	m.RemoveArcs(NewDataArc(1, 2, "foo"))
	c.Assert(g.HasDataArc(NewDataArc(1, 2, "foo")), Equals, false)
}

func (s *DataArcSetMutatorSuite) TestMultiAddRemoveHasArc(c *C) {
	g := s.Factory(NullGraph).(DataDigraph)
	m := g.(DataArcSetMutator)

	m.AddArcs(NewDataArc(1, 2, "foo"), NewDataArc(2, 3, "bar"))
	c.Assert(g.HasDataArc(NewDataArc(1, 2, "foo")), Equals, true)
	c.Assert(g.HasDataArc(NewDataArc(2, 3, "bar")), Equals, true)

	m.RemoveArcs(NewDataArc(1, 2, "foo"), NewDataArc(2, 3, "bar"))
	c.Assert(g.HasDataArc(NewDataArc(1, 2, "foo")), Equals, false)
	c.Assert(g.HasDataArc(NewDataArc(2, 3, "bar")), Equals, false)
}
