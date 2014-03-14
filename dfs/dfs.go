// Contains algos and logic related to depth-first graph traversal.
package dfs

import (
	"errors"

	"github.com/sdboyer/gogl"
)

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
	onInitializeVertex(vertex gogl.Vertex)
	onBackEdge(vertex gogl.Vertex)
	onStartVertex(vertex gogl.Vertex)
	onExamineEdge(edge gogl.Edge)
	onFinishVertex(vertex gogl.Vertex)
}

type DFTslVisitor struct {
	g   gogl.Graph
	tsl []gogl.Vertex
}

func (vis *DFTslVisitor) onInitializeVertex(vertex gogl.Vertex) {}

func (vis *DFTslVisitor) onBackEdge(vertex gogl.Vertex) {}

func (vis *DFTslVisitor) onStartVertex(vertex gogl.Vertex) {}

func (vis *DFTslVisitor) onExamineEdge(edge gogl.Edge) {}

func (vis *DFTslVisitor) onFinishVertex(vertex gogl.Vertex) {
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
		w.vis.onBackEdge(v)
	} else if _, visited := w.visited[v]; !visited {
		w.visiting[v] = struct{}{}
		w.vis.onStartVertex(v)

		w.g.EachAdjacent(v, func(to gogl.Vertex) {
			w.vis.onExamineEdge(gogl.BaseEdge{v, to})
			w.dfrecursive(to)
		})

		w.vis.onFinishVertex(v)
		w.visited[v] = struct{}{}
		delete(w.visiting, v)
	}
}

func (w *walker) dflist() {
	var v gogl.Vertex
	for v = w.ll.pop(); v != nil; v = w.ll.pop() {
		if _, visiting := w.visiting[v]; visiting {
			w.vis.onBackEdge(v)
		} else {

		}
	}
}
