// gogl provides a framework for representing and working with graphs.
package gogl

type Vertex interface{}

type Graph interface {
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(source Vertex, target Vertex))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
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

type Edge interface {
	Tail() Vertex
	Head() Vertex
}
