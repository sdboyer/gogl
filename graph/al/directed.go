package al

import (
	. "github.com/sdboyer/gogl"
)

type mutableDirected struct {
	al_basic_mut
}

/* mutableDirected additions */

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *mutableDirected) OutDegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Returns the indegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
//
// Note that getting indegree is inefficient for directed adjacency lists; it requires
// a full scan of the graph's edge set.
func (g *mutableDirected) InDegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return inDegreeOf(g, vertex)
}

// Returns the degree of the provided vertex, counting both in and out-edges.
func (g *mutableDirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	indegree, exists := inDegreeOf(g, vertex)
	outdegree, exists := g.OutDegreeOf(vertex)
	return indegree + outdegree, exists
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *mutableDirected) EachEdge(f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			if f(NewEdge(source, target)) {
				return
			}
		}
	}
}

// Traverses the set of arcs in the graph, passing each arc to the
// provided closure.
func (g *mutableDirected) EachArc(f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			if f(NewArc(source, target)) {
				return
			}
		}
	}
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *mutableDirected) EachEdgeIncidentTo(v Vertex, f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	eachEdgeIncidentToDirected(g, v, f)
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *mutableDirected) EachAdjacentTo(start Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.EachEdgeIncidentTo(start, func(e Edge) bool {
		u, v := e.Both()
		if u == start {
			return f(v)
		} else {
			return f(u)
		}
	})
}

// Enumerates the set of out-edges for the provided vertex.
func (g *mutableDirected) EachArcFrom(v Vertex, f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, _ := range g.list[v] {
		if f(NewArc(v, adjacent)) {
			return
		}
	}
}

func (g *mutableDirected) EachSuccessorOf(v Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachVertexInAdjacencyList(g.list, v, f)
}

// Enumerates the set of in-edges for the provided vertex.
func (g *mutableDirected) EachArcTo(v Vertex, f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for candidate, adjacent := range g.list {
		for target, _ := range adjacent {
			if target == v {
				if f(NewArc(candidate, target)) {
					return
				}
			}
		}
	}
}

func (g *mutableDirected) EachPredecessorOf(v Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachPredecessorOf(g.list, v, f)
}

// Indicates whether or not the given edge is present in the graph.
func (g *mutableDirected) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	u, v := edge.Both()
	_, exists := g.list[u][v]
	return exists
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *mutableDirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *mutableDirected) RemoveVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, vertex := range vertices {
		if g.hasVertex(vertex) {
			// TODO Is the expensive search good to do here and now...
			// while read-locked?
			g.size -= len(g.list[vertex])
			delete(g.list, vertex)

			// TODO consider chunking the list and parallelizing into goroutines
			for _, adjacent := range g.list {
				if _, has := adjacent[vertex]; has {
					delete(adjacent, vertex)
					g.size--
				}
			}
		}
	}
	return
}

// Adds arcs to the graph.
func (g *mutableDirected) AddArcs(arcs ...Arc) {
	if len(arcs) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addArcs(arcs...)
}

// Adds a new arc to the graph.
func (g *mutableDirected) addArcs(arcs ...Arc) {
	for _, arc := range arcs {
		g.ensureVertex(arc.Source(), arc.Target())

		if _, exists := g.list[arc.Source()][arc.Target()]; !exists {
			g.list[arc.Source()][arc.Target()] = keyExists
			g.size++
		}
	}
}

// Removes arcs from the graph. This does NOT remove vertex members of the
// removed arcs.
func (g *mutableDirected) RemoveArcs(arcs ...Arc) {
	if len(arcs) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, arc := range arcs {
		s, t := arc.Both()
		if _, exists := g.list[s][t]; exists {
			delete(g.list[s], t)
			g.size--
		}
	}
}

// Returns a graph with the same vertex and edge set, but with the
// directionality of all its edges reversed.
//
// This implementation returns a new graph object (doubling memory use),
// but not all implementations do so.
func (g *mutableDirected) Transpose() Digraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g2 := &mutableDirected{}
	g2.list = make(map[Vertex]map[Vertex]struct{})

	// Guess at average indegree by looking at ratio of edges to vertices, use that to initially size the adjacency maps
	startcap := int(g.Size() / g.Order())

	for source, adjacent := range g.list {
		if !g2.hasVertex(source) {
			g2.list[source] = make(map[Vertex]struct{}, startcap+1)
		}
		for target, _ := range adjacent {
			if !g2.hasVertex(target) {
				g2.list[target] = make(map[Vertex]struct{}, startcap+1)
			}
			g2.list[target][source] = keyExists
		}
	}

	return g2
}

