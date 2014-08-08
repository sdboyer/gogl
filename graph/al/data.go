package al

import (
	"sync"

	. "github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

// This is implemented as an adjacency list, because those are simple.
type baseData struct {
	list map[Vertex]map[Vertex]interface{}
	size int
	mu   sync.RWMutex
}

/* baseData shared methods */

// Traverses the graph's vertices in random order, passing each vertex to the
// provided closure.
func (g *baseData) EachVertex(f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for v := range g.list {
		if f(v) {
			return
		}
	}
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseData) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	exists = g.hasVertex(vertex)
	return
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseData) hasVertex(vertex Vertex) (exists bool) {
	_, exists = g.list[vertex]
	return
}

// Returns the order (number of vertices) in the graph.
func (g *baseData) Order() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.list)
}

// Returns the size (number of edges) in the graph.
func (g *baseData) Size() int {
	return g.size
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseData) EnsureVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.ensureVertex(vertices...)
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseData) ensureVertex(vertices ...Vertex) {
	for _, vertex := range vertices {
		if !g.hasVertex(vertex) {
			// TODO experiment with different lengths...possibly by analyzing existing density?
			g.list[vertex] = make(map[Vertex]interface{}, 10)
		}
	}

	return
}

/* DirectedData implementation */

type dataDirected struct {
	baseData
}

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *dataDirected) OutDegreeOf(vertex Vertex) (degree int, exists bool) {
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
func (g *dataDirected) InDegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return inDegreeOf(g, vertex)
}

// Returns the degree of the provided vertex, counting both in and out-edges.
func (g *dataDirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	indegree, exists := inDegreeOf(g, vertex)
	outdegree, exists := g.OutDegreeOf(vertex)
	return indegree + outdegree, exists
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *dataDirected) EachEdgeIncidentTo(v Vertex, f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	eachEdgeIncidentToDirected(g, v, f)
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *dataDirected) EachAdjacentTo(start Vertex, f VertexStep) {
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
func (g *dataDirected) EachArcFrom(v Vertex, f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, data := range g.list[v] {
		if f(NewDataArc(v, adjacent, data)) {
			return
		}
	}
}

func (g *dataDirected) EachSuccessorOf(v Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachVertexInAdjacencyList(g.list, v, f)
}

// Enumerates the set of in-edges for the provided vertex.
func (g *dataDirected) EachArcTo(v Vertex, f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for candidate, adjacent := range g.list {
		for target, data := range adjacent {
			if target == v {
				if f(NewDataArc(candidate, target, data)) {
					return
				}
			}
		}
	}
}

func (g *dataDirected) EachPredecessorOf(v Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachPredecessorOf(g.list, v, f)
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *dataDirected) EachEdge(f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, data := range adjacent {
			if f(NewDataEdge(source, target, data)) {
				return
			}
		}
	}
}

// Traverses the set of arcs in the graph, passing each arc to the
// provided closure.
func (g *dataDirected) EachArc(f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, data := range adjacent {
			if f(NewDataArc(source, target, data)) {
				return
			}
		}
	}
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge property.
func (g *dataDirected) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	u, v := edge.Both()
	_, exists := g.list[u][v]
	if !exists {
		_, exists = g.list[v][u]
	}
	return exists
}

// Indicates whether or not the given arc is present in the graph.
func (g *dataDirected) HasArc(arc Arc) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.list[arc.Source()][arc.Target()]
	return exists
}

// Indicates whether or not the given property edge is present in the graph.
// It will only match if the provided DataEdge has the same property as
// the edge contained in the graph.
func (g *dataDirected) HasDataEdge(edge DataEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	u, v := edge.Both()
	if data, exists := g.list[u][v]; exists {
		return data == edge.Data()
	} else if data, exists = g.list[v][u]; exists {
		return data == edge.Data()
	}
	return false
}

// Indicates whether or not the given data arc is present in the graph.
// It will only match if the provided DataEdge has the same data as
// the edge contained in the graph.
func (g *dataDirected) HasDataArc(arc DataArc) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if data, exists := g.list[arc.Source()][arc.Target()]; exists {
		return data == arc.Data()
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *dataDirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *dataDirected) RemoveVertex(vertices ...Vertex) {
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

// Adds arcs to the graph.
func (g *dataDirected) AddArcs(arcs ...DataArc) {
	if len(arcs) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addArcs(arcs...)
}

// Adds a new arc to the graph.
func (g *dataDirected) addArcs(arcs ...DataArc) {
	for _, arc := range arcs {
		u, v := arc.Both()
		g.ensureVertex(u, v)

		if _, exists := g.list[u][v]; !exists {
			g.list[u][v] = arc.Data()
			g.size++
		}
	}
}

// Removes arcs from the graph. This does NOT remove vertex members of the
// removed arcs.
func (g *dataDirected) RemoveArcs(arcs ...DataArc) {
	if len(arcs) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, arc := range arcs {
		s, t := arc.Both()
		if _, exists := g.list[s][t]; exists {
			delete(g.list[s], t)
			g.size--
		}
	}
}

func (g *dataDirected) Transpose() Digraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g2 := &dataDirected{}
	g2.list = make(map[Vertex]map[Vertex]interface{})

	// Guess at average indegree by looking at ratio of edges to vertices, use that to initially size the adjacency maps
	startcap := int(g.Size() / g.Order())

	for source, adjacent := range g.list {
		if !g2.hasVertex(source) {
			g2.list[source] = make(map[Vertex]interface{}, startcap+1)
		}

		for target, data := range adjacent {
			if !g2.hasVertex(target) {
				g2.list[target] = make(map[Vertex]interface{}, startcap+1)
			}
			g2.list[target][source] = data
		}
	}

	return g2
}

/* UndirectedData implementation */

type dataUndirected struct {
	baseData
}

// Returns the degree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *dataUndirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *dataUndirected) EachEdge(f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := set.NewNonTS()

	var e DataEdge
	for source, adjacent := range g.list {
		for target, data := range adjacent {
			e = NewDataEdge(source, target, data)
			if !visited.Has(NewEdge(e.Both())) {
				visited.Add(NewEdge(target, source))
				if f(e) {
					return
				}
			}
		}
	}
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *dataUndirected) EachEdgeIncidentTo(v Vertex, f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, data := range g.list[v] {
		if f(NewDataEdge(v, adjacent, data)) {
			return
		}
	}
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *dataUndirected) EachAdjacentTo(vertex Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachVertexInAdjacencyList(g.list, vertex, f)
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge property.
func (g *dataUndirected) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Spread it into two expressions to avoid evaluating the second if possible
	u, v := edge.Both()
	_, exists := g.list[u][v]
	return exists
}

// Indicates whether or not the given property edge is present in the graph.
// It will only match if the provided DataEdge has the same property as
// the edge contained in the graph.
func (g *dataUndirected) HasDataEdge(edge DataEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Spread it into two expressions to avoid evaluating the second if possible
	u, v := edge.Both()
	if data, exists := g.list[u][v]; exists {
		return edge.Data() == data
	} else if data, exists := g.list[v][u]; exists {
		return edge.Data() == data
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *dataUndirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *dataUndirected) RemoveVertex(vertices ...Vertex) {
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
func (g *dataUndirected) AddEdges(edges ...DataEdge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *dataUndirected) addEdges(edges ...DataEdge) {
	for _, edge := range edges {
		u, v := edge.Both()
		g.ensureVertex(u, v)

		if _, exists := g.list[u][v]; !exists {
			d := edge.Data()
			g.list[u][v] = d
			g.list[v][u] = d
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *dataUndirected) RemoveEdges(edges ...DataEdge) {
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
