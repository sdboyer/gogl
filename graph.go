// gogl provides a framework for representing and working with graphs.
package gogl

// gogl places *no* typing requirements on vertices. However, this
// does mean that calling code will need to directly convert vertices
// emitted from gogl methods back into their desired types.
type Vertex interface{}

/* Edge structures */

// A graph's behaviors is primarily a product of the constraints and
// capabilities it places on its edges. These constraints and
// capabilities determine whether certain types of operations are
// possible on the graph, as well as the efficiencies for various
// operations.

// gogl aims to provide a diverse range of graph implementations that
// can meet the varying constraints and implementation needs, but still
// achieve optimal performance given those constraints.

// TODO totally unclear whether or not defining capabilities in a
// bitfield like this will actually help us achieve the goal
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
	//  Properties() uint
}

// BaseEdge is a struct used internally to represent edges and meet the
// Edge interface requirements. It uses the standard notation, (u,v), for
// vertex pairs in an edge.
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

/* Graph structures */

type Graph interface {
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(edge Edge))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
	Order() uint
	Size() uint
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
