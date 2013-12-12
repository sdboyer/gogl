// gogl provides a framework for representing and working with graphs.
package gogl

// Constants defining graph capabilities and behaviors.
const (
  E_DIRECTED, EM_DIRECTED = 1 << iota, 1 << iota - 1
  E_UNDIRECTED, EM_UNDIRECTED
  E_WEIGHTED, EM_WEIGHTED
  E_TYPED, EM_TYPED
  E_SIGNED, EM_SIGNED
  E_LOOPS, EM_LOOPS
  E_MULTIGRAPH, EM_MULTIGRAPH
)

type Vertex interface{}

type Edge struct {
	Tail, Head Vertex
}

type Graph interface {
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(edge Edge))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
	Order() uint
	Size() uint
	GetSubgraph([]Vertex) Graph
}

type MutableGraph interface {
	Graph
	AddVertex(v interface{}) bool
	RemoveVertex(v interface{}) bool
}

type DirectedGraph interface {
	Graph
	Transpose() DirectedGraph
	IsAcyclic() bool
	GetCycles() [][]interface{}
}

type MutableDirectedGraph interface {
	MutableGraph
	DirectedGraph
	addDirectedEdge(source interface{}, target interface{}) bool
	removeDirectedEdge(source interface{}, target interface{}) bool
}
