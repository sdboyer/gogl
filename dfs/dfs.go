// Contains algos and logic related to depth-first graph traversal.
package dfs

import (
	"errors"

	"github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

const (
	white = iota
	grey
	black
)

// Performs a depth-first search for the given target vertex in the provided graph, beginning
// from the given start vertex.
//
// A slice of vertices is returned, identifying the path from the start to the target vertex.
// If no path can be found, the returned slice is nil and an error is returned instead.
func Search(g gogl.Graph, target gogl.Vertex, start gogl.Vertex) (path []gogl.Vertex, err error) {
	if !g.HasVertex(target) {
		return nil, errors.New("Target vertex is not present in graph.")
	}
	if !g.HasVertex(start) {
		return nil, errors.New("Start vertex is not present in graph.")
	}

	visitor := &searchVisitor{}

	w := walker{
		vis:    visitor,
		g:      g,
		colors: make(map[gogl.Vertex]uint),
		target: target,
	}

	if dg, ok := g.(gogl.Digraph); ok {
		w.dg = dg
	}

	w.dfsearch(start)

	return visitor.getPath(), nil
}

// Performs a depth-first search for the given target vertex in the provided graph, using the vertices
// provided to the third variadic parameter as starting points. For each given start vertex,
// the search proceeds in a parallel goroutine until a path to the target vertex is found.
// The resulting path is then sent back through the channel.
//
// If no starting vertices are provided, then a list of source vertices is built via FindSources(),
// and that set is used as the starting point. Because FindSources() requires a Digraph,
// an error will be returned if a non-directed graph is provided without any start vertices.
//
// TODO unexported until this is actually implemented. keeping it here as a note for now :)
func multiSearch(g gogl.Graph, target gogl.Vertex, start ...gogl.Vertex) (pathchan chan<- *searchPath, err error) {
	start, err = buildStartQueue(g)
	if err != nil {
		return nil, err
	}

	return
}

type searchPath struct {
	Start     gogl.Vertex
	Target    gogl.Vertex
	Reachable bool
	Path      []gogl.Vertex
}

// Performs a topological sort using the provided graph. The third variadic parameter defines a
// set of vertices to start from; vertices unreachable from this starting set will not be included
// in the final sorted output.
//
// If no starting vertices are provided, then a list of source vertices is built via FindSources(),
// and that set is used as the starting point. Because FindSources() requires a Digraph,
// an error will be returned if a non-directed graph is provided without any start vertices.
func Toposort(g gogl.Graph, start ...gogl.Vertex) ([]gogl.Vertex, error) {
	start, err := buildStartQueue(g, start...)
	if err != nil {
		return nil, err
	}

	stack := vstack{}
	for _, v := range start {
		stack.push(v)
	}

	if stack.length() == 0 {
		return nil, errors.New("No vertices provided as start points, cannot traverse.")
	}

	// Set the tsl capacity to the order of the graph. May be bigger than we need, but def not smaller
	var capacity int
	if c, ok := g.(gogl.VertexCounter); ok {
		capacity = c.Order()
	} else {
		capacity = 32
	}

	visitor := &TslVisitor{tsl: make([]gogl.Vertex, 0, capacity)}

	w := &walker{
		vis:    visitor,
		g:      g,
		colors: make(map[gogl.Vertex]uint),
	}

	var traverser func(*walker, gogl.Vertex)
	if dg, ok := g.(gogl.Digraph); ok {
		w.dg = dg
		traverser = (*walker).dftraverse
	} else {
		traverser = (*walker).dfutraverse
	}

	for v := stack.pop(); ; {
		traverser(w, v)

		if stack.length() == 0 {
			break
		}
	}

	return visitor.GetTsl()
}

// Traverses the given graph in a depth-first manner, using the given visitor
// and starting from the given vertices.
func Traverse(g gogl.Graph, visitor Visitor, start ...gogl.Vertex) (Visitor, error) {
	start, err := buildStartQueue(g, start...)
	if err != nil {
		return nil, err
	}

	stack := vstack{}
	for _, v := range start {
		stack.push(v)
	}

	w := &walker{
		vis:    visitor,
		g:      g,
		colors: make(map[gogl.Vertex]uint),
	}

	if dg, ok := g.(gogl.Digraph); ok {
		w.dg = dg
	}

	for v := stack.pop(); ; {
		w.dftraverse(v)

		if stack.length() == 0 {
			break
		}
	}

	return visitor, nil
}

// Finds all source vertices (vertices with no incoming edges) in the given directed graph.
func FindSources(g gogl.Digraph) (sources []gogl.Vertex, err error) {
	// TODO hardly the most efficient way to keep track, i'm sure
	incomings := set.New(set.NonThreadSafe)

	g.Arcs(func(e gogl.Arc) (terminate bool) {
		incomings.Add(e.Target())
		return
	})

	g.Vertices(func(v gogl.Vertex) (terminate bool) {
		if !incomings.Has(v) {
			sources = append(sources, v)
		}
		return
	})

	return
}

// Simple helper for shared traversal entry-point logic.
func buildStartQueue(g gogl.Graph, v ...gogl.Vertex) (start []gogl.Vertex, err error) {
	if len(v) == 0 {
		if dg, ok := g.(gogl.Digraph); ok {
			start, err = FindSources(dg)
		} else {
			return nil, errors.New("Undirected graphs do not have sources, a start point for traversal must be provided.")
		}
	} else {
		start = v
	}

	return
}

type vnode struct {
	v    gogl.Vertex
	next *vnode
}

type vqueue struct {
	front *vnode
	back  *vnode
	count int
}

type vstack struct {
	top   *vnode
	count int
}

type linkedlist interface {
	push(v gogl.Vertex)
	pop() gogl.Vertex
	length() int
}

func (q *vqueue) push(v gogl.Vertex) {
	n := &vnode{v: v}

	if q.back == nil {
		q.front = n
		q.back = n
	} else {
		q.back.next = n
		q.back = n
	}

	q.count++
}

func (q *vqueue) pop() gogl.Vertex {
	if q.front == nil {
		return nil
	}

	ret := q.front
	q.front = q.front.next

	q.count--
	return ret.v
}

func (q *vqueue) length() int {
	return q.count
}

func (s *vstack) push(v gogl.Vertex) {
	n := &vnode{v: v}

	if s.top == nil {
		s.top = n
	} else {
		n.next = s.top
		s.top = n
	}

	s.count++
}

func (s *vstack) pop() gogl.Vertex {
	if s.top == nil {
		return nil
	}

	ret := s.top
	s.top = s.top.next

	s.count--
	return ret.v
}

func (s *vstack) length() int {
	return s.count
}

type Visitor interface {
	OnBackEdge(vertex gogl.Vertex)
	OnStartVertex(vertex gogl.Vertex)
	OnExamineEdge(edge gogl.Edge)
	OnFinishVertex(vertex gogl.Vertex)
}

type searchVisitor struct {
	g     gogl.Graph
	stack vstack
}

func (sv *searchVisitor) OnBackEdge(vertex gogl.Vertex) {}

func (sv *searchVisitor) OnStartVertex(vertex gogl.Vertex) {
	sv.stack.push(vertex)
}

func (sv *searchVisitor) OnExamineEdge(edge gogl.Edge) {}

func (sv *searchVisitor) OnFinishVertex(vertex gogl.Vertex) {
	sv.stack.pop()
}

func (sv *searchVisitor) getPath() []gogl.Vertex {
	path := make([]gogl.Vertex, 0, sv.stack.length())

	stacklen := sv.stack.length()
	for i := 0; i < stacklen; i++ {
		vertex := sv.stack.pop()
		path = append(path, vertex)
	}

	return path
}

type TslVisitor struct {
	g   gogl.Graph
	tsl []gogl.Vertex
	err error
}

func (vis *TslVisitor) OnBackEdge(vertex gogl.Vertex) {
	vis.err = errors.New("Cycle detected in graph")
}

func (vis *TslVisitor) OnStartVertex(vertex gogl.Vertex) {}

func (vis *TslVisitor) OnExamineEdge(edge gogl.Edge) {}

func (vis *TslVisitor) OnFinishVertex(vertex gogl.Vertex) {
	vis.tsl = append(vis.tsl, vertex)
}

func (vis *TslVisitor) GetTsl() ([]gogl.Vertex, error) {
	return vis.tsl, vis.err
}

type walker struct {
	vis      Visitor
	g        gogl.Graph
	dg       gogl.Digraph
	complete bool
	target   gogl.Vertex
	// TODO is there ANY way to do this more efficiently without mutating/coloring the vertex objects directly? this means lots of hashtable lookups
	colors map[gogl.Vertex]uint
	ll     linkedlist
}

func (w *walker) dftraverse(v gogl.Vertex) {
	color, exists := w.colors[v]
	if !exists {
		color = white
	}

	if color == grey {
		w.vis.OnBackEdge(v)
	} else if color == white {
		w.colors[v] = grey
		w.vis.OnStartVertex(v)

		w.dg.ArcsFrom(v, func(e gogl.Arc) (terminate bool) {
			w.vis.OnExamineEdge(e)
			w.dftraverse(e.Target())
			return
		})

		w.vis.OnFinishVertex(v)
		w.colors[v] = black
	}
}

func (w *walker) dfsearch(v gogl.Vertex) {
	if v == w.target {
		w.complete = true
		w.vis.OnStartVertex(v)
		return
	}

	color, exists := w.colors[v]
	if !exists {
		color = white
	}

	if color == grey {
		w.vis.OnBackEdge(v)
	} else if color == white {
		w.colors[v] = grey
		w.vis.OnStartVertex(v)

		w.dg.ArcsFrom(v, func(e gogl.Arc) bool {
			// no more new visits if complete
			if !w.complete {
				w.vis.OnExamineEdge(e)
				w.dfsearch(e.Target())
			}
			return w.complete
		})
		// escape hatch
		if w.complete {
			return
		}

		w.vis.OnFinishVertex(v)
		w.colors[v] = black
	}
}

func (w *walker) dfutraverse(v gogl.Vertex) {
	if _, exists := w.colors[v]; !exists {
		w.colors[v] = grey
		w.vis.OnStartVertex(v)

		w.g.IncidentTo(v, func(e gogl.Edge) (terminate bool) {
			w.vis.OnExamineEdge(e)
			v1, v2 := e.Both()
			if v == v1 {
				w.dfutraverse(v2)
			} else {
				w.dfutraverse(v1)
			}
			return
		})

		w.vis.OnFinishVertex(v)
		w.colors[v] = black
	}
}
