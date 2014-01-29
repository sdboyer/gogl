package adjacency_list

import (
	"fmt"
	. "github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/test_bundle"
	//"math"
	"testing"
)

var _ = fmt.Println

var factory = test_bundle.GraphFactory{
	CreateMutableGraph: func() MutableGraph {
		return NewUndirected()
	},
	CreateGraph: func(edges []Edge) Graph {
		return NewUndirectedFromEdgeSet(edges)
	},
}

func Test_UVertexMembership(t *testing.T) {
	test_bundle.GraphTestVertexMembership(factory, t)
}
