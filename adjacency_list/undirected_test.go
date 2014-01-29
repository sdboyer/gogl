package adjacency_list

import (
	"fmt"
	. "github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/test_bundle"
	//"math"
	"testing"
)

var _ = fmt.Println

var u_fact = &GraphFactory{
	CreateMutableGraph: func() MutableGraph {
		return NewUndirected()
	},
	CreateGraph: func(edges []Edge) Graph {
		return NewUndirectedFromEdgeSet(edges)
	},
}
func Test_UVertexMembership(t *testing.T) {
	test_bundle.GraphTestVertexMembership(u_fact, t)
}

func Test_UVertexMultiOps(t *testing.T) {
	test_bundle.GraphTestVertexMultiOps(u_fact, t)
}

func Test_UVertexRemoveVertexWithEdges(t *testing.T) {
	test_bundle.GraphTestRemoveVertexWithEdges(u_fact, t)
}

func Test_UVertexTestEachVertex(t *testing.T) {
	test_bundle.GraphTestEachVertex(u_fact, t)
}

