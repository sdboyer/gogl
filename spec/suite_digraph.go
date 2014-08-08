package spec

import (
	"fmt"

	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

/* DigraphSuite - tests for directed graph methods */

type DigraphSuite struct {
	Factory func(GraphSource) Graph
}

func (s *DigraphSuite) SuiteLabel() string {
	return fmt.Sprintf("%T", s.Factory(NullGraph))
}

func (s *DigraphSuite) TestTranspose(c *C) {
	g := s.Factory(GraphFixtures["2e3v"]).(Digraph)

	g2 := g.Transpose()

	c.Assert(g2.HasArc(Swap(GraphFixtures["2e3v"].(ArcList)[0])), Equals, true)
	c.Assert(g2.HasArc(Swap(GraphFixtures["2e3v"].(ArcList)[1])), Equals, true)

	c.Assert(g2.HasArc(GraphFixtures["2e3v"].(ArcList)[0]), Equals, false)
	c.Assert(g2.HasArc(GraphFixtures["2e3v"].(ArcList)[1]), Equals, false)
}

func (s *DigraphSuite) TestOutDegreeOf(c *C) {
	g := s.Factory(GraphFixtures["3e5v1i"]).(Digraph)

	count, exists := g.OutDegreeOf("foo")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 2)

	count, exists = g.OutDegreeOf("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.OutDegreeOf("baz")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.OutDegreeOf("qux")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.DegreeOf("isolate")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.OutDegreeOf("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *DigraphSuite) TestInDegreeOf(c *C) {
	g := s.Factory(GraphFixtures["3e5v1i"]).(Digraph)

	count, exists := g.InDegreeOf("foo")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.InDegreeOf("bar")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.InDegreeOf("baz")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.InDegreeOf("qux")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 1)

	count, exists = g.DegreeOf("isolate")
	c.Assert(exists, Equals, true)
	c.Assert(count, Equals, 0)

	count, exists = g.InDegreeOf("missing")
	c.Assert(exists, Equals, false)
	c.Assert(count, Equals, 0)
}

func (s *DigraphSuite) TestEachArcTo(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	var hit int
	g.EachArcTo("foo", func(e Arc) (terminate bool) {
		c.Error("Vertex 'foo' should have no in-edges")
		c.FailNow()
		return
	})

	g.EachArcTo("bar", func(e Arc) (terminate bool) {
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewArc(e.Both()))
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[0]), Equals, true)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[1]), Equals, false)
	c.Assert(eset.Has(NewArc("qux", "bar")), Equals, true)
}

func (s *DigraphSuite) TestEachArcToTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachArcTo("baz", func(e Arc) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DigraphSuite) TestEachPredecessorOf(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	g.EachPredecessorOf("foo", func(v Vertex) (terminate bool) {
		c.Error("Vertex 'foo' should have no predecessors")
		c.FailNow()
		return
	})

	g.EachPredecessorOf("bar", func(v Vertex) (terminate bool) {
		eset.Add(v)
		return
	})

	c.Assert(eset.Size(), Equals, 2)
	c.Assert(eset.Has("foo"), Equals, true)
	c.Assert(eset.Has("qux"), Equals, true)

}

func (s *DigraphSuite) TestEachPredecessorOfTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachPredecessorOf("baz", func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DigraphSuite) TestEachArcFrom(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	var hit int
	g.EachArcFrom("baz", func(e Arc) (terminate bool) {
		c.Error("Vertex 'baz' should have no out-edges")
		c.FailNow()
		return
	})

	g.EachArcFrom("foo", func(e Arc) (terminate bool) {
		// A more specific edge type may be passed, but in this test we care only about the base
		eset.Add(NewArc(e.Both()))
		hit++
		return
	})

	c.Assert(hit, Equals, 2)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[0]), Equals, true)
	c.Assert(eset.Has(GraphFixtures["2e3v"].(ArcList)[1]), Equals, false)
	c.Assert(eset.Has(NewArc("foo", "qux")), Equals, true)
}

func (s *DigraphSuite) TestEachArcFromTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachArcFrom("foo", func(e Arc) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

func (s *DigraphSuite) TestEachSuccessorOf(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	eset := set.NewNonTS()
	g.EachSuccessorOf("baz", func(v Vertex) (terminate bool) {
		c.Error("Vertex 'foo' should have no successors")
		c.FailNow()
		return
	})

	g.EachSuccessorOf("foo", func(v Vertex) (terminate bool) {
		eset.Add(v)
		return
	})

	c.Assert(eset.Size(), Equals, 2)
	c.Assert(eset.Has("qux"), Equals, true)
	c.Assert(eset.Has("bar"), Equals, true)

}

func (s *DigraphSuite) TestEachSuccessorOfTermination(c *C) {
	g := s.Factory(GraphFixtures["arctest"]).(Digraph)

	var hit int
	g.EachSuccessorOf("foo", func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)
}

