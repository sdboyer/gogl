package dfs

import (
	"testing"

	"github.com/sdboyer/gogl"
)

var dfEdgeSet = []gogl.Edge{
	&gogl.BaseEdge{"foo", "bar"},
	&gogl.BaseEdge{"bar", "baz"},
	&gogl.BaseEdge{"baz", "qux"},
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
	g := gogl.NewDirectedFromEdgeSet(dfEdgeSet)

	vis := &DFTslVisitor{}
	DepthFirstFromVertices(g, vis, "foo")

	if !sliceEquals(vis.GetTsl(), []gogl.Vertex{"qux", "baz", "bar", "foo"}) {
		t.Error("TSL is not correct.")
	}
}
