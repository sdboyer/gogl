package gogl

import (
	. "github.com/sdboyer/gocheck"
	"gopkg.in/fatih/set.v0"
)

// Tests for collection functors
type CollectionFunctorsSuite struct{}

var _ = Suite(&CollectionFunctorsSuite{})

func (s *CollectionFunctorsSuite) TestCollectVertices(c *C) {
	slice := CollectVertices(graphFixtures["3e5v1i"].(Graph))

	c.Assert(len(slice), Equals, 5)

	set := set.NewNonTS()
	for _, v := range slice {
		set.Add(v)
	}

	c.Assert(set.Has("foo"), Equals, true)
	c.Assert(set.Has("bar"), Equals, true)
	c.Assert(set.Has("baz"), Equals, true)
	c.Assert(set.Has("qux"), Equals, true)
	c.Assert(set.Has("isolate"), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectAdjacentVertices(c *C) {
	slice := CollectVerticesAdjacentTo("foo", graphFixtures["3e5v1i"].(Graph))

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, v := range slice {
		set.Add(v)
	}

	c.Assert(set.Has("bar"), Equals, true)
	c.Assert(set.Has("qux"), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectEdges(c *C) {
	slice := CollectEdges(graphFixtures["3e5v1i"].(Graph))

	c.Assert(len(slice), Equals, 3)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(BaseEdge{"foo", "bar"}), Equals, true)
	c.Assert(set.Has(BaseEdge{"bar", "baz"}), Equals, true)
	c.Assert(set.Has(BaseEdge{"foo", "qux"}), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectEdgesIncidentTo(c *C) {
	slice := CollectEdgesIncidentTo("foo", graphFixtures["3e5v1i"].(Graph))

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(BaseEdge{"foo", "bar"}), Equals, true)
	c.Assert(set.Has(BaseEdge{"foo", "qux"}), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectArcsFrom(c *C) {
	digraph := BuildGraph().Directed().Using(graphFixtures["arctest"]).Create(AdjacencyList).(DirectedGraph)
	slice := CollectArcsFrom("foo", digraph)

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(BaseEdge{"foo", "qux"}), Equals, true)
	c.Assert(set.Has(BaseEdge{"foo", "bar"}), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectArcsTo(c *C) {
	digraph := BuildGraph().Directed().Using(graphFixtures["arctest"]).Create(AdjacencyList).(DirectedGraph)
	slice := CollectArcsTo("bar", digraph)

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(BaseEdge{"foo", "bar"}), Equals, true)
	c.Assert(set.Has(BaseEdge{"qux", "bar"}), Equals, true)
}
