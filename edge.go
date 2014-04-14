package gogl

// A graph's behaviors are primarily a product of the constraints and
// capabilities it places on its edges. These constraints and capabilities
// determine whether certain types of operations are possible on the graph, as
// well as the efficiencies for various operations.

// gogl aims to provide a range of graph implementations that can meet
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

type LabeledEdge interface {
	Edge
	Label() string
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

type BaseLabeledEdge struct {
	BaseEdge
	L string
}

func (e BaseLabeledEdge) Label() string {
	return e.L
}
