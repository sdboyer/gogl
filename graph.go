package gogl

import "fmt"

// for great justice
var fml = fmt.Println

/* Composite graph interfaces */

// Graph is gogl's most basic interface: it contains only the methods that
// *every* type of graph implements.
//
// Graph is intentionally underspecified: both directed and undirected graphs
// implement it; simple graphs, multigraphs, weighted, labeled, or any
// combination thereof.
//
// The semantics of some of these methods vary slightly from one graph type
// to another, but in general, the basic Graph methods are supplemented, not
// superceded, by the methods in more specific interfaces.
//
// Graph is a purely read oriented interface; the various Mutable*Graph
// interfaces contain the methods for writing.
type Graph interface {
	VertexEnumerator        // Enumerates vertices to an injected lambda
	EdgeEnumerator          // Enumerates edges to an injected lambda
	AdjacencyEnumerator     // Enumerates a vertex's adjacent vertices to an injected lambda
	IncidentEdgeEnumerator  // Enumerates a vertex's incident edges to an injected lambda
	VertexMembershipChecker // Allows inspection of contained vertices
	EdgeMembershipChecker   // Allows inspection of contained edges
	DegreeChecker           // Reports degree of vertices
	Order() int             // Reports total number of vertices in the graph
	Size() int              // Reports total number of edges in the graph
}

// GraphSource is a subinterface of Graph, describing the minimal set of methods
// necessary to accomplish a naive full graph traversal and copy.
type GraphSource interface {
	VertexEnumerator
	EdgeEnumerator
	Order() int
}

// DirectedGraph describes a Graph all of whose edges are directed.
//
// gogl treats edge directionality as a property of the graph, not the edge itself.
// Thus, implementing this interface is gogl's only signal that a graph's edges are directed.
type DirectedGraph interface {
	Graph
	IncidentArcEnumerator // Enumerates a vertex's incident in- and out-arcs to an injected lambda
	DirectedDegreeChecker // Reports in- and out-degree of vertices
	Transposer            // DirectedGraphs can produce a transpose of themselves
}

// MutableGraph describes a graph with basic edges (no weighting, labeling, etc.)
// that can be modified freely by adding or removing vertices or edges.
type MutableGraph interface {
	Graph
	VertexSetMutator
	EdgeSetMutator
}

// A simple graph is in opposition to a multigraph: it disallows loops and
// parallel edges.
type SimpleGraph interface {
	Graph
	Density() float64
}

// A weighted graph is a graph subtype where the edges have a numeric weight;
// as described by the WeightedEdge interface, this weight is a signed int.
//
// WeightedGraphs have both the HasEdge() and HasWeightedEdge() methods.
// Correct implementations should treat the difference as a matter of strictness:
//
// HasEdge() should return true as long as an edge exists
// connecting the two given vertices (respecting directed or undirected as
// appropriate), regardless of its weight.
//
// HasWeightedEdge() should return true iff an edge exists connecting the
// two given vertices (respecting directed or undirected as appropriate),
// AND if the edge weights are the same.
type WeightedGraph interface {
	Graph
	HasWeightedEdge(e WeightedEdge) bool
}

// MutableWeightedGraph is the mutable version of a weighted graph. Its
// AddEdges() method is incompatible with MutableGraph, guaranteeing
// only weighted edges can be present in the graph.
type MutableWeightedGraph interface {
	WeightedGraph
	VertexSetMutator
	WeightedEdgeSetMutator
}

// A labeled graph is a graph subtype where the edges have an identifier;
// as described by the LabeledEdge interface, this identifier is a string.
//
// LabeledGraphs have both the HasEdge() and HasLabeledEdge() methods.
// Correct implementations should treat the difference as a matter of strictness:
//
// HasEdge() should return true as long as an edge exists
// connecting the two given vertices (respecting directed or undirected as
// appropriate), regardless of its label.
//
// HasLabeledEdge() should return true iff an edge exists connecting the
// two given vertices (respecting directed or undirected as appropriate),
// AND if the edge labels are the same.
type LabeledGraph interface {
	Graph
	HasLabeledEdge(e LabeledEdge) bool
}

// LabeledWeightedGraph is the mutable version of a labeled graph. Its
// AddEdges() method is incompatible with MutableGraph, guaranteeing
// only labeled edges can be present in the graph.
type MutableLabeledGraph interface {
	LabeledGraph
	VertexSetMutator
	LabeledEdgeSetMutator
}

