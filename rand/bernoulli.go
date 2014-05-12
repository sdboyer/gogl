package rand

import (
	"github.com/sdboyer/gogl"
	stdrand "math/rand"
	"time"
)

// Generates a random graph of order n with Bernoulli distribution probability p of edge existing between any two vertices.
//
// This produces simple graphs only - no loops, no multiple edges. Graphs can be either directed or undirected, governed
// by the appropriately named parameter.
//
// If a stable graph is requested (stable == true), then the edge set presented by calling EachEdge() on the returned graph
// will be the same on every call. To provide stability, however, a memory allocation of n^2 * (int width) bytes
// is required to store the generated graph.
//
// Unstable graphs will create a new probabilistic edge set on the fly each time EachEdge(). It thus makes only minimal
// allocations, but is still CPU intensive for successive runs (and produces a different edge set). Given these
// characteristics, unstable graphs should always be used for single-use random graphs.
//
// Note that calling the Size() method on an unstable graph will create a prediction based on the Bernoulli number, but
// is not guaranteed to be exactly the same as the number of edges traversed through EachEdge().
func BernoulliDistribution(n uint, ρ int, directed bool, stable bool, src stdrand.Source) gogl.GraphEnumerator {
	if ρ > 100 || ρ < 1 {
		panic("Must designate an edge probability between 1 and 100")
	}

	if src == nil {
		src = stdrand.NewSource(time.Now().UnixNano())
	}

	if stable {
		return &stableBernoulliGraph{order: n, edgeProbability: ρ, directed: directed, source: src}
	} else {
		return unstableBernoulliGraph{order: n, edgeProbability: ρ, directed: directed, source: src}
	}
}

type stableBernoulliGraph struct {
	order uint
	size int
	directed bool
	edgeProbability int
	source stdrand.Source
	list [][]struct{}
}

func (g *stableBernoulliGraph) EachVertex(f gogl.VertexLambda) {
	o := int(g.order)
	for i := 0; i < o; i++ {
		if f(i) {
			return
		}
	}
}

func (g *stableBernoulliGraph) EachEdge(f gogl.EdgeLambda) {
	if g.list == nil {
		r := stdrand.New(g.source)
		g.list = make([][]struct{}, g.order, g.order)

		// Wrapping edge lambda; records edges into the adjacency list, then passes edge along
		ff := func (e gogl.Edge) bool {
			g.list[e.Source().(int)][e.Target().(int)] = struct{}{}
			g.size++
			return f(e)
		}

		if g.directed {
			bernoulliArcCreator(ff, int(g.order), g.edgeProbability, r)
		} else {
			bernoulliEdgeCreator(ff, int(g.order), g.edgeProbability, r)
		}
	} else {
		var e gogl.BaseEdge
		for u, adj := range g.list {
			for v, _ := range adj {
				e.U, e.V = u, v
				if f(e) {
					return
				}
			}
		}
	}
}

func (g *stableBernoulliGraph) Order() int {
	return int(g.order)
}

func (g *stableBernoulliGraph) Size() int {
	g.EachEdge(func (e gogl.Edge) (terminate bool) {
		return
	})
	return g.size
}

type unstableBernoulliGraph struct {
	order uint
	directed bool
	edgeProbability int
	source stdrand.Source
}

func (g unstableBernoulliGraph) EachVertex(f gogl.VertexLambda) {
	o := int(g.order)
	for i := 0; i < o; i++ {
		if f(i) {
			return
		}
	}
}

func (g unstableBernoulliGraph) EachEdge(f gogl.EdgeLambda) {
	if g.directed {
		bernoulliArcCreator(f, int(g.order), g.edgeProbability, stdrand.New(g.source))
	} else {
		bernoulliEdgeCreator(f, int(g.order), g.edgeProbability, stdrand.New(g.source))
	}
}

func (g unstableBernoulliGraph) Order() int {
	return int(g.order)
}

// The return value here is hogwash; as the generator is rerun with each passthrough, there is
// no guarantee the size will actually be exactly the same as the size produced by iterating EachEdge().
// It should be reasonably close...but rarely exactly correct, with the likelihood inversely proportional
// to the order of the graph.
func (g unstableBernoulliGraph) Size() int {
	var cs int

	cs = int(g.order) * (int(g.order) - 1)
	if !g.directed {
		cs = cs/2
	}

	return int(float64(cs) * (float64(g.edgeProbability) / 100))
}

var bernoulliEdgeCreator = func(el gogl.EdgeLambda, order int, edgeProbability int, r *stdrand.Rand) {
	var e gogl.BaseEdge
	for u := 0; u < order; u++ {
		// Set target vertex to one more than current source vertex. This guarantees
		// we only evaluate each unique edge pair once, as gogl's implicit contract requires.
		for v := u + 1; v < order; v++ {
			if r.Intn(100) <= edgeProbability {
				e.U, e.V = u, v
				if el(e) {
					return
				}
			}
		}
	}
}

var bernoulliArcCreator = func(el gogl.EdgeLambda, order int, edgeProbability int, r *stdrand.Rand) {
	var e gogl.BaseEdge
	for u := 0; u < order; u++ {
		for v := 0; v < order; v++ {
			if u != v && r.Intn(100) <= edgeProbability {
				e.U, e.V = u, v
				if el(e) {
					return
				}
			}
		}
	}
}

//func UnstableBernoulliDistribution(order uint, edgeProbability int, src stdrand.Source)
