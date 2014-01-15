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

type DFVisitor interface {
	onInitializeVertex(vertex Vertex)
	onBackEdge(vertex Vertex)
	onStartVertex(vertex Vertex)
	onExamineEdge(edge Edge)
	onFinishVertex(vertex Vertex)
}

type DFTslVisitor struct {
	g   Graph
	tsl []Vertex
}

func (vis *DFTslVisitor) onInitializeVertex(vertex Vertex) {}

func (vis *DFTslVisitor) onBackEdge(vertex Vertex) {}

func (vis *DFTslVisitor) onStartVertex(vertex Vertex) {}

func (vis *DFTslVisitor) onExamineEdge(edge Edge) {}

func (vis *DFTslVisitor) onFinishVertex(vertex Vertex) {
	vis.tsl = append(vis.tsl, vertex)
}

func (vis *DFTslVisitor) GetTsl() []Vertex {
	return vis.tsl
}

type EdgeFilterer interface {
	FilterEdge(edge Edge) bool
	FilterEdges(edges []Edge) []Edge
}

type walker struct {
	vis      DFVisitor
	g        Graph
	visited  map[Vertex]struct{}
	visiting map[Vertex]struct{}
	ll       linkedlist
}

func DepthFirstFromVertices(g Graph, vis DFVisitor, vertices ...Vertex) error {
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
		visited:  make(map[Vertex]struct{}),
		visiting: make(map[Vertex]struct{}),
	}

	for v := stack.pop(); ; {
		w.dfrecursive(v)

		if stack.length() == 0 {
			break
		}
	}

	return nil
}

func (w *walker) dfrecursive(v Vertex) {
	if _, visiting := w.visiting[v]; visiting {
		w.vis.onBackEdge(v)
	} else if _, visited := w.visited[v]; !visited {
		w.visiting[v] = struct{}{}
		w.vis.onStartVertex(v)

		w.g.EachAdjacent(v, func(to Vertex) {
			w.vis.onExamineEdge(&BaseEdge{v, to})
			w.dfrecursive(to)
		})

		w.vis.onFinishVertex(v)
		w.visited[v] = struct{}{}
		delete(w.visiting, v)
	}
}

func (w *walker) dflist() {
	var v Vertex
	for v = w.ll.pop(); v != nil; v = w.ll.pop() {
		if _, visiting := w.visiting[v]; visiting {
			w.vis.onBackEdge(v)
		} else {

		}
	}
}
