package gogl

/* Graph creation */

// These interfaces are convenient shorthand
type gm interface {
	EdgeSetMutator
	VertexSetMutator
}

type wgm interface {
	WeightedEdgeSetMutator
	VertexSetMutator
}

type pgm interface {
	DataEdgeSetMutator
	VertexSetMutator
}

type lgm interface {
	LabeledEdgeSetMutator
	VertexSetMutator
}

// Copies the first graph's edges and vertices into the second, and returns the second.
//
// The second argument must be some flavor of MutableGraph, otherwise this function will panic.
//
// Generally speaking, this is a very naive copy operation; it should not be relied on for complex cases.
func CopyGraph(from Graph, to interface{}) interface{} {
	var el EdgeLambda
	var g VertexSetMutator

	// Establish the mutable type of the second graph, then dispatch to
	// a specialized copyer depending on whether the types correspond.
	if g, ok := to.(gm); ok {
		el = func(e Edge) (terminate bool) {
			g.AddEdges(e)
			return
		}
	} else if g, ok := to.(wgm); ok {
		el = func(e Edge) (terminate bool) {
			if ee, ok := e.(WeightedEdge); ok {
				g.AddEdges(ee)
			} else {
				// TODO should this case panic?
				g.AddEdges(BaseWeightedEdge{BaseEdge{U: e.Source(), V: e.Target()}, 0})
			}
			return
		}
	} else if g, ok := to.(lgm); ok {
		el = func(e Edge) (terminate bool) {
			if ee, ok := e.(LabeledEdge); ok {
				g.AddEdges(ee)
			} else {
				// TODO should this case panic?
				g.AddEdges(BaseLabeledEdge{BaseEdge{U: e.Source(), V: e.Target()}, ""})
			}
			return
		}
	} else if g, ok := to.(pgm); ok {
		el = func(e Edge) (terminate bool) {
			if ee, ok := e.(DataEdge); ok {
				g.AddEdges(ee)
			} else {
				// TODO should this case panic?
				g.AddEdges(BaseDataEdge{BaseEdge{U: e.Source(), V: e.Target()}, struct{}{}})
			}
			return
		}
	} else {
		panic("Second graph passed to CopyGraph must be mutable.")
	}

	// Do the simplistic copy
	from.EachEdge(el)

	// Ensure vertex isolates come, too
	from.EachVertex(func(v Vertex) (terminate bool) {
		g.EnsureVertex(v)
		return
	})

	return g
}
