package al

import (
	"sync"
	. "github.com/sdboyer/gogl"
)

/*
Adjacency lists are a relatively simple graph representation. They maintain
a list of vertices, storing information about edge membership relative to
those vertices. This makes vertex-centric operations generally more
efficient, and edge-centric operations generally less efficient, as edges
are represented implicitly. It also makes them inappropriate for more
complex graph types, such as multigraphs.

gogl's adjacency lists are space-efficient; in a directed graph, the memory
cost for the entire graph G is proportional to V + E; in an undirected graph,
it is V + 2E.
*/

var alCreators = map[GraphProperties]func() Graph{
	GraphProperties(G_IMMUTABLE | G_DIRECTED | G_BASIC | G_SIMPLE): func() Graph {
		return &immutableDirected{al_basic_immut{al_basic{list: make(map[Vertex]map[Vertex]struct{})}}}
	},
	GraphProperties(G_MUTABLE | G_DIRECTED | G_BASIC | G_SIMPLE): func() Graph {
		return &mutableDirected{al_basic_mut{al_basic{list: make(map[Vertex]map[Vertex]struct{})}, sync.RWMutex{}}}
	},
	GraphProperties(G_MUTABLE | G_UNDIRECTED | G_BASIC | G_SIMPLE): func() Graph {
		return &mutableUndirected{al_basic_mut{al_basic{list: make(map[Vertex]map[Vertex]struct{})}, sync.RWMutex{}}}
	},
	GraphProperties(G_MUTABLE | G_DIRECTED | G_WEIGHTED | G_SIMPLE): func() Graph {
		return &weightedDirected{baseWeighted{list: make(map[Vertex]map[Vertex]float64), size: 0, mu: sync.RWMutex{}}}
	},
	GraphProperties(G_MUTABLE | G_UNDIRECTED | G_WEIGHTED | G_SIMPLE): func() Graph {
		return &weightedUndirected{baseWeighted{list: make(map[Vertex]map[Vertex]float64), size: 0, mu: sync.RWMutex{}}}
	},
	GraphProperties(G_MUTABLE | G_DIRECTED | G_LABELED | G_SIMPLE): func() Graph {
		return &labeledDirected{baseLabeled{list: make(map[Vertex]map[Vertex]string), size: 0, mu: sync.RWMutex{}}}
	},
	GraphProperties(G_MUTABLE | G_UNDIRECTED | G_LABELED | G_SIMPLE): func() Graph {
		return &labeledUndirected{baseLabeled{list: make(map[Vertex]map[Vertex]string), size: 0, mu: sync.RWMutex{}}}
	},
	GraphProperties(G_MUTABLE | G_DIRECTED | G_DATA | G_SIMPLE): func() Graph {
		return &dataDirected{baseData{list: make(map[Vertex]map[Vertex]interface{}), size: 0, mu: sync.RWMutex{}}}
	},
	GraphProperties(G_MUTABLE | G_UNDIRECTED | G_DATA | G_SIMPLE): func() Graph {
		return &dataUndirected{baseData{list: make(map[Vertex]map[Vertex]interface{}), size: 0, mu: sync.RWMutex{}}}
	},
}

// Create a graph implementation in the adjacency list style from the provided GraphSpec.
//
// If the GraphSpec contains a GraphSource, it will be imported into the provided graph.
// If the GraphSpec indicates a graph type that is not currently implemented, this function
// will panic.
func G(gs GraphSpec) Graph {
	for gp, gf := range alCreators {

		// TODO satisfiability here is not so narrow
		if gp&^gs.Props == 0 {
			if gs.Source != nil {
				if gs.Props&G_DIRECTED == G_DIRECTED {
					if dgs, ok := gs.Source.(DigraphSource); ok {
						return functorToDirectedAdjacencyList(dgs, gf().(al_digraph))
					} else {
						panic("Cannot create a digraph from a graph.")
					}
				} else {
					return functorToAdjacencyList(gs.Source, gf().(al_graph))
				}
			} else {
				return gf()
			}
		}
	}

	panic("No graph implementation found for spec")
}

type al_basic struct {
	list map[Vertex]map[Vertex]struct{}
	size int
}

// Helper to not have to write struct{} everywhere.
var keyExists = struct{}{}

// Indicates whether or not the given vertex is present in the graph.
func (g *al_basic) hasVertex(vertex Vertex) (exists bool) {
	_, exists = g.list[vertex]
	return
}

// Returns the size (number of edges) in the graph.
func (g *al_basic) Size() int {
	return g.size
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *al_basic) ensureVertex(vertices ...Vertex) {
	for _, vertex := range vertices {
		if !g.hasVertex(vertex) {
			// TODO experiment with different lengths...possibly by analyzing existing density?
			g.list[vertex] = make(map[Vertex]struct{}, 10)
		}
	}

	return
}

type al_basic_immut struct {
	al_basic
}

// Traverses the graph's vertices in random order, passing each vertex to the
// provided closure.
func (g *al_basic_immut) EachVertex(f VertexStep) {
	for v := range g.list {
		if f(v) {
			return
		}
	}
}

// Indicates whether or not the given vertex is present in the graph.
func (g *al_basic_immut) HasVertex(vertex Vertex) bool {
	return g.hasVertex(vertex)
}

// Returns the order (number of vertices) in the graph.
func (g *al_basic_immut) Order() int {
	return len(g.list)
}

type al_basic_mut struct {
	al_basic
	mu sync.RWMutex
}

/* Base al_basic_mut methods */

// Traverses the graph's vertices in random order, passing each vertex to the
// provided closure.
func (g *al_basic_mut) EachVertex(f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for v := range g.list {
		if f(v) {
			return
		}
	}
}

// Indicates whether or not the given vertex is present in the graph.
func (g *al_basic_mut) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	exists = g.hasVertex(vertex)
	return
}

// Returns the order (number of vertices) in the graph.
func (g *al_basic_mut) Order() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.list)
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *al_basic_mut) EnsureVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.ensureVertex(vertices...)
}