// A data graph is a graph subtype where the edges carry arbitrary Go data;
// as described by the DataEdge interface, this identifier is an interface{}.
//
// DataGraphs have both the HasEdge() and HasDataEdge() methods.
// Correct implementations should treat the difference as a matter of strictness:
//
// HasEdge() should return true as long as an edge exists
// connecting the two given vertices (respecting directed or undirected as
// appropriate), regardless of its label.
//
// HasDataEdge() should return true iff an edge exists connecting the
// two given vertices (respecting directed or undirected as appropriate),
// AND if the edge data is the same. Simple comparison will typically be used
// to establish data equality, which means that using noncomparables (a slice,
// map, or non-pointer struct containing a slice or a map) for the data will
// cause a panic.
type DataGraph interface {
	Graph
	HasDataEdge(e DataEdge) bool
}

// MutableDataGraph is the mutable version of a propety graph. Its
// AddEdges() method is incompatible with MutableGraph, guaranteeing
// only property edges can be present in the graph.
type MutableDataGraph interface {
	DataGraph
	VertexSetMutator
	DataEdgeSetMutator
}

/* Atomic graph interfaces */

// EdgeLambdas are used as arguments to various enumerators. They are called once for each edge produced by the enumerator.
//
// If the lambda returns true, the calling enumerator is expected to end enumeration and return control to its caller.
type EdgeLambda func(Edge) (terminate bool)

// VertexLambdas are used as arguments to various enumerators. They are called once for each vertex produced by the enumerator.
//
// If the lambda returns true, the calling enumerator is expected to end enumeration and return control to its caller.
type VertexLambda func(Vertex) (terminate bool)

// A VertexEnumerator iteratively enumerates vertices, and can indicate the number of vertices present.
type VertexEnumerator interface {
	EachVertex(VertexLambda)
}

// An EdgeEnumerator iteratively enumerates edges, and can indicate the number of edges present.
type EdgeEnumerator interface {
	EachEdge(EdgeLambda)
}

// An IncidentEdgeEnumerator iteratively enumerates a given vertex's incident edges.
type IncidentEdgeEnumerator interface {
	EachEdgeIncidentTo(v Vertex, incidentEdgeLambda EdgeLambda)
}

// An IncidentArcEnumerator iteratively enumerates a given vertex's incident arcs (directed edges).
// One enumerator provides inbound edges, the other outbound edges.
type IncidentArcEnumerator interface {
	EachArcFrom(v Vertex, outEdgeLambda EdgeLambda)
	EachArcTo(v Vertex, inEdgeLambda EdgeLambda)
}

// An AdjacencyEnumerator iteratively enumerates a given vertex's adjacent vertices.
type AdjacencyEnumerator interface {
	EachAdjacentTo(start Vertex, adjacentVertexLambda VertexLambda)
}

// A VertexMembershipChecker can indicate the presence of a vertex.
type VertexMembershipChecker interface {
	HasVertex(Vertex) bool // Whether or not the vertex is present in the set
}

// A DegreeChecker reports the number of edges incident to a given vertex.
type DegreeChecker interface {
	DegreeOf(Vertex) (degree int, exists bool) // Number of incident edges; if vertex is present
}

// A DirectedDegreeChecker reports the number of in or out-edges incident to given vertex.
type DirectedDegreeChecker interface {
	InDegreeOf(Vertex) (degree int, exists bool)  // Number of in-edges; if vertex is present
	OutDegreeOf(Vertex) (degree int, exists bool) // Number of out-edges; if vertex is present
}

// An EdgeMembershipChecker can indicate the presence of an edge.
type EdgeMembershipChecker interface {
	HasEdge(Edge) bool
}

// A VertexSetMutator allows the addition and removal of vertices from a set.
type VertexSetMutator interface {
	EnsureVertex(...Vertex)
	RemoveVertex(...Vertex)
}

// An EdgeSetMutator allows the addition and removal of edges from a set.
type EdgeSetMutator interface {
	AddEdges(edges ...Edge)
	RemoveEdges(edges ...Edge)
}

// A WeightedEdgeSetMutator allows the addition and removal of weighted edges from a set.
type WeightedEdgeSetMutator interface {
	AddEdges(edges ...WeightedEdge)
	RemoveEdges(edges ...WeightedEdge)
}

// A LabeledEdgeSetMutator allows the addition and removal of labeled edges from a set.
type LabeledEdgeSetMutator interface {
	AddEdges(edges ...LabeledEdge)
	RemoveEdges(edges ...LabeledEdge)
}

// A DataEdgeSetMutator allows the addition and removal of data edges from a set.
type DataEdgeSetMutator interface {
	AddEdges(edges ...DataEdge)
	RemoveEdges(edges ...DataEdge)
}

// A Transposer produces a transposed version of a DirectedGraph.
type Transposer interface {
	Transpose() DirectedGraph
}
