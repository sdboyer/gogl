package al

import (
	"sync"

	. "github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

// This is implemented as an adjacency list, because those are simple.
type baseWeighted struct {
	list map[Vertex]map[Vertex]float64
	size int
	mu   sync.RWMutex
}

/* baseWeighted shared methods */

// Traverses the graph's vertices in random order, passing each vertex to the
// provided closure.
func (g *baseWeighted) EachVertex(f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for v := range g.list {
		if f(v) {
			return
		}
	}
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseWeighted) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	exists = g.hasVertex(vertex)
	return
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseWeighted) hasVertex(vertex Vertex) (exists bool) {
	_, exists = g.list[vertex]
	return
}

// Returns the order (number of vertices) in the graph.
func (g *baseWeighted) Order() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.list)
}

// Returns the size (number of edges) in the graph.
func (g *baseWeighted) Size() int {
	return g.size
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseWeighted) EnsureVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.ensureVertex(vertices...)
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseWeighted) ensureVertex(vertices ...Vertex) {
	for _, vertex := range vertices {
		if !g.hasVertex(vertex) {
			// TODO experiment with different lengths...possibly by analyzing existing density?
			g.list[vertex] = make(map[Vertex]float64, 10)
		}
	}

	return
}

/* DirectedWeighted implementation */

type weightedDirected struct {
	baseWeighted
}

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *weightedDirected) OutDegreeOf(vertex Vertex) (degree int, exists bool) {
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
func (g *weightedDirected) InDegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return inDegreeOf(g, vertex)
}

// Returns the degree of the given vertex, counting both in and out-edges.
func (g *weightedDirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	indegree, exists := inDegreeOf(g, vertex)
	outdegree, exists := g.OutDegreeOf(vertex)
	return indegree + outdegree, exists
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *weightedDirected) EachEdgeIncidentTo(v Vertex, f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	eachEdgeIncidentToDirected(g, v, f)
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *weightedDirected) EachAdjacentTo(start Vertex, f VertexStep) {
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
func (g *weightedDirected) EachArcFrom(v Vertex, f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, weight := range g.list[v] {
		if f(NewWeightedArc(v, adjacent, weight)) {
			return
		}
	}
}

func (g *weightedDirected) EachSuccessorOf(v Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachVertexInAdjacencyList(g.list, v, f)
}

// Enumerates the set of in-edges for the provided vertex.
func (g *weightedDirected) EachArcTo(v Vertex, f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for candidate, adjacent := range g.list {
		for target, weight := range adjacent {
			if target == v {
				if f(NewWeightedArc(candidate, target, weight)) {
					return
				}
			}
		}
	}
}

func (g *weightedDirected) EachPredecessorOf(v Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachPredecessorOf(g.list, v, f)
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *weightedDirected) EachEdge(f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, weight := range adjacent {
			if f(NewWeightedEdge(source, target, weight)) {
				return
			}
		}
	}
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge weight.
func (g *weightedDirected) HasEdge(edge Edge) bool {
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
func (g *weightedDirected) HasArc(arc Arc) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.list[arc.Source()][arc.Target()]
	return exists
}

// Indicates whether or not the given weighted edge is present in the graph.
// It will only match if the provided WeightedEdge has the same weight as
// the edge contained in the graph.
func (g *weightedDirected) HasWeightedEdge(edge WeightedEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	u, v := edge.Both()
	if weight, exists := g.list[u][v]; exists {
		return weight == edge.Weight()
	} else if weight, exists = g.list[v][u]; exists {
		return weight == edge.Weight()
	}
	return false
}

// Indicates whether or not the given weighted arc is present in the graph.
// It will only match if the provided LabeledEdge has the same label as
// the edge contained in the graph.
func (g *weightedDirected) HasWeightedArc(arc WeightedArc) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if weight, exists := g.list[arc.Source()][arc.Target()]; exists {
		return weight == arc.Weight()
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *weightedDirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *weightedDirected) RemoveVertex(vertices ...Vertex) {
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
func (g *weightedDirected) AddArcs(arcs ...WeightedArc) {
	if len(arcs) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addArcs(arcs...)
}

// Adds a new arc to the graph.
func (g *weightedDirected) addArcs(arcs ...WeightedArc) {
	for _, arc := range arcs {
		g.ensureVertex(arc.Source(), arc.Target())

		if _, exists := g.list[arc.Source()][arc.Target()]; !exists {
			g.list[arc.Source()][arc.Target()] = arc.Weight()
			g.size++
		}
	}
}

// Removes arcs from the graph. This does NOT remove vertex members of the
// removed arcs.
func (g *weightedDirected) RemoveArcs(arcs ...WeightedArc) {
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

func (g *weightedDirected) Transpose() Digraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g2 := &weightedDirected{}
	g2.list = make(map[Vertex]map[Vertex]float64)

	// Guess at average indegree by looking at ratio of edges to vertices, use that to initially size the adjacency maps
	startcap := int(g.Size() / g.Order())

	for source, adjacent := range g.list {
		if !g2.hasVertex(source) {
			g2.list[source] = make(map[Vertex]float64, startcap+1)
		}

		for target, weight := range adjacent {
			if !g2.hasVertex(target) {
				g2.list[target] = make(map[Vertex]float64, startcap+1)
			}
			g2.list[target][source] = weight
		}
	}

	return g2
}

/* UndirectedWeighted implementation */

type weightedUndirected struct {
	baseWeighted
}

// Returns the degree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *weightedUndirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *weightedUndirected) EachEdge(f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := set.NewNonTS()

	var e WeightedEdge
	for source, adjacent := range g.list {
		for target, weight := range adjacent {
			e = NewWeightedEdge(source, target, weight)
			if !visited.Has(NewEdge(e.Both())) {
				visited.Add(NewEdge(target, source))
				if f(e) {
					return
				}
			}
		}
	}
}

// Traverses the set of arcs in the graph, passing each arc to the
// provided closure.
func (g *weightedDirected) EachArc(f ArcStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, weight := range adjacent {
			if f(NewWeightedArc(source, target, weight)) {
				return
			}
		}
	}
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *weightedUndirected) EachEdgeIncidentTo(v Vertex, f EdgeStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, weight := range g.list[v] {
		if f(NewWeightedEdge(v, adjacent, weight)) {
			return
		}
	}
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *weightedUndirected) EachAdjacentTo(vertex Vertex, f VertexStep) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachVertexInAdjacencyList(g.list, vertex, f)
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge weight.
func (g *weightedUndirected) HasEdge(edge Edge) bool {
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

// Indicates whether or not the given weighted edge is present in the graph.
// It will only match if the provided WeightedEdge has the same weight as
// the edge contained in the graph.
func (g *weightedUndirected) HasWeightedEdge(edge WeightedEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Spread it into two expressions to avoid evaluating the second if possible
	u, v := edge.Both()
	if weight, exists := g.list[u][v]; exists {
		return edge.Weight() == weight
	} else if weight, exists := g.list[v][u]; exists {
		return edge.Weight() == weight
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *weightedUndirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *weightedUndirected) RemoveVertex(vertices ...Vertex) {
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
func (g *weightedUndirected) AddEdges(edges ...WeightedEdge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *weightedUndirected) addEdges(edges ...WeightedEdge) {
	for _, edge := range edges {
		u, v := edge.Both()
		g.ensureVertex(u, v)

		if _, exists := g.list[u][v]; !exists {
			w := edge.Weight()
			g.list[u][v] = w
			g.list[v][u] = w
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *weightedUndirected) RemoveEdges(edges ...WeightedEdge) {
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
