package adjacency_list

import (
	"sync"
	. "github.com/sdboyer/gogl"
)

type al map[Vertex]VertexSet

// Helper to not have to write struct{} everywhere.
var keyExists = struct{}{}

type adjacencyList struct {
	list al
	size int
	mu   sync.RWMutex
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

