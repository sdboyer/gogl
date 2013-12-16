package gogl

import "sync"

// VertexSet uses maps to express a value-less (empty struct), indexed
// unordered list. See
// https://groups.google.com/forum/#!searchin/golang-nuts/map/golang-nuts/H2cXpwisEUE/1X2FV-rODfIJ
type VertexSet map[Vertex]struct{}
type adjacencyList map[Vertex]VertexSet

// Helper to not have to write struct{} everywhere.
var keyExists = struct{}{}

type AdjacencyList struct {
	adjacencyList
	size uint
	mu   sync.RWMutex
}

// Composite literal to create a new AdjacencyList.
func NewAdjacencyList() *AdjacencyList {
	return &AdjacencyList{
		adjacencyList: make(map[Vertex]VertexSet)}
}

func (g AdjacencyList) EachVertex(f func(vertex Vertex)) {
	g.mu.RLock()

	for v := range g.adjacencyList {
		f(v)
	}

	g.mu.RUnlock()
}

func (g AdjacencyList) EachEdge(f func(edge Edge)) {
	g.mu.RLock()

	for source, adjacent := range g.adjacencyList {
		for _, target := range adjacent {
			f(BaseEdge{u: source, v: target})
		}
	}

	g.mu.RUnlock()
}

func (g AdjacencyList) EachAdjacent(vertex Vertex, f func(target Vertex)) {
	g.mu.RLock()

	if _, exists := g.adjacencyList[vertex]; exists {
		for _, adjacent := range g.adjacencyList[vertex] {
			f(adjacent)
		}
	}

	g.mu.RUnlock()
}

func (g AdjacencyList) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()

	_, exists = g.adjacencyList[vertex]

	g.mu.RUnlock()
	return
}

func (g AdjacencyList) Order() (length uint) {
	g.mu.RLock()

	length = uint(len(g.adjacencyList))

	g.mu.RUnlock()
	return
}

func (g AdjacencyList) Size() uint {
	return g.size
}

func (g AdjacencyList) Density() (density float64) {
	g.mu.RLock()

	order := float64(g.Order())
	density = (2 * float64(g.Size())) / (order * (order - 1))

	g.mu.RUnlock()
	return
}

func (g AdjacencyList) AddVertex(vertex Vertex) (success bool) {
	g.mu.Lock()

	if _, exists := g.adjacencyList[vertex]; !exists {
		// TODO experiment with different lengths...possibly by analyzing existing density?
		g.adjacencyList[vertex] = make(VertexSet, 10)
		success = true
	}

	g.mu.Unlock()
	return
}

func (g AdjacencyList) RemoveVertex(vertex Vertex) (success bool) {
	g.mu.Lock()
	if g.HasVertex(vertex) {
		// TODO Is the expensive search good to do here and now...
		// while read-locked?
		delete(g.adjacencyList, vertex)

		// TODO consider chunking the list and parallelizing into goroutines
		for _, adjacent := range g.adjacencyList {
			if _, has := adjacent[vertex]; has {
				delete(adjacent, vertex)
				g.size--
			}
		}
	}

	g.mu.Unlock()
	return
}

func (g AdjacencyList) AddEdge(edge Edge) (exists bool) {
	g.mu.Lock()

	g.AddVertex(edge.Source())
	g.AddVertex(edge.Target())

	if _, exists = g.adjacencyList[edge.Source()][edge.Target]; !exists {
		g.adjacencyList[edge.Source()][edge.Target()] = keyExists
	}

	g.mu.Unlock()
	return !exists
}

func (g AdjacencyList) RemoveEdge(edge Edge) {
	g.mu.Lock()

	delete(g.adjacencyList, edge.Source())

	g.mu.Unlock()
}
