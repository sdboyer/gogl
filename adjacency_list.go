package gogl

import "sync"

// VertexSet uses maps to express a value-less (empty struct), indexed unordered list
type VertexSet map[Vertex]struct{}
type adjacencyList map[Vertex]VertexSet

type AdjacencyList struct {
	adjacencyList
	mu sync.RWMutex
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

func (g AdjacencyList) EachEdge(f func(source Vertex, target Vertex)) {
	g.mu.RLock()

	for source, adjacent := range g.adjacencyList {
		for _, target := range adjacent {
			f(source, target)
		}
	}
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

func (g AdjacencyList) AddVertex(vertex Vertex) (success bool) {
	g.mu.Lock()

	if exists := g.HasVertex(vertex); !exists {
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

		var wg sync.WaitGroup
		for _, adjacent := range g.adjacencyList {
			// TODO parallelization here will probably have some cost for sparse graphs
			// but considerable benefit for dense graphs
			wg.Add(1)
			go func(haystack VertexSet) {
				defer wg.Done()

				for adj := range haystack {
					if adj == vertex {
						delete(haystack, adj)
					}
				}
			}(adjacent)

		}
		// Wait for all vertex pruning to complete.
		wg.Wait()
	}

	g.mu.Unlock()
	return
}
