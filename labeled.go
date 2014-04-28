package gogl

import (
	"sync"

	"gopkg.in/fatih/set.v0"
)

// This is implemented as an adjacency list, because those are simple.
type baseLabeled struct {
	list map[Vertex]map[Vertex]string
	size int
	mu   sync.RWMutex
}

/* baseLabeled shared methods */

// Traverses the graph's vertices in random order, passing each vertex to the
// provided closure.
func (g *baseLabeled) EachVertex(f VertexLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for v := range g.list {
		if f(v) {
			return
		}
	}
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseLabeled) HasVertex(vertex Vertex) (exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	exists = g.hasVertex(vertex)
	return
}

// Indicates whether or not the given vertex is present in the graph.
func (g *baseLabeled) hasVertex(vertex Vertex) (exists bool) {
	_, exists = g.list[vertex]
	return
}

// Returns the order (number of vertices) in the graph.
func (g *baseLabeled) Order() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return len(g.list)
}

// Returns the size (number of edges) in the graph.
func (g *baseLabeled) Size() int {
	return g.size
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseLabeled) EnsureVertex(vertices ...Vertex) {
	if len(vertices) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.ensureVertex(vertices...)
}

// Adds the provided vertices to the graph. If a provided vertex is
// already present in the graph, it is a no-op (for that vertex only).
func (g *baseLabeled) ensureVertex(vertices ...Vertex) {
	for _, vertex := range vertices {
		if !g.hasVertex(vertex) {
			// TODO experiment with different lengths...possibly by analyzing existing density?
			g.list[vertex] = make(map[Vertex]string, 10)
		}
	}

	return
}

/* DirectedLabeled implementation */

type labeledDirected struct {
	baseLabeled
}

func NewLabeledDirected() MutableLabeledGraph {
	list := &labeledDirected{}
	// Cannot assign to promoted fields in a composite literals.
	list.list = make(map[Vertex]map[Vertex]string)

	// Type assertions to ensure interfaces are met
	var _ Graph = list
	var _ SimpleGraph = list
	var _ DirectedGraph = list
	var _ LabeledGraph = list
	var _ MutableLabeledGraph = list

	return list
}

// Returns the outdegree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *labeledDirected) OutDegreeOf(vertex Vertex) (degree int, exists bool) {
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
func (g *labeledDirected) InDegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		g.EachEdge(func(e Edge) (terminate bool) {
			if vertex == e.Target() {
				degree++
			}
			return
		})
	}

	return
}

// Returns the degree of the provided vertex, counting both in and out-edges.
func (g *labeledDirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	indegree, exists := g.InDegreeOf(vertex)
	outdegree, exists := g.OutDegreeOf(vertex)
	return indegree + outdegree, exists
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *labeledDirected) EachEdge(f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for source, adjacent := range g.list {
		for target, label := range adjacent {
			if f(BaseLabeledEdge{BaseEdge{U: source, V: target}, label}) {
				return
			}
		}
	}
}

