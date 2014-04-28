package gogl

import (
	"sync"

	"gopkg.in/fatih/set.v0"
)

// This is implemented as an adjacency list, because those are simple.
type baseProperty struct {
	list map[Vertex]map[Vertex]interface{}
	size int
	mu   sync.RWMutex
}

/* baseProperty shared methods */

// Traverses the graph's vertices in random order, passing each vertex to the
// provided closure.
func (g *baseProperty) EachVertex(f VertexLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for v := range g.list {
		if f(v) {
			return
		}
	}
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseProperty) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	exists = g.hasVertex(vertex)
	return
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseProperty) hasVertex(vertex Vertex) (exists bool) {
	_, exists = g.list[vertex]
	return
}

// Returns the order (number of vertices) in the graph.
func (g *baseProperty) Order() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.list)
}

// Returns the size (number of edges) in the graph.
func (g *baseProperty) Size() int {
	return g.size
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseProperty) EnsureVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.ensureVertex(vertices...)
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseProperty) ensureVertex(vertices ...Vertex) {
	for _, vertex := range vertices {
		if !g.hasVertex(vertex) {
			// TODO experiment with different lengths...possibly by analyzing existing density?
			g.list[vertex] = make(map[Vertex]interface{}, 10)
		}
	}

	return
}

/* DirectedProperty implementation */

type propertyDirected struct {
	baseProperty
}

func NewPropertyDirected() MutablePropertyGraph {
	list := &propertyDirected{}
	// Cannot assign to promoted fields in a composite literals.
	list.list = make(map[Vertex]map[Vertex]interface{})

	// Type assertions to ensure interfaces are met
	var _ Graph = list
	var _ SimpleGraph = list
	var _ DirectedGraph = list
	var _ PropertyGraph = list
	var _ MutablePropertyGraph = list

	return list
}

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *propertyDirected) OutDegreeOf(vertex Vertex) (degree int, exists bool) {
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
func (g *propertyDirected) InDegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return inDegreeOf(g, vertex)
}

// Returns the degree of the provided vertex, counting both in and out-edges.
func (g *propertyDirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	indegree, exists := inDegreeOf(g, vertex)
	outdegree, exists := g.OutDegreeOf(vertex)
	return indegree + outdegree, exists
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *propertyDirected) EachEdgeIncidentTo(v Vertex, f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var terminate bool
	interloper := func(e Edge) bool {
		terminate = terminate || f(e)
		return terminate
	}

	g.EachArcFrom(v, interloper)
	g.EachArcTo(v, interloper)
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *propertyDirected) EachAdjacentTo(start Vertex, f VertexLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.EachEdgeIncidentTo(start, func(e Edge) bool {
		u, v := e.Both()
		if u == start {
			return f(v)
		} else {
			return f(u)
		}
	})
}

// Enumerates the set of out-edges for the provided vertex.
func (g *propertyDirected) EachArcFrom(v Vertex, f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, property := range g.list[v] {
		if f(BasePropertyEdge{BaseEdge{U: v, V: adjacent}, property}) {
			return
		}
	}
}

// Enumerates the set of in-edges for the provided vertex.
func (g *propertyDirected) EachArcTo(v Vertex, f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for candidate, adjacent := range g.list {
		for target, property := range adjacent {
			if target == v {
				if f(BasePropertyEdge{BaseEdge{U: candidate, V: target}, property}) {
					return
				}
			}
		}
	}
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *propertyDirected) EachEdge(f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, property := range adjacent {
			if f(BasePropertyEdge{BaseEdge{U: source, V: target}, property}) {
				return
			}
		}
	}
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge property.
func (g *propertyDirected) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.list[edge.Source()][edge.Target()]
	return exists
}

