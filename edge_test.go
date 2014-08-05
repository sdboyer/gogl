package gogl_test

import (
	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/spec"
	"gopkg.in/fatih/set.v0"
)

type EdgeListSuite struct{}

var _ = Suite(&EdgeListSuite{})

func (s *EdgeListSuite) TestEachVertex(c *C) {
	set1 := set.NewNonTS()

	spec.GraphFixtures["2e3v"].EachVertex(func(v Vertex) (terminate bool) {
		set1.Add(v)
		return
	})

	c.Assert(set1.Size(), Equals, 3)
	c.Assert(set1.Has("foo"), Equals, true)
	c.Assert(set1.Has("bar"), Equals, true)
	c.Assert(set1.Has("baz"), Equals, true)
}

func (s *EdgeListSuite) TestEachVertexTermination(c *C) {
	var hit int
	spec.GraphFixtures["2e3v"].EachVertex(func(v Vertex) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 1)

	spec.GraphFixtures["w-2e3v"].EachVertex(func(v Vertex) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 2)

	spec.GraphFixtures["l-2e3v"].EachVertex(func(v Vertex) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 3)

	spec.GraphFixtures["d-2e3v"].EachVertex(func(v Vertex) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 4)
}

func (s *EdgeListSuite) TestEachEdgeTermination(c *C) {
	var hit int
	spec.GraphFixtures["2e3v"].EachEdge(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 1)

	spec.GraphFixtures["w-2e3v"].EachEdge(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 2)

	spec.GraphFixtures["l-2e3v"].EachEdge(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 3)

	spec.GraphFixtures["d-2e3v"].EachEdge(func(e Edge) (terminate bool) {
		hit++
		return true
	})
	c.Assert(hit, Equals, 4)
}
