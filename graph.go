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

/* Graph structures */

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
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(edge Edge))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
	HasEdge(e Edge) bool
	Order() int
	Size() int
	InDegree(vertex Vertex) (int, bool)
	OutDegree(vertex Vertex) (int, bool)
}

// DirectedGraph describes a Graph all of whose edges are directed.
//
// Implementing DirectedGraph is the only unambiguous signal gogl provides
// that a graph's edges are directed.
type DirectedGraph interface {
	Graph
	Transpose() DirectedGraph
}

// MutableGraph describes a graph with basic edges (no weighting, labeling, etc.)
// that can be modified freely by adding or removing vertices or edges.
type MutableGraph interface {
	Graph
	EnsureVertex(vertices ...Vertex)
	RemoveVertex(vertices ...Vertex)
	AddEdges(edges ...Edge)
	RemoveEdges(edges ...Edge)
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
	EnsureVertex(vertices ...Vertex)
	RemoveVertex(vertices ...Vertex)
	AddEdges(edges ...WeightedEdge)
	RemoveEdges(edges ...WeightedEdge)
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
	EnsureVertex(vertices ...Vertex)
	RemoveVertex(vertices ...Vertex)
	AddEdges(edges ...LabeledEdge)
	RemoveEdges(edges ...LabeledEdge)
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
