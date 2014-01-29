package adjacency_list

import (
	"fmt"
	. "github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/test_bundle"
	//"math"
	"testing"
)

var fml = fmt.Println

func newUndirected(edges ...Edge) (g MutableGraph) {
	if len(edges) > 0 {
		g = NewUndirectedFromEdgeSet(edges)
	} else {
		g = NewUndirected()
	}
	return g
}

func newDirected(edges ...Edge) (g MutableGraph) {
	if len(edges) > 0 {
		g = NewDirectedFromEdgeSet(edges)
	} else {
		g = NewDirected()
	}
	return g
}

func Test_UVertexMembership(t *testing.T) {
	test_bundle.GraphTestVertexMembership(newUndirected, t)
}
