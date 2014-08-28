package gogl

// Returns the number of vertices in a graph.
//
// If available, this function will take advantage of the optional optimization Order() method.
// Otherwise, it will iterate through all vertices in the graph. Thus, if your use case involves
// iterating through all the graph's vertices, it is better to simply check for the VertexCounter
// interface yourself.
func Order(g VertexEnumerator) int {
	if c, ok := g.(VertexCounter); ok {
		return c.Order()
	} else {
		var order int
		g.Vertices(func(v Vertex) (terminate bool) {
			order++
			return
		})

		return order
	}
}

// Returns the number of edges in a graph.
//
// If available, this function will take advantage of the optional optimization Size() method.
// Otherwise, it will iterate through all edges in the graph. Thus, if your use case involves
// iterating through all the graph's edges, it is better to simply check for the EdgeCounter
// interface yourself.
func Size(g EdgeEnumerator) int {
	if c, ok := g.(EdgeCounter); ok {
		return c.Size()
	} else {
		var size int
		g.EachEdge(func(e Edge) (terminate bool) {
			size++
			return
		})
		return size
	}
}

/* Enumerator to slice/collection functors */

// Collects all of a graph's vertices into a vertex slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectVertices(g VertexEnumerator) (vertices []Vertex) {
	if c, ok := g.(VertexCounter); ok {
		// If possible, size the slice based on the number of vertices the graph reports it has
		vertices = make([]Vertex, 0, c.Order())
	} else {
		// Otherwise just pick something...reasonable?
		vertices = make([]Vertex, 0, 32)
	}

	g.Vertices(func(v Vertex) (terminate bool) {
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
	if c, ok := g.(EdgeCounter); ok {
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

	g.IncidentTo(v, func(e Edge) (terminate bool) {
		edges = append(edges, e)
		return
	})

	return edges
}

// Collects all of a given vertex's out-arcs into an arc slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectArcsFrom(v Vertex, g IncidentArcEnumerator) (arcs []Arc) {
	if c, ok := g.(DirectedDegreeChecker); ok {
		// If possible, size the slice based on the number of out-arcs the graph reports
		deg, _ := c.OutDegreeOf(v)
		arcs = make([]Arc, 0, deg)
	} else {
		// Otherwise just pick something...reasonable?
		arcs = make([]Arc, 0, 8)
	}

	g.ArcsFrom(v, func(e Arc) (terminate bool) {
		arcs = append(arcs, e)
		return
	})

	return arcs
}

// Collects all of a given vertex's in-arcs into an arc slice, for easy range-ing.
//
// This is a convenience function. Avoid it on very large graphs or in performance critical sections.
func CollectArcsTo(v Vertex, g IncidentArcEnumerator) (arcs []Arc) {
	if c, ok := g.(DirectedDegreeChecker); ok {
		// If possible, size the slice based on the number of in-arcs the graph reports
		deg, _ := c.InDegreeOf(v)
		arcs = make([]Arc, 0, deg)
	} else {
		// Otherwise just pick something...reasonable?
		arcs = make([]Arc, 0, 8)
	}

	g.EachArcTo(v, func(e Arc) (terminate bool) {
		arcs = append(arcs, e)
		return
	})

	return arcs
}
