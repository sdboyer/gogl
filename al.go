/*
Adjacency lists are a relatively simple graph representation. They maintain
a list of vertices, storing information about edge membership relative to
those vertices. This makes vertex-centric operations generally more
efficient, and edge-centric operations generally less efficient, as edges
are represented implicitly. It also makes them inappropriate for more
complex graph types, such as multigraphs or those with edges having multiple
properties.

gogl's adjacency lists are space-efficient; in a directed graph, the memory
cost for the entire graph G is proportional to V + E; in an undirected graph,
it is V + 2E.

*/
package gogl

import (
	"sync"
)

type al map[Vertex]map[Vertex]struct{}

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

	g.eachAdjacent(vertex, f)
}

// Internal adjacency traverser that bypasses locking.
func (g *adjacencyList) eachAdjacent(vertex Vertex, f func(target Vertex)) {
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

	g.ensureVertex(vertices...)
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *adjacencyList) ensureVertex(vertices ...Vertex) {
	// TODO this is horrible, but the reflection approach in the testing harness requires it...for now
	if g.list == nil {
		g.list = make(map[Vertex]map[Vertex]struct{})
	}

	for _, vertex := range vertices {
		if !g.hasVertex(vertex) {
			// TODO experiment with different lengths...possibly by analyzing existing density?
			g.list[vertex] = make(map[Vertex]struct{}, 10)
		}
	}

	return
}
