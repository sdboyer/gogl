package gogl

import (
	"math"
	. "github.com/sdboyer/gocheck"
)

type NullGraphSuite bool

var _ = Suite(NullGraphSuite(false))

func (s NullGraphSuite) TestEnumerators(c *C) {
	NullGraph.EachVertex(func(v Vertex) (terminate bool) {
		c.Error("The NullGraph should not have any vertices.")
		return
	})

	NullGraph.EachEdge(func(e Edge) (terminate bool) {
		c.Error("The NullGraph should not have any edges.")
		return
	})

	NullGraph.EachEdgeIncidentTo("foo", func(e Edge) (terminate bool) {
		c.Error("The NullGraph should be empty of edges and vertices.")
		return
	})

	NullGraph.EachAdjacentTo("foo", func(v Vertex) (terminate bool) {
		c.Error("The NullGraph should be empty of edges and vertices.")
		return
	})

	NullGraph.EachArcFrom("foo", func(e Edge) (terminate bool) {
		c.Error("The NullGraph should be empty of edges and vertices.")
		return
	})

	NullGraph.EachArcTo("foo", func(e Edge) (terminate bool) {
		c.Error("The NullGraph should be empty of edges and vertices.")
		return
	})
}

func (s NullGraphSuite) TestMembership(c *C) {
	c.Assert(NullGraph.HasVertex("foo"), Equals, false)                                                       // must always be false
	c.Assert(NullGraph.HasEdge(BaseEdge{"qux", "quark"}), Equals, false)                                      // must always be false
	c.Assert(NullGraph.HasWeightedEdge(BaseWeightedEdge{BaseEdge{"qux", "quark"}, 0}), Equals, false)         // must always be false
	c.Assert(NullGraph.HasLabeledEdge(BaseLabeledEdge{BaseEdge{"qux", "quark"}, ""}), Equals, false)          // must always be false
	c.Assert(NullGraph.HasPropertyEdge(BasePropertyEdge{BaseEdge{"qux", "quark"}, func() {}}), Equals, false) // must always be false
}

func (s NullGraphSuite) TestDegree(c *C) {
	deg, exists := NullGraph.DegreeOf("foo")
	c.Assert(exists, Equals, false) // vertex is not present in graph
	c.Assert(deg, Equals, 0)        // always will have degree 0

	deg, exists = NullGraph.InDegreeOf("foo")
	c.Assert(exists, Equals, false) // vertex is not present in graph
	c.Assert(deg, Equals, 0)        // always will have indegree 0

	deg, exists = NullGraph.OutDegreeOf("foo")
	c.Assert(exists, Equals, false) // vertex is not present in graph
	c.Assert(deg, Equals, 0)        // always will have outdegree 0
}

func (s NullGraphSuite) TestSizingOps(c *C) {
	c.Assert(NullGraph.Size(), Equals, 0)
	c.Assert(NullGraph.Order(), Equals, 0)
	c.Assert(math.IsNaN(NullGraph.Density()), Equals, true)
}

func (s NullGraphSuite) TestTranspose(c *C) {
	c.Assert(NullGraph.Transpose(), Equals, NullGraph)
}
