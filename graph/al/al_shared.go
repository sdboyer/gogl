package al

import (
	. "github.com/sdboyer/gogl"
)

// Contains behaviors shared across adjacency list implementations.

type al_graph interface {
	Graph
	ensureVertex(...Vertex)
	hasVertex(Vertex) bool
}

type al_digraph interface {
	Digraph
	ensureVertex(...Vertex)
	hasVertex(Vertex) bool
}

type al_ea interface {
	al_graph
	addEdges(...Edge)
}

type al_dea interface {
	al_digraph
	addArcs(...Arc)
}

type al_wea interface {
	al_graph
	addEdges(...WeightedEdge)
}

type al_dwea interface {
	al_digraph
	addArcs(...WeightedArc)
}

type al_lea interface {
	al_graph
	addEdges(...LabeledEdge)
}

type al_dlea interface {
	al_digraph
	addArcs(...LabeledArc)
}

type al_pea interface {
	al_graph
	addEdges(...DataEdge)
}

type al_dpea interface {
	al_digraph
	addArcs(...DataArc)
}

// Copies an incoming graph into any of the implemented adjacency list types.
//
// This encapsulates the full matrix of conversion possibilities between
// different graph edge types, for undirected graphs.
func functorToAdjacencyList(from GraphSource, to al_graph) Graph {
	vf := func(from GraphSource, to al_graph) {
		if Order(to) != Order(from) {
			from.Vertices(func(vertex Vertex) (terminate bool) {
				to.ensureVertex(vertex)
				return
			})

		}
	}

	if g, ok := to.(al_ea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			g.addEdges(edge)
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_wea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(WeightedEdge); ok {
				g.addEdges(e)
			} else {
				u, v := edge.Both()
				g.addEdges(NewWeightedEdge(u, v, 0))
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_lea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(LabeledEdge); ok {
				g.addEdges(e)
			} else {
				u, v := edge.Both()
				g.addEdges(NewLabeledEdge(u, v, ""))
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_pea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(DataEdge); ok {
				g.addEdges(e)
			} else {
				u, v := edge.Both()
				g.addEdges(NewDataEdge(u, v, nil))
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_ea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			g.addEdges(edge)
			return
		})
		vf(from, g)
	} else {
		panic("Target graph did not implement a recognized adjacency list internal type")
	}

	return to.(Graph)
}

// Copies an incoming graph into any of the implemented adjacency list types.
//
// This encapsulates the full matrix of conversion possibilities between
// different graph edge types, for directed graphs.
func functorToDirectedAdjacencyList(from DigraphSource, to al_digraph) Digraph {
	vf := func(from GraphSource, to al_graph) {
		if Order(to) != Order(from) {
			from.Vertices(func(vertex Vertex) (terminate bool) {
				to.ensureVertex(vertex)
				return
			})

		}
	}

	if g, ok := to.(al_dea); ok {
		from.EachArc(func(arc Arc) (terminate bool) {
			g.addArcs(arc)
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_dwea); ok {
		from.EachArc(func(arc Arc) (terminate bool) {
			if e, ok := arc.(WeightedArc); ok {
				g.addArcs(e)
			} else {
				g.addArcs(NewWeightedArc(arc.Source(), arc.Target(), 0))
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_dlea); ok {
		from.EachArc(func(arc Arc) (terminate bool) {
			if e, ok := arc.(LabeledArc); ok {
				g.addArcs(e)
			} else {
				g.addArcs(NewLabeledArc(arc.Source(), arc.Target(), ""))
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_dpea); ok {
		from.EachArc(func(arc Arc) (terminate bool) {
			if e, ok := arc.(DataArc); ok {
				g.addArcs(e)
			} else {
				g.addArcs(NewDataArc(arc.Source(), arc.Target(), nil))
			}
			return
		})
		vf(from, g)
	} else {
		panic("Target graph did not implement a recognized adjacency list internal type")
	}

	return to.(Digraph)
}

func eachVertexInAdjacencyList(list interface{}, vertex Vertex, vs VertexStep) {
	switch l := list.(type) {
	case map[Vertex]map[Vertex]struct{}:
		if _, exists := l[vertex]; exists {
			for adjacent := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	case map[Vertex]map[Vertex]float64:
		if _, exists := l[vertex]; exists {
			for adjacent := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	case map[Vertex]map[Vertex]string:
		if _, exists := l[vertex]; exists {
			for adjacent := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	case map[Vertex]map[Vertex]interface{}:
		if _, exists := l[vertex]; exists {
			for adjacent := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	default:
		panic("Unrecognized adjacency list map type.")
	}
}

func eachPredecessorOf(list interface{}, vertex Vertex, vs VertexStep) {
	switch l := list.(type) {
	case map[Vertex]map[Vertex]struct{}:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	case map[Vertex]map[Vertex]float64:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	case map[Vertex]map[Vertex]string:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	case map[Vertex]map[Vertex]interface{}:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	default:
		panic("Unrecognized adjacency list map type.")
	}

}

func inDegreeOf(g al_digraph, v Vertex) (degree int, exists bool) {
	if exists = g.hasVertex(v); exists {
		g.EachArc(func(e Arc) (terminate bool) {
			if v == e.Target() {
				degree++
			}
			return
		})
	}
	return
}

func eachEdgeIncidentToDirected(g al_digraph, v Vertex, f EdgeStep) {
	if !g.hasVertex(v) {
		return
	}

	var terminate bool
	interloper := func(e Arc) bool {
		terminate = terminate || f(e)
		return terminate
	}

	g.EachArcFrom(v, interloper)
	g.EachArcTo(v, interloper)
}
