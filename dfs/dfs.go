// Contains algos and logic related to depth-first graph traversal.
package dfs

import (
	"errors"

	"github.com/fatih/set"
	"github.com/sdboyer/gogl"
)

const (
	white = iota
	grey
	black
)

// Performs a depth-first search on the provided graph for the given target vertex, beginning
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
		vis: visitor,
		g:   g,
		// TODO is there ANY way to do this more efficiently without mutating/coloring the vertex objects directly? this means lots of hashtable lookups
		colors: make(map[gogl.Vertex]uint),
		target: target,
	}

	w.dfsearch(start)

	return visitor.getPath(), nil
}

// Performs a depth-first search on the provided graph for the given vertex, using the vertices
// provided to the third variadic parameter as starting points. For each given start vertex,
// the search proceeds in a parallel goroutine until a path to the target vertex is found.
// The resulting path is then sent back through the channel.
//
// If no starting vertices are provided, then a list of source vertices is built via FindSources(),
// and that set is used as the starting point. Because FindSources() requires a DirectedGraph,
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
// and that set is used as the starting point. Because FindSources() requires a DirectedGraph,
// an error will be returned if a non-directed graph is provided without any start vertices.
func Toposort(g gogl.Graph, start ...gogl.Vertex) (tsl []gogl.Vertex, err error) {
	start, err = buildStartQueue(g)
	if err != nil {
		return nil, err
	}

	return
}

// Traverses the given graph in a depth-first manner, using the provided visitor.
func Traverse(g gogl.Graph, visitor Visitor) (Visitor, error) {

	return visitor, nil
}

// Finds all source vertices (vertices with no incoming edges) in the given directed graph.
func FindSources(g gogl.DirectedGraph) (sources []gogl.Vertex, err error) {
	// TODO hardly the most efficient way to keep track, i'm sure
	incomings := set.NewNonTS()

	g.EachEdge(func(e gogl.Edge) {
		incomings.Add(e.Target())
	})

	g.EachVertex(func(v gogl.Vertex) {
		if !incomings.Has(v) {
			sources = append(sources, v)
		}
	})

	return
}

// Simple helper for shared traversal entry-point logic.
func buildStartQueue(g gogl.Graph, v ...gogl.Vertex) (start []gogl.Vertex, err error) {
	if len(v) == 0 {
		if dg, ok := g.(gogl.DirectedGraph); ok {
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

type DFTslVisitor struct {
	g   gogl.Graph
	tsl []gogl.Vertex
}

func (vis *DFTslVisitor) OnInitializeVertex(vertex gogl.Vertex) {}

func (vis *DFTslVisitor) OnBackEdge(vertex gogl.Vertex) {}

func (vis *DFTslVisitor) OnStartVertex(vertex gogl.Vertex) {}

func (vis *DFTslVisitor) OnExamineEdge(edge gogl.Edge) {}

func (vis *DFTslVisitor) OnFinishVertex(vertex gogl.Vertex) {
	vis.tsl = append(vis.tsl, vertex)
}

func (vis *DFTslVisitor) GetTsl() []gogl.Vertex {
	return vis.tsl
}

type EdgeFilterer interface {
	FilterEdge(edge gogl.Edge) bool
	FilterEdges(edges []gogl.Edge) []gogl.Edge
}

type walker struct {
	vis      Visitor
	g        gogl.Graph
	complete bool
	target   gogl.Vertex
	colors   map[gogl.Vertex]uint
	visited  map[gogl.Vertex]struct{}
	visiting map[gogl.Vertex]struct{}
	ll       linkedlist
}

func DepthFirstFromVertices(g gogl.Graph, vis Visitor, vertices ...gogl.Vertex) error {
	stack := vstack{}
	for _, v := range vertices {
		stack.push(v)
	}

	if stack.length() == 0 {
		return errors.New("No vertices provided as start points, cannot traverse.")
	}

	w := walker{
		vis:      vis,
		g:        g,
		visited:  make(map[gogl.Vertex]struct{}),
		visiting: make(map[gogl.Vertex]struct{}),
	}

	for v := stack.pop(); ; {
		w.dfrecursive(v)

		if stack.length() == 0 {
			break
		}
	}

	return nil
}

func (w *walker) dfrecursive(v gogl.Vertex) {
	if _, visiting := w.visiting[v]; visiting {
		w.vis.OnBackEdge(v)
	} else if _, visited := w.visited[v]; !visited {
		w.visiting[v] = struct{}{}
		w.vis.OnStartVertex(v)

		w.g.EachAdjacent(v, func(to gogl.Vertex) {
			w.vis.OnExamineEdge(gogl.BaseEdge{v, to})
			w.dfrecursive(to)
		})

		w.vis.OnFinishVertex(v)
		w.visited[v] = struct{}{}
		delete(w.visiting, v)
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

		w.g.EachAdjacent(v, func(to gogl.Vertex) {
			// no more new visits if complete
			if !w.complete {
				w.vis.OnExamineEdge(gogl.BaseEdge{v, to})
				w.dfsearch(to)
			}
		})
		// escape hatch
		if w.complete {
			return
		}

		w.vis.OnFinishVertex(v)
		w.colors[v] = black
	}
}

func (w *walker) dflist() {
	var v gogl.Vertex
	for v = w.ll.pop(); v != nil; v = w.ll.pop() {
		if _, visiting := w.visiting[v]; visiting {
			w.vis.OnBackEdge(v)
		} else {

		}
	}
}
