package gogl

import (
	. "github.com/sdboyer/gocheck"
	"gopkg.in/fatih/set.v0"
)

// Define a graph literal fixture for testing here.
// The literal has two edges and four vertices; one vertex is an isolate.
//
// bool state indicates whether using a transpose or not.
type graphLiteralFixture bool

func (g graphLiteralFixture) EachVertex(f VertexStep) {
	vl := []Vertex{"foo", "bar", "baz", "isolate"}
	for _, v := range vl {
		if f(v) {
			return
		}
	}
}

func (g graphLiteralFixture) EachEdge(f EdgeStep) {
	var el []Edge
	if g {
		el = []Edge{
			NewEdge("foo", "bar"),
			NewEdge("bar", "baz"),
		}
	} else {
		el = []Edge{
			NewEdge("bar", "foo"),
			NewEdge("baz", "bar"),
		}
	}

	for _, e := range el {
		if f(e) {
			return
		}
	}
}

func (g graphLiteralFixture) EachEdgeIncidentTo(v Vertex, f EdgeStep) {
	if g {
		switch v {
		case "foo":
			f(NewEdge("foo", "bar"))
		case "bar":
			terminate := f(NewEdge("foo", "bar"))
			if !terminate {
				f(NewEdge("bar", "baz"))
			}
		case "baz":
			f(NewEdge("bar", "baz"))
		default:
		}
	} else {
		switch v {
		case "foo":
			f(NewEdge("bar", "foo"))
		case "bar":
			terminate := f(NewEdge("bar", "foo"))
			if !terminate {
				f(NewEdge("baz", "bar"))
			}
		case "baz":
			f(NewEdge("baz", "bar"))
		default:
		}
	}
}

func (g graphLiteralFixture) EachArcFrom(v Vertex, f EdgeStep) {
	if g {
		switch v {
		case "foo":
			f(NewEdge("foo", "bar"))
		case "bar":
			f(NewEdge("bar", "baz"))
		default:
		}
	} else {
		switch v {
		case "bar":
			f(NewEdge("bar", "foo"))
		case "baz":
			f(NewEdge("baz", "bar"))
		default:
		}
	}
}

func (g graphLiteralFixture) EachArcTo(v Vertex, f EdgeStep) {
	if g {
		switch v {
		case "bar":
			f(NewEdge("foo", "bar"))
		case "baz":
			f(NewEdge("bar", "baz"))
		default:
		}
	} else {
		switch v {
		case "foo":
			f(NewEdge("bar", "foo"))
		case "bar":
			f(NewEdge("baz", "bar"))
		default:
		}
	}
}

func (g graphLiteralFixture) EachPredecessorOf(v Vertex, f VertexStep) {
	if g {
		switch v {
		case "bar":
			f("foo")
		case "baz":
			f("bar")
		default:
		}
	} else {
		switch v {
		case "foo":
			f("bar")
		case "bar":
			f("baz")
		default:
		}
	}
}
func (g graphLiteralFixture) EachSuccessorOf(v Vertex, f VertexStep) {
	if g {
		switch v {
		case "foo":
			f("bar")
		case "bar":
			f("baz")
		default:
		}
	} else {
		switch v {
		case "bar":
			f("foo")
		case "baz":
			f("bar")
		default:
		}
	}
}

func (g graphLiteralFixture) EachAdjacentTo(v Vertex, f VertexStep) {
	switch v {
	case "foo":
		f("bar")
	case "bar":
		terminate := f("foo")
		if !terminate {
			f("baz")
		}
	case "baz":
		f("bar")
	default:
	}
}

func (g graphLiteralFixture) HasVertex(v Vertex) bool {
	switch v {
	case "foo", "bar", "baz", "isolate":
		return true
	default:
		return false
	}
}

func (g graphLiteralFixture) InDegreeOf(v Vertex) (degree int, exists bool) {
	if g {
		switch v {
		case "foo":
			return 0, true
		case "bar":
			return 1, true
		case "baz":
			return 1, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}
	} else {
		switch v {
		case "foo":
			return 1, true
		case "bar":
			return 1, true
		case "baz":
			return 0, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}
	}
}

func (g graphLiteralFixture) OutDegreeOf(v Vertex) (degree int, exists bool) {
	if g {
		switch v {
		case "foo":
			return 1, true
		case "bar":
			return 1, true
		case "baz":
			return 0, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}

	} else {
		switch v {
		case "foo":
			return 0, true
		case "bar":
			return 1, true
		case "baz":
			return 1, true
		case "isolate":
			return 0, true
		default:
			return 0, false
		}
	}
}

func (g graphLiteralFixture) DegreeOf(v Vertex) (degree int, exists bool) {
	switch v {
	case "foo":
		return 1, true
	case "bar":
		return 2, true
	case "baz":
		return 1, true
	case "isolate":
		return 0, true
	default:
		return 0, false
	}
}

func (g graphLiteralFixture) HasEdge(e Edge) bool {
	u, v := e.Both()

	// TODO this is a little hinky until Arc is introduced
	switch u {
	case "foo":
		return v == "bar"
	case "bar":
		return v == "baz" || v == "foo"
	case "baz":
		return v == "bar"
	default:
		return false
	}
}

func (g graphLiteralFixture) Density() float64 {
	return 2 / 12 // 2 edges of maximum 12 in a 4-vertex digraph
}

func (g graphLiteralFixture) Transpose() Digraph {
	return graphLiteralFixture(!g)
}

func (g graphLiteralFixture) Size() int {
	return 2
}

func (g graphLiteralFixture) Order() int {
	return 4
}

// Tests for collection functors
type CollectionFunctorsSuite struct{}

var _ = Suite(&CollectionFunctorsSuite{})

func (s *CollectionFunctorsSuite) TestCollectVertices(c *C) {
	slice := CollectVertices(graphLiteralFixture(true))

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
	slice := CollectVerticesAdjacentTo("bar", graphLiteralFixture(true))

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, v := range slice {
		set.Add(v)
	}

	c.Assert(set.Has("foo"), Equals, true)
	c.Assert(set.Has("baz"), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectEdges(c *C) {
	slice := CollectEdges(graphLiteralFixture(true))

	c.Assert(len(slice), Equals, 2)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(NewEdge("foo", "bar")), Equals, true)
	c.Assert(set.Has(NewEdge("bar", "baz")), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectEdgesIncidentTo(c *C) {
	slice := CollectEdgesIncidentTo("foo", graphLiteralFixture(true))

	c.Assert(len(slice), Equals, 1)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(NewEdge("foo", "bar")), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectArcsFrom(c *C) {
	slice := CollectArcsFrom("foo", graphLiteralFixture(true))

	c.Assert(len(slice), Equals, 1)

	set := set.NewNonTS()
	for _, e := range slice {
		set.Add(e)
	}

	c.Assert(set.Has(NewEdge("foo", "bar")), Equals, true)
}

func (s *CollectionFunctorsSuite) TestCollectArcsTo(c *C) {
	slice := CollectArcsTo("bar", graphLiteralFixture(true))

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
	c.Assert(Order(graphFixtures["3e5v1i"]), Equals, 5)
	c.Assert(Order(graphLiteralFixture(true)), Equals, 4)
}

func (s *CountingFunctorsSuite) TestSize(c *C) {
	c.Assert(Size(graphFixtures["3e5v1i"]), Equals, 3)
	c.Assert(Size(graphLiteralFixture(true)), Equals, 2)
}
