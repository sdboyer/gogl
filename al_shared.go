package gogl

// Wraps an EdgeLambda into a closure for deferred execution.
//
// This encapsulates the full matrix of conversion possibilities between
// different graph edge types.
func createDeferredEdgeLambda(from Graph, to interface{}) func() {
	var f func()
	if g, ok := to.(al_ea); ok {
		f = func() {
			from.EachEdge(func(edge Edge) (terminate bool) {
				g.addEdges(edge)
				return
			})
		}
	} else if g, ok := to.(al_wea); ok {
		f = func() {
			from.EachEdge(func(edge Edge) (terminate bool) {
				if e, ok := edge.(WeightedEdge); ok {
					g.addEdges(e)
				} else {
					g.addEdges(BaseWeightedEdge{BaseEdge{U: e.Source(), V: e.Target()}, 0})
				}
				return
			})
		}
	} else if g, ok := to.(al_lea); ok {
		f = func() {
			from.EachEdge(func(edge Edge) (terminate bool) {
				if e, ok := edge.(LabeledEdge); ok {
					g.addEdges(e)
				} else {
					g.addEdges(BaseLabeledEdge{BaseEdge{U: e.Source(), V: e.Target()}, ""})
				}
				return
			})
		}
	} else if g, ok := to.(al_pea); ok {
		f = func() {
			from.EachEdge(func(edge Edge) (terminate bool) {
				if e, ok := edge.(PropertyEdge); ok {
					g.addEdges(e)
				} else {
					g.addEdges(BasePropertyEdge{BaseEdge{U: e.Source(), V: e.Target()}, struct{}{}})
				}
				return
			})
		}
	}

	return f
}

type al_ea interface {
	addEdges(...Edge)
}

type al_wea interface {
	addEdges(...WeightedEdge)
}

type al_lea interface {
	addEdges(...LabeledEdge)
}

type al_pea interface {
	addEdges(...PropertyEdge)
}