/* immutableDirected implementation */

type immutableDirected struct {
	al_basic_immut
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *immutableDirected) EachEdge(f EdgeStep) {
	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			if f(NewEdge(source, target)) {
				return
			}
		}
	}
}

// Traverses the set of arcs in the graph, passing each arc to the
// provided closure.
func (g *immutableDirected) EachArc(f ArcStep) {
	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			if f(NewArc(source, target)) {
				return
			}
		}
	}
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *immutableDirected) EachEdgeIncidentTo(v Vertex, f EdgeStep) {
	eachEdgeIncidentToDirected(g, v, f)
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *immutableDirected) EachAdjacentTo(start Vertex, f VertexStep) {
	g.EachEdgeIncidentTo(start, func(e Edge) bool {
		u, v := e.Both()
		if u == start {
			return f(v)
		} else {
			return f(u)
		}
	})
}

// Enumerates the set of out-edges for the provided vertex.
func (g *immutableDirected) EachArcFrom(v Vertex, f ArcStep) {
	if !g.hasVertex(v) {
		return
	}

	for adjacent, _ := range g.list[v] {
		if f(NewArc(v, adjacent)) {
			return
		}
	}
}

func (g *immutableDirected) EachSuccessorOf(v Vertex, f VertexStep) {
	eachVertexInAdjacencyList(g.list, v, f)
}

// Enumerates the set of in-edges for the provided vertex.
func (g *immutableDirected) EachArcTo(v Vertex, f ArcStep) {
	if !g.hasVertex(v) {
		return
	}

	for candidate, adjacent := range g.list {
		for target, _ := range adjacent {
			if target == v {
				if f(NewArc(candidate, target)) {
					return
				}
			}
		}
	}
}

func (g *immutableDirected) EachPredecessorOf(v Vertex, f VertexStep) {
	eachPredecessorOf(g.list, v, f)
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *immutableDirected) Density() float64 {
	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Indicates whether or not the given edge is present in the graph.
func (g *immutableDirected) HasEdge(edge Edge) bool {
	u, v := edge.Both()
	_, exists := g.list[u][v]
	return exists
}

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *immutableDirected) OutDegreeOf(vertex Vertex) (degree int, exists bool) {
	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Returns the indegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
//
// Note that getting indegree is inefficient for directed adjacency lists; it requires
// a full scan of the graph's edge set.
func (g *immutableDirected) InDegreeOf(vertex Vertex) (degree int, exists bool) {
	if exists = g.hasVertex(vertex); exists {
		g.EachArc(func(e Arc) (terminate bool) {
			if vertex == e.Target() {
				degree++
			}
			return
		})
	}

	return
}

// Returns the degree of the vertex, counting both in and out-edges.
func (g *immutableDirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	indegree, exists := g.InDegreeOf(vertex)
	outdegree, exists := g.OutDegreeOf(vertex)
	return indegree + outdegree, exists
}

// Returns a graph with the same vertex and edge set, but with the
// directionality of all its edges reversed.
//
// This implementation returns a new graph object (doubling memory use),
// but not all implementations do so.
func (g *immutableDirected) Transpose() Digraph {
	g2 := &immutableDirected{}
	g2.list = make(map[Vertex]map[Vertex]struct{})

	// Guess at average indegree by looking at ratio of edges to vertices, use that to initially size the adjacency maps
	startcap := int(g.Size() / g.Order())

	for source, adjacent := range g.list {
		if !g2.hasVertex(source) {
			g2.list[source] = make(map[Vertex]struct{}, startcap+1)
		}
		for target, _ := range adjacent {
			if !g2.hasVertex(target) {
				g2.list[target] = make(map[Vertex]struct{}, startcap+1)
			}
			g2.list[target][source] = keyExists
		}
	}

	return g2
}

// Adds a new arc to the graph.
func (g *immutableDirected) addArcs(arcs ...Arc) {
	for _, arc := range arcs {
		g.ensureVertex(arc.Source(), arc.Target())

		if _, exists := g.list[arc.Source()][arc.Target()]; !exists {
			g.list[arc.Source()][arc.Target()] = keyExists
			g.size++
		}
	}
}
