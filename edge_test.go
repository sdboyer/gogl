package gogl_test

import (
	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/spec"
	"gopkg.in/fatih/set.v0"
)

type EdgeListSuite struct{}

var _ = Suite(&EdgeListSuite{})

func (s *EdgeListSuite) TestVertices(c *C) {
	set1 := set.NewNonTS()

	spec.GraphFixtures["2e3v"].Vertices(func(v Vertex) (terminate bool) {
		set1.Add(v)
		return
	})

	c.Assert(set1.Size(), Equals, 3)
	c.Assert(set1.Has("foo"), Equals, true)
	c.Assert(set1.Has("bar"), Equals, true)
	c.Assert(set1.Has("baz"), Equals, true)
}

func (s *EdgeListSuite) TestVerticesTermination(c *C) {
	var hit int
	spec.GraphFixtures["2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)

	spec.GraphFixtures["w-2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 2)

	spec.GraphFixtures["l-2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 3)

	spec.GraphFixtures["d-2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 4)
}

func (s *EdgeListSuite) TestEdgesTermination(c *C) {
	var hit int
	spec.GraphFixtures["2e3v"].Edges(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 1)

	spec.GraphFixtures["w-2e3v"].Edges(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 2)

	spec.GraphFixtures["l-2e3v"].Edges(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 3)

	spec.GraphFixtures["d-2e3v"].Edges(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 4)
}

type ArcListSuite struct{}

var _ = Suite(&ArcListSuite{})

func (s *ArcListSuite) TestVertices(c *C) {
	set1 := set.NewNonTS()

	spec.GraphFixtures["2e3v"].Vertices(func(v Vertex) (terminate bool) {
		set1.Add(v)
		return
	})

	c.Assert(set1.Size(), Equals, 3)
	c.Assert(set1.Has("foo"), Equals, true)
	c.Assert(set1.Has("bar"), Equals, true)
	c.Assert(set1.Has("baz"), Equals, true)
}

func (s *ArcListSuite) TestVerticesTermination(c *C) {
	var hit int
	spec.GraphFixtures["2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 1)

	spec.GraphFixtures["w-2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 2)

	spec.GraphFixtures["l-2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 3)

	spec.GraphFixtures["d-2e3v"].Vertices(func(v Vertex) (terminate bool) {
		hit++
		return true
	})

	c.Assert(hit, Equals, 4)
}

func (s *ArcListSuite) TestEachArcTermination(c *C) {
	var hit int
	spec.GraphFixtures["2e3v"].(DigraphSource).EachArc(func(e Arc) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 1)

	spec.GraphFixtures["w-2e3v"].(DigraphSource).EachArc(func(e Arc) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 2)

	spec.GraphFixtures["l-2e3v"].(DigraphSource).EachArc(func(e Arc) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 3)

	spec.GraphFixtures["d-2e3v"].(DigraphSource).EachArc(func(e Arc) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 4)
}

type EdgeSuite struct{}

var _ = Suite(&EdgeSuite{})

func (s *EdgeSuite) TestEdges(c *C) {
	e := NewEdge("a", "b")

	a, b := e.Both()
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")

	we := NewWeightedEdge("a", "b", 4.2)

	a, b = we.Both()
	c.Assert(we.Weight(), Equals, 4.2)
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")

	le := NewLabeledEdge("a", "b", "foobar")

	a, b = le.Both()
	c.Assert(le.Label(), Equals, "foobar")
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")

	de := NewDataEdge("a", "b", NullGraph)

	a, b = de.Both()
	c.Assert(de.Data(), Equals, NullGraph)
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")
}

type ArcSuite struct{}

var _ = Suite(&ArcSuite{})

func (s *ArcSuite) TestArcs(c *C) {
	e := NewArc("a", "b")

	a, b := e.Both()
	c.Assert(e.Source(), Equals, "a")
	c.Assert(e.Target(), Equals, "b")
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")

	we := NewWeightedArc("a", "b", 4.2)

	a, b = we.Both()
	c.Assert(we.Source(), Equals, "a")
	c.Assert(we.Target(), Equals, "b")
	c.Assert(we.Weight(), Equals, 4.2)
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")

	le := NewLabeledArc("a", "b", "foobar")

	a, b = le.Both()
	c.Assert(le.Source(), Equals, "a")
	c.Assert(le.Target(), Equals, "b")
	c.Assert(le.Label(), Equals, "foobar")
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")

	de := NewDataArc("a", "b", NullGraph)

	a, b = de.Both()
	c.Assert(de.Source(), Equals, "a")
	c.Assert(de.Target(), Equals, "b")
	c.Assert(de.Data(), Equals, NullGraph)
	c.Assert(a, Equals, "a")
	c.Assert(b, Equals, "b")
}
