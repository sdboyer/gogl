package gogl

type Directed struct {
	al_basic_mut
}

// Creates a new mutable, directed graph.
func NewDirected() *Directed {
	list := &Directed{}
	// Cannot assign to promoted fields in a composite literals.
	list.list = make(map[Vertex]map[Vertex]struct{})

	// Type assertions to ensure interfaces are met
	var _ Graph = list
	var _ SimpleGraph = list
	var _ MutableGraph = list
	var _ DirectedGraph = list

	return list
}

/* Directed additions */

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *Directed) OutDegree(vertex Vertex) (degree int, exists bool) {
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
func (g *Directed) InDegree(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		// This results in a double read-lock. Should be fine.
		for e := range g.list {
			g.EachAdjacent(e, func(v Vertex) {
				if v == vertex {
					degree++
				}
			})
		}
	}

	return
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *Directed) EachEdge(f func(edge Edge)) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			f(BaseEdge{U: source, V: target})
		}
	}
}

// Indicates whether or not the given edge is present in the graph.
func (g *Directed) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.list[edge.Source()][edge.Target()]
	return exists
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *Directed) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *Directed) RemoveVertex(vertices ...Vertex) {
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

// Adds edges to the graph.
func (g *Directed) AddEdges(edges ...Edge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *Directed) addEdges(edges ...Edge) {
	for _, edge := range edges {
		g.ensureVertex(edge.Source(), edge.Target())

		if _, exists := g.list[edge.Source()][edge.Target()]; !exists {
			g.list[edge.Source()][edge.Target()] = keyExists
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *Directed) RemoveEdges(edges ...Edge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, edge := range edges {
		s, t := edge.Both()
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
func (g *Directed) Transpose() DirectedGraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g2 := &Directed{}
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

/* ImmutableDirected implementation */

type immutableDirected struct {
	al_basic_immut
}

// Creates a new mutable, directed graph.
func NewImmutableDirected(g DirectedGraph) DirectedGraph {
	list := &immutableDirected{}
	// Cannot assign to promoted fields in a composite literals.
	list.list = make(map[Vertex]map[Vertex]struct{})

	g.EachEdge(func(edge Edge) {
		list.ensureVertex(edge.Source(), edge.Target())

		if _, exists := list.list[edge.Source()][edge.Target()]; !exists {
			list.list[edge.Source()][edge.Target()] = keyExists
			list.size++
		}
	})

	if list.Order() != g.Order() {
		g.EachVertex(func(vertex Vertex) {
			list.ensureVertex(vertex)
		})
	}

	// Type assertions to ensure interfaces are met
	var _ Graph = list
	var _ SimpleGraph = list
	var _ DirectedGraph = list

	return list
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *immutableDirected) EachEdge(f func(edge Edge)) {
	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			f(BaseEdge{U: source, V: target})
		}
	}
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *immutableDirected) Density() float64 {
	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Indicates whether or not the given edge is present in the graph.
func (g *immutableDirected) HasEdge(edge Edge) bool {
	_, exists := g.list[edge.Source()][edge.Target()]
	return exists
}

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *immutableDirected) OutDegree(vertex Vertex) (degree int, exists bool) {
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
func (g *immutableDirected) InDegree(vertex Vertex) (degree int, exists bool) {
	if exists = g.hasVertex(vertex); exists {
		// This results in a double read-lock. Should be fine.
		for e := range g.list {
			g.EachAdjacent(e, func(v Vertex) {
				if v == vertex {
					degree++
				}
			})
		}
	}

	return
}

// Returns a graph with the same vertex and edge set, but with the
// directionality of all its edges reversed.
//
// This implementation returns a new graph object (doubling memory use),
// but not all implementations do so.
func (g *immutableDirected) Transpose() DirectedGraph {
	g2 := &Directed{}
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
