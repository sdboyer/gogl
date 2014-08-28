package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

/* GraphSuite - tests for non-mutable graph methods */

type GraphSuite struct {
	Factory  func(GraphSource) Graph
	Directed bool
}

func (s *GraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *GraphSuite) TestHasVertex(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])
	c.Assert(g.HasVertex("qux"), Equals, false)
	c.Assert(g.HasVertex("foo"), Equals, true)
}

func (s *GraphSuite) TestHasEdge(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])
	// Testing match with minimum possible specificity here
	c.Assert(g.HasEdge(GraphFixtures["2e3v"].(ArcList)[0]), Equals, true)
	c.Assert(g.HasEdge(NewEdge("qux", "quark")), Equals, false)
}

func (s *GraphSuite) TestVertices(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	vset := set.NewNonTS()
	var hit int
	g.Vertices(func(v Vertex) (terminate bool) {
		hit++
		vset.Add(v)
		return
	})

	c.Assert(vset.Has("foo"), Equals, true)
	c.Assert(vset.Has("bar"), Equals, true)
	c.Assert(vset.Has("baz"), Equals, true)
	c.Assert(hit, Equals, 3)
}

func (s *GraphSuite) TestVerticesTermination(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.Vertices(func(v Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachEdge(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.EachEdge(func(e Edge) (terminate bool) {
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
}

func (s *GraphSuite) TestEachEdgeTermination(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.EachEdge(func(e Edge) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestEachAdjacentTo(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	vset := set.NewNonTS()
	var hit int
	g.EachAdjacentTo("bar", func(adj Vertex) (terminate bool) {
		hit++
		vset.Add(adj)
		return
	})

	c.Assert(vset.Has("foo"), Equals, true)
	c.Assert(vset.Has("bar"), Equals, false)
	c.Assert(vset.Has("baz"), Equals, true)
	c.Assert(hit, Equals, 2)
}

func (s *GraphSuite) TestEachAdjacentToTermination(c *C) {
	g := s.Factory(GraphFixtures["3e4v"])

	var hit int
	g.EachAdjacentTo("foo", func(adjacent Vertex) bool {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestIncidentTo(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	flipset := []Edge{
		Swap(GraphFixtures["2e3v"].(ArcList)[0]),
		Swap(GraphFixtures["2e3v"].(ArcList)[1]),
	}

	eset := set.NewNonTS()
	var hit int
	g.IncidentTo("foo", func(e Edge) (terminate bool) {
		hit++
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewArc(e.Both()))
		return
	})

	c.Assert(hit, Equals, 1)
	if s.Directed {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[1]), Equals, false)
	} else {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[0]) != eset.Has(flipset[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[1]) != eset.Has(flipset[1]), Equals, false)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[1]), Equals, false)
	}

	eset = set.NewNonTS()
	g.IncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewArc(e.Both()))
		return
	})

	c.Assert(hit, Equals, 3)
	if s.Directed {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[1]), Equals, true)
	} else {
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[0]) != eset.Has(flipset[0]), Equals, true)
		c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[1]) != eset.Has(flipset[1]), Equals, true)
	}
}

func (s *GraphSuite) TestIncidentToTermination(c *C) {
	g := s.Factory(GraphFixtures["2e3v"])

	var hit int
	g.IncidentTo("bar", func(e Edge) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *GraphSuite) TestDegreeOf(c *C) {
	g := s.Factory(GraphFixtures["3e5v1i"])

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

	count, exists = g.DegreeOf("isolate")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.DegreeOf("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}
