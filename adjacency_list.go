package gogl

import "sync"

type al map[Vertex]VertexSet

// Helper to not have to write struct{} everywhere.
var keyExists = struct{}{}

type adjacencyList struct {
	list al
	size int
	mu   sync.RWMutex
}

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

func NewDirectedFromEdgeSet(set []Edge) *Directed {
	g := NewDirected()

	for _, edge := range set {
		g.addEdge(edge)
	}

	return g
}

/* Base adjacencyList methods */
// Traverses the graph's vertices in random order, passing each vertex to the
// provided closure.
func (g *adjacencyList) EachVertex(f func(vertex Vertex)) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for v := range g.list {
		f(v)
	}
}

// Given a vertex present in the graph, passes each vertex adjacent to the
// provided vertex to the provided closure.
func (g *adjacencyList) EachAdjacent(vertex Vertex, f func(target Vertex)) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.list[vertex]; exists {
		for adjacent, _ := range g.list[vertex] {
			f(adjacent)
		}
	}
}

// Indicates whether or not the given vertex is present in the graph.
func (g *adjacencyList) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	exists = g.hasVertex(vertex)
	return
}

// Indicates whether or not the given vertex is present in the graph.
func (g *adjacencyList) hasVertex(vertex Vertex) (exists bool) {
	_, exists = g.list[vertex]
	return
}

// Returns the order (number of vertices) in the graph.
func (g *adjacencyList) Order() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.list)
}

// Returns the size (number of edges) in the graph.
func (g *adjacencyList) Size() int {
	return g.size
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *adjacencyList) EnsureVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, vertex := range vertices {
		g.ensureVertex(vertex)
	}
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *adjacencyList) ensureVertex(vertex Vertex) (success bool) {
	if exists := g.hasVertex(vertex); !exists {
		// TODO experiment with different lengths...possibly by analyzing existing density?
		g.list[vertex] = make(VertexSet, 10)
		success = true
	}

	return
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
			f(&BaseEdge{u: source, v: target})
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

// Adds a new edge to the graph.
func (g *Directed) AddEdge(edge Edge) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.addEdge(edge)
}

// Adds a new edge to the graph.
func (g *Directed) addEdge(edge Edge) (exists bool) {
	g.ensureVertex(edge.Source())
	g.ensureVertex(edge.Target())

	if _, exists = g.list[edge.Source()][edge.Target()]; !exists {
		g.list[edge.Source()][edge.Target()] = keyExists
		g.size++
	}
	return !exists
}

// Removes an edge from the graph. This does NOT remove vertex members of the
// removed edge.
func (g *Directed) RemoveEdge(edge Edge) {
	g.mu.Lock()
	defer g.mu.Unlock()

	s, t := edge.Both()
	if _, exists := g.list[s][t]; exists {
		delete(g.list[s], t)
		g.size--
	}
}

