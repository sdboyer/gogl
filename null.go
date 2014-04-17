package gogl

// The null graph is a graph without any edges or vertices. It implements all possible (non-mutable) graph interfaces.
//
// In effect, it is the zero-value of all possible graph types.
const NullGraph = nullGraph(false)

type nullGraph bool

var _ Graph = nullGraph(false)
var _ DirectedGraph = nullGraph(false)
var _ SimpleGraph = nullGraph(false)
var _ WeightedGraph = nullGraph(false)
var _ LabeledGraph = nullGraph(false)
var _ PropertyGraph = nullGraph(false)

func (g nullGraph) EachVertex(f func(v Vertex))                        {}
func (g nullGraph) EachEdge(f func(e Edge))                            {}
func (g nullGraph) EachWeightedEdge(f func(e WeightedEdge))            {}
func (g nullGraph) EachLabeledEdge(f func(e LabeledEdge))              {}
func (g nullGraph) EachPropertyEdge(f func(e PropertyEdge))            {}
func (g nullGraph) EachAdjacent(start Vertex, f func(adjacent Vertex)) {}

func (g nullGraph) HasVertex(v Vertex) bool {
	return false
}

func (g nullGraph) InDegree(Vertex) (degree int, exists bool) {
	return 0, false
}

func (g nullGraph) OutDegree(Vertex) (degree int, exists bool) {
	return 0, false
}

func (g nullGraph) HasEdge(e Edge) bool {
	return false
}

func (g nullGraph) HasWeightedEdge(e WeightedEdge) bool {
	return false
}

func (g nullGraph) HasLabeledEdge(e LabeledEdge) bool {
	return false
}

func (g nullGraph) HasPropertyEdge(e PropertyEdge) bool {
	return false
}

func (g nullGraph) Size() int {
	return 0
}

func (g nullGraph) Order() int {
	return 0
}

func (g nullGraph) Density() float64 {
	return 0
}

func (g nullGraph) Transpose() DirectedGraph {
	return g
}
