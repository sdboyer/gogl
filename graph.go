// gogl provides a framework for representing and working with graphs.
package gogl

/* Vertex structures */

// As a rule, gogl tries to place as low a requirement on its vertices as
// possible. This is because, from a purely graph theoretic perspective,
// vertices are inert. Boring, even. Graphs are more about the topology, the
// characteristics of the edges connecting the points than the points
// themselves. Your use case cares about the content of your vertices, but gogl
// does not.  Consequently, anything can act as a vertex.
type Vertex interface {}

// However, as a practical, internal matter, some approaches to representing
// graphs may benefit significantly (in terms of memory use) from having a uint
// that uniquely identifies each vertex. This is essentially an implementation
// detail, but if your use case merits it (large graph processing is busting
// memory), you can trade off a bit of speed for potentially nontrivial memory
// savings.
// TODO implement a datastructure that does this, and test these outlandish claims!
type IdentifiableVertex interface {
	id() uint64
}


/* Edge structures */

// A graph's behaviors is primarily a product of the constraints and
// capabilities it places on its edges. These constraints and capabilities
// determine whether certain types of operations are possible on the graph, as
// well as the efficiencies for various operations.

// gogl aims to provide a diverse range of graph implementations that can meet
// the varying constraints and implementation needs, but still achieve optimal
// performance given those constraints.

// TODO totally unclear whether or not defining capabilities in a bitfield like
// this will actually help us achieve the goal
const (
	E_DIRECTED, EM_DIRECTED = 1 << iota, 1<<iota - 1
	E_UNDIRECTED, EM_UNDIRECTED
	E_WEIGHTED, EM_WEIGHTED
	E_TYPED, EM_TYPED
	E_SIGNED, EM_SIGNED
	E_LOOPS, EM_LOOPS
	E_MULTIGRAPH, EM_MULTIGRAPH
)

type Edge interface {
	Source() Vertex
	Target() Vertex
	Both() (Vertex, Vertex)
	//  Properties() uint
}

// BaseEdge is a struct used internally to represent edges and meet the Edge
// interface requirements. It uses the standard notation, (u,v), for vertex
// pairs in an edge.
// TODO what does having non-exported fields, and no write methods, mean
// for mutability outside of the package?
type BaseEdge struct {
	u Vertex
	v Vertex
}

func (e BaseEdge) Source() Vertex {
	return e.u
}

func (e BaseEdge) Target() Vertex {
	return e.v
}

func (e BaseEdge) Both() (Vertex, Vertex) {
	return e.u, e.v
}

/* Graph structures */

type Graph interface {
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(edge Edge))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
	Order() uint
	Size() uint
	InDegree(vertex Vertex) (uint, bool)
	OutDegree(vertex Vertex) (uint, bool)
	AddVertex(v Vertex) bool
	RemoveVertex(v Vertex) bool
	AddEdge(edge Edge) bool
}

// A simple graph is in opposition to a multigraph: it disallows loops
// and parallel edges.
type SimpleGraph interface {
	Density() float64
}

type DirectedGraph interface {
	Graph
	Transpose() DirectedGraph
	IsAcyclic() bool
	GetCycles() [][]Vertex
	addDirectedEdge(source Vertex, target Vertex) bool
	removeDirectedEdge(source Vertex, target Vertex) bool
}
