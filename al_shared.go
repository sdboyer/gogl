package gogl

// Copies an incoming graph into any of the implemented adjacency list types.
//
// This encapsulates the full matrix of conversion possibilities between
// different graph edge types.
func functorToAdjacencyList(from Graph, to interface{}) {
	vf := func(from Graph, to al_mutver) {
		if to.Order() != from.Order() {
			from.EachVertex(func(vertex Vertex) (terminate bool) {
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
				g.addEdges(BaseWeightedEdge{BaseEdge{U: e.Source(), V: e.Target()}, 0})
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_lea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(LabeledEdge); ok {
				g.addEdges(e)
			} else {
				g.addEdges(BaseLabeledEdge{BaseEdge{U: e.Source(), V: e.Target()}, ""})
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_pea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(PropertyEdge); ok {
				g.addEdges(e)
			} else {
				g.addEdges(BasePropertyEdge{BaseEdge{U: e.Source(), V: e.Target()}, struct{}{}})
			}
			return
		})
		vf(from, g)
	} else {
		panic("Target graph did not implement a recognized adjacency list internal type")
	}

}

type al_mutver interface {
	ensureVertex(...Vertex)
	Order() int
}

type al_ea interface {
	al_mutver
	addEdges(...Edge)
}

type al_wea interface {
	al_mutver
	addEdges(...WeightedEdge)
}

type al_lea interface {
	al_mutver
	addEdges(...LabeledEdge)
}

type al_pea interface {
	al_mutver
	addEdges(...PropertyEdge)
}
