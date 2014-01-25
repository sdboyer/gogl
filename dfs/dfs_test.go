package gogl

import (
	"testing"
	al "github.com/sdboyer/gogl/adjacency_list"
	. "github.com/sdboyer/gogl"
)

var dfEdgeSet = []Edge{
	&BaseEdge{"foo", "bar"},
	&BaseEdge{"bar", "baz"},
	&BaseEdge{"baz", "qux"},
}

func sliceEquals(a, b []Vertex) bool {
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
	g := al.NewDirectedFromEdgeSet(dfEdgeSet)

	vis := &DFTslVisitor{}
	DepthFirstFromVertices(g, vis, "foo")

	if !sliceEquals(vis.GetTsl(), []Vertex{"qux", "baz", "bar", "foo"}) {
		t.Error("TSL is not correct.")
	}
}
