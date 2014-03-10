package gogl

type Directed struct {
	adjacencyList
}

// Composite literal to create a new Directed.
func NewDirected() *Directed {
	list := &Directed{}
	// Cannot assign to promoted fields in a composite literals.
	list.list = make(map[Vertex]VertexSet)

	// Type assertions to ensure interfaces are met
	var _ Graph = list
	var _ SimpleGraph = list
	var _ MutableGraph = list

	return list
}

// Creates a new Directed graph from an edge set.
func NewDirectedFromEdgeSet(set []Edge) *Directed {
	g := NewDirected()

	g.addEdges(set...)

	return g
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

		f := func(v Vertex) {
			if v == vertex {
				degree++
			}
		}

		// This results in a double read-lock. Should be fine.
		for e := range g.list {
			g.EachAdjacent(e, f)
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

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *Directed) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
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
