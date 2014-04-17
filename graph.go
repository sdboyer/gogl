package gogl

import (
	"fmt"
)

/* Vertex structures */

// As a rule, gogl tries to place as low a requirement on its vertices as
// possible. This is because, from a purely graph theoretic perspective,
// vertices are inert. Boring, even. Graphs are more about the topology, the
// characteristics of the edges connecting the points than the points
// themselves. Your use case cares about the content of your vertices, but gogl
// does not.  Consequently, anything can act as a vertex.
type Vertex interface{}

/* Atomic graph interfaces */

// A VertexEnumerator iteratively enumerates vertices.
type VertexEnumerator interface {
	EachVertex(f func(Vertex))
}

// An EdgeEnumerator iteratively enumerates edges.
type EdgeEnumerator interface {
	EachEdge(f func(Edge))
}

// An AdjacencyEnumerator iteratively enumerates a given vertex's adjacent vertices.
type AdjacencyEnumerator interface {
	EachAdjacent(start Vertex, f func(adjacent Vertex))
}

// A VertexMembershipChecker can indicate the presence of a vertex.
type VertexMembershipChecker interface {
	HasVertex(Vertex) bool // Whether or not the vertex is present in the set
}

// An InOutDegreeChecker reports the number of edges incident to a given vertex.
// TODO use this
type DegreeChecker interface {
	Degree(Vertex) (degree int, exists bool) // Number of incident edges; if vertex is present
}

// An InOutDegreeChecker reports the number of in or out-edges a given vertex has.
type InOutDegreeChecker interface {
	InDegree(Vertex) (degree int, exists bool)  // Number of in-edges; if vertex is present
	OutDegree(Vertex) (degree int, exists bool) // Number of out-edges; if vertex is present
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

/* Aggregate graph interfaces */

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
	VertexEnumerator        // Allows enumerated traversal of vertices
	EdgeEnumerator          // Allows enumerated traversal of edges
	AdjacencyEnumerator     // Allows enumerated traversal of a vertex's adjacent vertices
	VertexMembershipChecker // Allows inspection of contained vertices
	EdgeMembershipChecker   // Allows inspection of contained edges
	InOutDegreeChecker      // Reports in- and out-degree of vertices
	Order() int             // Total number of vertices in the graph
	Size() int              // Total number of edges in the graph
}

// DirectedGraph describes a Graph all of whose edges are directed.
//
// Implementing DirectedGraph is the only unambiguous signal gogl provides
// that a graph's edges are directed.
type DirectedGraph interface {
	Graph
	Transposer // DirectedGraphs can produce a transpose of themselves
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
	EachWeightedEdge(f func(edge WeightedEdge))
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
	EachLabeledEdge(f func(edge LabeledEdge))
}

// LabeledWeightedGraph is the mutable version of a labeled graph. Its
// AddEdges() method is incompatible with MutableGraph, guaranteeing
// only weighted edges can be present in the graph.
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
	EachDataEdge(f func(edge DataEdge))
}

// DataWeightedGraph is the mutable version of a labeled graph. Its
// AddEdges() method is incompatible with MutableGraph, guaranteeing
// only weighted edges can be present in the graph.
type MutableDataGraph interface {
	DataGraph
	VertexSetMutator
	DataEdgeSetMutator
}

/* Graph creation */

type GraphFactory func() interface{}

var Graphs = make(map[string]GraphFactory, 0)

// Creates a new graph instance.
//
// You will need to type assert the returned graph to the interface appropriate
// for your use case: Graph, DirectedGraph, MutableGraph, WeightedGraph, etc.
func New(name string) (graph interface{}, err error) {
	if _, exists := Graphs[name]; !exists {
		return nil, fmt.Errorf("No graph is registered with the name %q", name)
	}

	return Graphs[name](), nil
}

func RegisterGraph(name string, factory GraphFactory) error {
	if _, exists := Graphs[name]; exists {
		return fmt.Errorf("A graph is already registered with the name %q", name)
	}

	g := factory()

	if _, ok := g.(Graph); ok {
		return nil
	} else if _, ok := g.(WeightedGraph); ok {
		return nil
	}

	return fmt.Errorf("Value returned from factory does not implement a known Graph interface")
}

func init() {
	RegisterGraph("basic.directed", func() interface{} {
		return NewDirected()
	})
	RegisterGraph("basic.undirected", func() interface{} {
		return NewUndirected()
	})
	RegisterGraph("weighted.directed", func() interface{} {
		return NewWeightedDirected()
	})
	RegisterGraph("weighted.undirected", func() interface{} {
		return NewWeightedUndirected()
	})
}
