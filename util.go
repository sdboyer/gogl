package gogl

/* Enumerator to slice/collection functors */

// Internal interface used for granular checks on whether a graphish can report vertex count.
type vertex_counter interface {
	Order() int
}

// Internal interface used for granular checks on whether a graphish can report edge count.
type edge_counter interface {
	Size() int
}

// Collects all of a graph's vertices into a vertex slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectVertices(g VertexEnumerator) (vertices []Vertex) {
	if c, ok := g.(vertex_counter); ok {
		// If possible, size the slice based on the number of vertices the graph reports it has
		vertices = make([]Vertex, 0, c.Order())
	} else {
		// Otherwise just pick something...reasonable?
		vertices = make([]Vertex, 0, 32)
	}

	g.EachVertex(func(v Vertex) (terminate bool) {
		vertices = append(vertices, v)
		return
	})

	return vertices
}

// Collects all of a given vertex's adjacent vertices into a vertex slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectVerticesAdjacentTo(v Vertex, g AdjacencyEnumerator) (vertices []Vertex) {
	if c, ok := g.(DegreeChecker); ok {
		// If possible, size the slice based on the number of adjacent vertices the graph reports
		deg, _ := c.DegreeOf(v)
		vertices = make([]Vertex, 0, deg)
	} else {
		// Otherwise just pick something...reasonable?
		vertices = make([]Vertex, 0, 8)
	}

	g.EachAdjacentTo(v, func(v Vertex) (terminate bool) {
		vertices = append(vertices, v)
		return
	})

	return vertices
}

// Collects all of a graph's edges into an edge slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectEdges(g EdgeEnumerator) (edges []Edge) {
	if c, ok := g.(edge_counter); ok {
		// If possible, size the slice based on the number of edges the graph reports it has
		edges = make([]Edge, 0, c.Size())
	} else {
		// Otherwise just pick something...reasonable?
		edges = make([]Edge, 0, 32)
	}

	g.EachEdge(func(e Edge) (terminate bool) {
		edges = append(edges, e)
		return
	})

	return edges
}

// Collects all of a given vertex's incident edges into an edge slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectEdgesIncidentTo(v Vertex, g IncidentEdgeEnumerator) (edges []Edge) {
	if c, ok := g.(DegreeChecker); ok {
		// If possible, size the slice based on the number of incident edges the graph reports
		deg, _ := c.DegreeOf(v)
		edges = make([]Edge, 0, deg)
	} else {
		// Otherwise just pick something...reasonable?
		edges = make([]Edge, 0, 8)
	}

	g.EachEdgeIncidentTo(v, func(e Edge) (terminate bool) {
		edges = append(edges, e)
		return
	})

	return edges
}

// Collects all of a given vertex's out-edges into an edge slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectArcsFrom(v Vertex, g IncidentArcEnumerator) (edges []Edge) {
	if c, ok := g.(DirectedDegreeChecker); ok {
		// If possible, size the slice based on the number of out-edges the graph reports
		deg, _ := c.OutDegreeOf(v)
		edges = make([]Edge, 0, deg)
	} else {
		// Otherwise just pick something...reasonable?
		edges = make([]Edge, 0, 8)
	}

	g.EachArcFrom(v, func(e Edge) (terminate bool) {
		edges = append(edges, e)
		return
	})

	return edges
}

// Collects all of a given vertex's in-edges into an edge slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectArcsTo(v Vertex, g IncidentArcEnumerator) (edges []Edge) {
	if c, ok := g.(DirectedDegreeChecker); ok {
		// If possible, size the slice based on the number of in-edges the graph reports
		deg, _ := c.InDegreeOf(v)
		edges = make([]Edge, 0, deg)
	} else {
		// Otherwise just pick something...reasonable?
		edges = make([]Edge, 0, 8)
	}

	g.EachArcTo(v, func(e Edge) (terminate bool) {
		edges = append(edges, e)
		return
	})

	return edges
}
