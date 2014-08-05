package gogl_test

import (
	"testing"

	"github.com/sdboyer/gogl/spec"
	. "github.com/sdboyer/gogl"
	. "github.com/sdboyer/gocheck"
	"gopkg.in/fatih/set.v0"

)

// Hook gocheck into the go test runner
func TestHookup(t *testing.T) { TestingT(t) }

// Tests for collection functors
type CollectionFunctorsSuite struct{}

var _ = Suite(&CollectionFunctorsSuite{})

func (s *CollectionFunctorsSuite) TestCollectVertices(c *C) {
	slice := CollectVertices(spec.GraphLiteralFixture(true))

	c.Assert(len(slice), Equals, 4)

	set := set.NewNonTS()
	for _, v := range slice {
		set.Add(v)
	}

	c.Assert(set.Has("foo"), Equals, true)
	c.Assert(set.Has("bar"), Equals, true)
	c.Assert(set.Has("baz"), Equals, true)
	c.Assert(set.Has("isolate"), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectAdjacentVertices(c *C) {
	slice := CollectVerticesAdjacentTo("bar", spec.GraphLiteralFixture(true))

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, v := range slice {
		set.Add(v)
	}

	c.Assert(set.Has("foo"), Equals, true)
	c.Assert(set.Has("baz"), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectEdges(c *C) {
	slice := CollectEdges(spec.GraphLiteralFixture(true))

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(NewEdge("foo", "bar")), Equals, true)
	c.Assert(set.Has(NewEdge("bar", "baz")), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectEdgesIncidentTo(c *C) {
	slice := CollectEdgesIncidentTo("foo", spec.GraphLiteralFixture(true))

	c.Assert(len(slice), Equals, 1)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(NewEdge("foo", "bar")), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectArcsFrom(c *C) {
	slice := CollectArcsFrom("foo", spec.GraphLiteralFixture(true))

	c.Assert(len(slice), Equals, 1)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(NewEdge("foo", "bar")), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectArcsTo(c *C) {
	slice := CollectArcsTo("bar", spec.GraphLiteralFixture(true))

	c.Assert(len(slice), Equals, 1)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(NewEdge("foo", "bar")), Equals, true)
}

type CountingFunctorsSuite struct{}

var _ = Suite(&CountingFunctorsSuite{})

func (s *CountingFunctorsSuite) TestOrder(c *C) {
	el := EdgeList{
		NewEdge("foo", "bar"),
		NewEdge("bar", "baz"),
		NewEdge("foo", "qux"),
		NewEdge("qux", "bar"),
	}
	c.Assert(Order(el), Equals, 4)
	c.Assert(Order(spec.GraphLiteralFixture(true)), Equals, 4)
}

func (s *CountingFunctorsSuite) TestSize(c *C) {
	el := EdgeList{
		NewEdge("foo", "bar"),
		NewEdge("bar", "baz"),
		NewEdge("foo", "qux"),
		NewEdge("qux", "bar"),
	}
	c.Assert(Size(el), Equals, 4)
	c.Assert(Size(spec.GraphLiteralFixture(true)), Equals, 2)
}
