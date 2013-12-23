package gogl

// Contains algos and logic related to depth-first traversal.

type vnode struct {
	v Vertex
	next *vnode
}

type vqueue struct {
	front *vnode
	back *vnode
	count int
}

type vstack struct {
	top *vnode
	count int
}

type linkedlist interface {
	push(v Vertex)
	pop() Vertex
	length() int
}

func (q *vqueue) push(v Vertex) {
	n  := &vnode{v: v}

	if q.back == nil {
		q.front = n
		q.back = n
	} else {
		q.back.next = n
		q.back = n
	}

	q.count++
}

func (q *vqueue) pop() Vertex {
	ret := q.front
	q.front = q.front.next

	q.count--
	return ret.v
}

func (q *vqueue) length() int {
	return q.count
}

func (s *vstack) push(v Vertex) {
	n  := &vnode{v: v}

	if s.top == nil {
		s.top = n
	} else {
		n.next = s.top
		s.top = n
	}

	s.count++
}

func (s *vstack) pop() Vertex {
	ret := s.top
	s.top = s.top.next

	s.count--
	return ret.v
}

func (s *vstack) length() int {
	return s.count
}

type walker struct {
	vis DFVisitor
	walker dfWalker
	g Graph
	visited map[Vertex]struct{}
	visiting map[Vertex]struct{}
	ll linkedlist
}

type DFVisitor interface {
	BindGraph(g Graph) bool
	HasGraph() bool
	OnInitializeVertex(vertex Vertex)
	OnBackEdge(vertex Vertex)
	OnStartVertex(vertex Vertex)
	OnExamineEdge(edge Edge)
	OnFinishVertex(vertex Vertex)
}

type dfWalker func (v Vertex)

func (w *walker) dfrecursive(v Vertex) {
	if _, visiting := w.visiting[v]; visiting {
		w.vis.OnBackEdge(v)
	} else if _, visited := w.visited[v]; !visited {
		w.visiting[v] = struct{}{}
		w.vis.OnStartVertex(v)

		w.g.EachAdjacent(v, func(to Vertex) {
			w.vis.OnExamineEdge(&BaseEdge{v, to})
			w.dfrecursive(to)
		})

		w.vis.OnFinishVertex(v)
		w.visited[v] = struct{}{}
		delete(w.visiting, v)
	}
}

func (w *walker) dflist() {
	var v Vertex
	for v = w.ll.pop(); v != nil; v = w.ll.pop() {
		if _, visiting := w.visiting[v]; visiting {
			w.vis.OnBackEdge(v)
		} else {
			
		}
	}
}


