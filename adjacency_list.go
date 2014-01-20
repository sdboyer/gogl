package gogl

import "sync"

type al map[Vertex]VertexSet

// Helper to not have to write struct{} everywhere.
var keyExists = struct{}{}

type adjacencyList struct {
	list al
	size uint
	mu   sync.RWMutex
}

type DirectedAdjacencyList struct {
	adjacencyList
}

// Composite literal to create a new DirectedAdjacencyList.
func NewDirectedAdjacencyList() *DirectedAdjacencyList {
	// Cannot assign to promoted fields in a composite literals.
	list := &DirectedAdjacencyList{}
	list.list = make(map[Vertex]VertexSet)
	return list
}

func NewDirectedAdjacencyListFromEdgeSet(set []Edge) *DirectedAdjacencyList {
	g := NewDirectedAdjacencyList()

	for _, edge := range set {
		g.addEdge(edge)
	}

	return g
}

/* Base adjacencyList methods */

func (g *adjacencyList) EachVertex(f func(vertex Vertex)) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for v := range g.list {
		f(v)
	}
}

func (g *adjacencyList) EachAdjacent(vertex Vertex, f func(target Vertex)) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if _, exists := g.list[vertex]; exists {
		for adjacent, _ := range g.list[vertex] {
			f(adjacent)
		}
	}
}

func (g *adjacencyList) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	exists = g.hasVertex(vertex)
	return
}

func (g *adjacencyList) hasVertex(vertex Vertex) (exists bool) {
	_, exists = g.list[vertex]
	return
}

func (g *adjacencyList) Order() uint {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return uint(len(g.list))
}

func (g *adjacencyList) Size() uint {
	return g.size
}

func (g *adjacencyList) EnsureVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, vertex := range vertices {
		g.addVertex(vertex)
	}
}

func (g *adjacencyList) addVertex(vertex Vertex) (success bool) {
	if exists := g.hasVertex(vertex); !exists {
		// TODO experiment with different lengths...possibly by analyzing existing density?
		g.list[vertex] = make(VertexSet, 10)
		success = true
	}

	return
}

/* DirectedAdjacencyList additions (TODO - call it Directed) */

func (g *DirectedAdjacencyList) OutDegree(vertex Vertex) (degree uint, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = uint(len(g.list[vertex]))
	}
	return
}

// Getting InDegree is inefficient for directed adjacency lists
func (g *DirectedAdjacencyList) InDegree(vertex Vertex) (degree uint, exists bool) {
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

func (g *DirectedAdjacencyList) EachEdge(f func(edge Edge)) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			f(&BaseEdge{u: source, v: target})
		}
	}
}

func (g *DirectedAdjacencyList) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
}

func (g *DirectedAdjacencyList) RemoveVertex(vertices ...Vertex) {
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

func (g *DirectedAdjacencyList) AddEdge(edge Edge) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.addEdge(edge)
}

func (g *DirectedAdjacencyList) addEdge(edge Edge) (exists bool) {
	g.addVertex(edge.Source())
	g.addVertex(edge.Target())

	if _, exists = g.list[edge.Source()][edge.Target()]; !exists {
		g.list[edge.Source()][edge.Target()] = keyExists
		g.size++
	}
	return !exists
}

func (g *DirectedAdjacencyList) RemoveEdge(edge Edge) {
	g.mu.Lock()
	defer g.mu.Unlock()

	s, t := edge.Both()
	if _, exists := g.list[s][t]; exists {
		delete(g.list[s], t)
		g.size--
	}
}