// Indicates whether or not the given property edge is present in the graph.
// It will only match if the provided PropertyEdge has the same property as
// the edge contained in the graph.
func (g *propertyDirected) HasPropertyEdge(edge PropertyEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if property, exists := g.list[edge.Source()][edge.Target()]; exists {
		return property == edge.Property()
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *propertyDirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *propertyDirected) RemoveVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, vertex := range vertices {
		if g.hasVertex(vertex) {
			g.size -= len(g.list[vertex])
			delete(g.list, vertex)

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
func (g *propertyDirected) AddEdges(edges ...PropertyEdge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *propertyDirected) addEdges(edges ...PropertyEdge) {
	for _, edge := range edges {
		g.ensureVertex(edge.Source(), edge.Target())

		if _, exists := g.list[edge.Source()][edge.Target()]; !exists {
			g.list[edge.Source()][edge.Target()] = edge.Property()
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *propertyDirected) RemoveEdges(edges ...PropertyEdge) {
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

func (g *propertyDirected) Transpose() DirectedGraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g2 := &propertyDirected{}
	g2.list = make(map[Vertex]map[Vertex]interface{})

	// Guess at average indegree by looking at ratio of edges to vertices, use that to initially size the adjacency maps
	startcap := int(g.Size() / g.Order())

	for source, adjacent := range g.list {
		if !g2.hasVertex(source) {
			g2.list[source] = make(map[Vertex]interface{}, startcap+1)
		}

		for target, property := range adjacent {
			if !g2.hasVertex(target) {
				g2.list[target] = make(map[Vertex]interface{}, startcap+1)
			}
			g2.list[target][source] = property
		}
	}

	return g2
}

/* UndirectedProperty implementation */

type propertyUndirected struct {
	baseProperty
}

func NewPropertyUndirected() MutablePropertyGraph {
	g := &propertyUndirected{}
	// Cannot assign to promoted fields in a composite literals.
	g.list = make(map[Vertex]map[Vertex]interface{})

	// Type assertions to ensure interfaces are met
	var _ Graph = g
	var _ SimpleGraph = g
	var _ PropertyGraph = g
	var _ MutablePropertyGraph = g

	return g
}

// Returns the degree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *propertyUndirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *propertyUndirected) EachEdge(f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := set.NewNonTS()

	for source, adjacent := range g.list {
		for target, property := range adjacent {
			be := BaseEdge{U: source, V: target}
			e := BasePropertyEdge{be, property}
			if !visited.Has(BaseEdge{U: target, V: source}) {
				visited.Add(be)
				if f(e) {
					return
				}
			}
		}
	}
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *propertyUndirected) EachEdgeIncidentTo(v Vertex, f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, property := range g.list[v] {
		if f(BasePropertyEdge{BaseEdge{U: v, V: adjacent}, property}) {
			return
		}
	}
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *propertyUndirected) EachAdjacentTo(vertex Vertex, f VertexLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachAdjacentToUndirected(g.list, vertex, f)
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge property.
func (g *propertyUndirected) HasEdge(edge Edge) bool {
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

// Indicates whether or not the given property edge is present in the graph.
// It will only match if the provided PropertyEdge has the same property as
// the edge contained in the graph.
func (g *propertyUndirected) HasPropertyEdge(edge PropertyEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Spread it into two expressions to avoid evaluating the second if possible
	if property, exists := g.list[edge.Source()][edge.Target()]; exists {
		return edge.Property() == property
	} else if property, exists := g.list[edge.Target()][edge.Source()]; exists {
		return edge.Property() == property
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *propertyUndirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *propertyUndirected) RemoveVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, vertex := range vertices {
		if g.hasVertex(vertex) {
			eachAdjacentToUndirected(g.list, vertex, func(adjacent Vertex) (terminate bool) {
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
func (g *propertyUndirected) AddEdges(edges ...PropertyEdge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *propertyUndirected) addEdges(edges ...PropertyEdge) {
	for _, edge := range edges {
		g.ensureVertex(edge.Source(), edge.Target())

		if _, exists := g.list[edge.Source()][edge.Target()]; !exists {
			d := edge.Property()
			g.list[edge.Source()][edge.Target()] = d
			g.list[edge.Target()][edge.Source()] = d
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *propertyUndirected) RemoveEdges(edges ...PropertyEdge) {
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
