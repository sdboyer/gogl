package gogl

import (
	"gopkg.in/fatih/set.v0"
)

type Undirected struct {
	adjacencyList
}

func NewUndirected() *Undirected {
	list := &Undirected{}
	// Cannot assign to promoted fields in a composite literals.
	list.list = make(map[Vertex]map[Vertex]struct{})

	// Type assertions to ensure interfaces are met
	var _ Graph = list
	var _ SimpleGraph = list
	var _ MutableGraph = list

	return list
}

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *Undirected) OutDegree(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Returns the indegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *Undirected) InDegree(vertex Vertex) (degree int, exists bool) {
	return g.OutDegree(vertex)
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *Undirected) EachEdge(f func(edge Edge)) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := set.NewNonTS()

	for source, adjacent := range g.list {
		for target, _ := range adjacent {
			e := BaseEdge{U: source, V: target}
			if !visited.Has(BaseEdge{U: target, V: source}) {
				visited.Add(e)
				f(e)
			}
		}
	}
}

// Indicates whether or not the given edge is present in the graph.
func (g *Undirected) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Spread it into two expressions to avoid evaluating the second if possible
	if _, exists := g.list[edge.Source()][edge.Target()]; exists {
		return true
	} else if _, exists := g.list[edge.Target()][edge.Source()]; exists {
		return true
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *Undirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *Undirected) RemoveVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, vertex := range vertices {
		if g.hasVertex(vertex) {
			g.eachAdjacent(vertex, func(adjacent Vertex) {
				delete(g.list[adjacent], vertex)
			})
			g.size -= len(g.list[vertex])
			delete(g.list, vertex)
		}
	}
	return
}

// Adds edges to the graph.
func (g *Undirected) AddEdges(edges ...Edge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *Undirected) addEdges(edges ...Edge) {
	for _, edge := range edges {
		g.ensureVertex(edge.Source(), edge.Target())

		if _, exists := g.list[edge.Source()][edge.Target()]; !exists {
			g.list[edge.Source()][edge.Target()] = keyExists
			g.list[edge.Target()][edge.Source()] = keyExists
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *Undirected) RemoveEdges(edges ...Edge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, edge := range edges {
		s, t := edge.Both()
		if _, exists := g.list[s][t]; exists {
			delete(g.list[s], t)
			delete(g.list[t], s)
			g.size--
		}
	}
}
