package dfs

import (
	"fmt"
	"testing"

	"github.com/sdboyer/gogl"
	. "launchpad.net/gocheck"
)

// Hook gocheck into the go test runner
func Test(t *testing.T) { TestingT(t) }

var dfEdgeSet = []gogl.Edge{
	&gogl.BaseEdge{"foo", "bar"},
	&gogl.BaseEdge{"bar", "baz"},
	&gogl.BaseEdge{"baz", "qux"},
}

type DepthFirstSearchSuite struct {
}

var _ = Suite(&DepthFirstSearchSuite{})

// Basic test of outermost search functionality.
func (s *DepthFirstSearchSuite) TestSearch(c *C) {
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)

	path, err := Search(g, "qux", "bar")
	c.Assert(fmt.Sprint(path), Equals, fmt.Sprint([]gogl.Vertex{"qux", "baz", "bar"}))
	c.Assert(err, IsNil)
}

func sliceEquals(a, b []gogl.Vertex) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestTslGeneration(t *testing.T) {
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)

	vis := &DFTslVisitor{}
	DepthFirstFromVertices(g, vis, "foo")

	if !sliceEquals(vis.GetTsl(), []gogl.Vertex{"qux", "baz", "bar", "foo"}) {
		t.Error("TSL is not correct.")
	}
}