// Enumerates the set of all edges incident to the provided vertex.
func (g *labeledDirected) EachEdgeIncidentTo(v Vertex, f EdgeLambda) {
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
func (g *labeledDirected) EachAdjacentTo(start Vertex, f VertexLambda) {
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
func (g *labeledDirected) EachArcFrom(v Vertex, f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, label := range g.list[v] {
		if f(BaseLabeledEdge{BaseEdge{U: v, V: adjacent}, label}) {
			return
		}
	}
}

// Enumerates the set of in-edges for the provided vertex.
func (g *labeledDirected) EachArcTo(v Vertex, f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for candidate, adjacent := range g.list {
		for target, label := range adjacent {
			if target == v {
				if f(BaseLabeledEdge{BaseEdge{U: candidate, V: target}, label}) {
					return
				}
			}
		}
	}
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge label.
func (g *labeledDirected) HasEdge(edge Edge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, exists := g.list[edge.Source()][edge.Target()]
	return exists
}

// Indicates whether or not the given labeled edge is present in the graph.
// It will only match if the provided LabeledEdge has the same label as
// the edge contained in the graph.
func (g *labeledDirected) HasLabeledEdge(edge LabeledEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if label, exists := g.list[edge.Source()][edge.Target()]; exists {
		return label == edge.Label()
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *labeledDirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *labeledDirected) RemoveVertex(vertices ...Vertex) {
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
func (g *labeledDirected) AddEdges(edges ...LabeledEdge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *labeledDirected) addEdges(edges ...LabeledEdge) {
	for _, edge := range edges {
		g.ensureVertex(edge.Source(), edge.Target())

		if _, exists := g.list[edge.Source()][edge.Target()]; !exists {
			g.list[edge.Source()][edge.Target()] = edge.Label()
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *labeledDirected) RemoveEdges(edges ...LabeledEdge) {
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

func (g *labeledDirected) Transpose() DirectedGraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g2 := &labeledDirected{}
	g2.list = make(map[Vertex]map[Vertex]string)

	// Guess at average indegree by looking at ratio of edges to vertices, use that to initially size the adjacency maps
	startcap := int(g.Size() / g.Order())

	for source, adjacent := range g.list {
		if !g2.hasVertex(source) {
			g2.list[source] = make(map[Vertex]string, startcap+1)
		}

		for target, label := range adjacent {
			if !g2.hasVertex(target) {
				g2.list[target] = make(map[Vertex]string, startcap+1)
			}
			g2.list[target][source] = label
		}
	}

	return g2
}

/* UndirectedLabeled implementation */

type labeledUndirected struct {
	baseLabeled
}

func NewLabeledUndirected() MutableLabeledGraph {
	g := &labeledUndirected{}
	// Cannot assign to promoted fields in a composite literals.
	g.list = make(map[Vertex]map[Vertex]string)

	// Type assertions to ensure interfaces are met
	var _ Graph = g
	var _ SimpleGraph = g
	var _ LabeledGraph = g
	var _ MutableLabeledGraph = g

	return g
}

// Returns the degree of the provided vertex. If the vertex is not present in the
// graph, the second return value will be false.
func (g *labeledUndirected) DegreeOf(vertex Vertex) (degree int, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if exists = g.hasVertex(vertex); exists {
		degree = len(g.list[vertex])
	}
	return
}

// Traverses the set of edges in the graph, passing each edge to the
// provided closure.
func (g *labeledUndirected) EachEdge(f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := set.NewNonTS()

	for source, adjacent := range g.list {
		for target, label := range adjacent {
			be := BaseEdge{U: source, V: target}
			e := BaseLabeledEdge{be, label}
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
func (g *labeledUndirected) EachEdgeIncidentTo(v Vertex, f EdgeLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.hasVertex(v) {
		return
	}

	for adjacent, label := range g.list[v] {
		if f(BaseLabeledEdge{BaseEdge{U: v, V: adjacent}, label}) {
			return
		}
	}
}

// Enumerates the vertices adjacent to the provided vertex.
func (g *labeledUndirected) EachAdjacentTo(vertex Vertex, f VertexLambda) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	eachAdjacentToUndirected(g.list, vertex, f)
}

// Indicates whether or not the given edge is present in the graph. It matches
// based solely on the presence of an edge, disregarding edge label.
func (g *labeledUndirected) HasEdge(edge Edge) bool {
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

// Indicates whether or not the given labeled edge is present in the graph.
// It will only match if the provided LabeledEdge has the same label as
// the edge contained in the graph.
func (g *labeledUndirected) HasLabeledEdge(edge LabeledEdge) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Spread it into two expressions to avoid evaluating the second if possible
	if label, exists := g.list[edge.Source()][edge.Target()]; exists {
		return edge.Label() == label
	} else if label, exists := g.list[edge.Target()][edge.Source()]; exists {
		return edge.Label() == label
	}
	return false
}

// Returns the density of the graph. Density is the ratio of edge count to the
// number of edges there would be in complete graph (maximum edge count).
func (g *labeledUndirected) Density() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order := g.Order()
	return 2 * float64(g.Size()) / float64(order*(order-1))
}

// Removes a vertex from the graph. Also removes any edges of which that
// vertex is a member.
func (g *labeledUndirected) RemoveVertex(vertices ...Vertex) {
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
func (g *labeledUndirected) AddEdges(edges ...LabeledEdge) {
	if len(edges) == 0 {
		return
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.addEdges(edges...)
}

// Adds a new edge to the graph.
func (g *labeledUndirected) addEdges(edges ...LabeledEdge) {
	for _, edge := range edges {
		g.ensureVertex(edge.Source(), edge.Target())

		if _, exists := g.list[edge.Source()][edge.Target()]; !exists {
			l := edge.Label()
			g.list[edge.Source()][edge.Target()] = l
			g.list[edge.Target()][edge.Source()] = l
			g.size++
		}
	}
}

// Removes edges from the graph. This does NOT remove vertex members of the
// removed edges.
func (g *labeledUndirected) RemoveEdges(edges ...LabeledEdge) {
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
