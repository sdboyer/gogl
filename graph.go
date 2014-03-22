package gogl

/* Vertex structures */

// As a rule, gogl tries to place as low a requirement on its vertices as
// possible. This is because, from a purely graph theoretic perspective,
// vertices are inert. Boring, even. Graphs are more about the topology, the
// characteristics of the edges connecting the points than the points
// themselves. Your use case cares about the content of your vertices, but gogl
// does not.  Consequently, anything can act as a vertex.
type Vertex interface{}

/* Edge structures */

// A graph's behaviors is primarily a product of the constraints and
// capabilities it places on its edges. These constraints and capabilities
// determine whether certain types of operations are possible on the graph, as
// well as the efficiencies for various operations.

// gogl aims to provide a diverse range of graph implementations that can meet
// the varying constraints and implementation needs, but still achieve optimal
// performance given those constraints.

type Edge interface {
	Source() Vertex
	Target() Vertex
	Both() (Vertex, Vertex)
}

type WeightedEdge interface {
	Edge
	Weight() int
}

// BaseEdge is a struct used internally to represent edges and meet the Edge
// interface requirements. It uses the standard notation, (u,v), for vertex
// pairs in an edge.
type BaseEdge struct {
	U Vertex
	V Vertex
}

func (e BaseEdge) Source() Vertex {
	return e.U
}

func (e BaseEdge) Target() Vertex {
	return e.V
}

func (e BaseEdge) Both() (Vertex, Vertex) {
	return e.U, e.V
}

type BaseWeightedEdge struct {
	BaseEdge
	W int
}

func (e BaseWeightedEdge) Weight() int {
	return e.W
}

/* Graph structures */

type Graph interface {
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(edge Edge))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
	HasEdge(e Edge) bool
	Order() int
	Size() int
	InDegree(vertex Vertex) (int, bool)
	OutDegree(vertex Vertex) (int, bool)
}

type MutableGraph interface {
	Graph
	EnsureVertex(vertices ...Vertex)
	RemoveVertex(vertices ...Vertex)
	AddEdges(edges ...Edge)
	RemoveEdges(edges ...Edge)
}

// A simple graph is in opposition to a multigraph: it disallows loops and
// parallel edges.
type SimpleGraph interface {
	Graph
	Density() float64
}

type DirectedGraph interface {
	Graph
	Transpose() DirectedGraph
}

type WeightedGraph interface {
	Graph
	HasWeightedEdge(e WeightedEdge) bool
	EachWeightedEdge(f func(edge WeightedEdge))
}

type MutableWeightedGraph interface {
	WeightedGraph
	EnsureVertex(vertices ...Vertex)
	RemoveVertex(vertices ...Vertex)
	AddEdges(edges ...WeightedEdge)
	RemoveEdges(edges ...WeightedEdge)
}
