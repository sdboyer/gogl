package al

import (
	. "github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

type mutableUndirected struct {
	al_basic_mut
}

// Returns the degree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *mutableUndirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *mutableUndirected) Edges(f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := set.New(set.NonThreadSafe)

	for source, adjacent := range g.list {
		for target := range adjacent {
			e := NewEdge(source, target)
			if !visited.Has(NewEdge(target, source)) {
				visited.Add(e)
				if f(e) {
					return
				}
			}
		}
	}
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *mutableUndirected) IncidentTo(v Vertex, f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent := range g.list[v] {
		if f(NewEdge(v, adjacent)) {
			return
		}
	}
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *mutableUndirected) AdjacentTo(vertex Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachVertexInAdjacencyList(g.list, vertex, f)
}

// Indicates whether or not the given edge is present in the graph.
func (g *mutableUndirected) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Spread it into two expressions to avoid evaluating the second if possible
	u, v := edge.Both()
	if _, exists := g.list[u][v]; exists {
		return true
	} else if _, exists := g.list[v][u]; exists {
		return true
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *mutableUndirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *mutableUndirected) RemoveVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, vertex := range vertices {
		if g.hasVertex(vertex) {
			eachVertexInAdjacencyList(g.list, vertex, func(adjacent Vertex) (terminate bool) {
				delete(g.list[adjacent], vertex)
				return
			})
			g.size -= len(g.list[vertex])
			delete(g.list, vertex)
		}
	}
	return
}

// Adds edges to the graph.
func (g *mutableUndirected) AddEdges(edges ...Edge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *mutableUndirected) addEdges(edges ...Edge) {
	for _, edge := range edges {
		u, v := edge.Both()
		g.ensureVertex(u, v)

		if _, exists := g.list[u][v]; !exists {
			g.list[u][v] = keyExists
			g.list[v][u] = keyExists
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *mutableUndirected) RemoveEdges(edges ...Edge) {
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
