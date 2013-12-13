// gogl provides a framework for representing and working with graphs.
package gogl

// Constants defining graph capabilities and behaviors.
const (
	E_DIRECTED, EM_DIRECTED = 1 << iota, 1<<iota - 1
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
	AddVertex(v Vertex) bool
	RemoveVertex(v Vertex) bool
	AddEdge(edge Edge) (bool, error)
}

type DirectedGraph interface {
	Graph
	Transpose() DirectedGraph
	IsAcyclic() bool
	GetCycles() [][]Vertex
	addDirectedEdge(source Vertex, target Vertex) bool
	removeDirectedEdge(source Vertex, target Vertex) bool
}
