package gogl

import (
	"fmt"
	"reflect"
)

/* Vertex structures */

// As a rule, gogl tries to place as low a requirement on its vertices as
// possible. This is because, from a purely graph theoretic perspective,
// vertices are inert. Boring, even. Graphs are more about the topology, the
// characteristics of the edges connecting the points than the points
// themselves. Your use case cares about the content of your vertices, but gogl
// does not.  Consequently, anything can act as a vertex.
type Vertex interface{}

/* Edge structures */

// A graph's behaviors is primarily a product of the constraints and
// capabilities it places on its edges. These constraints and capabilities
// determine whether certain types of operations are possible on the graph, as
// well as the efficiencies for various operations.

// gogl aims to provide a diverse range of graph implementations that can meet
// the varying constraints and implementation needs, but still achieve optimal
// performance given those constraints.

type Edge interface {
	Source() Vertex
	Target() Vertex
	Both() (Vertex, Vertex)
}

type WeightedEdge interface {
	Edge
	Weight() int
}

// BaseEdge is a struct used internally to represent edges and meet the Edge
// interface requirements. It uses the standard notation, (u,v), for vertex
// pairs in an edge.
type BaseEdge struct {
	U Vertex
	V Vertex
}

func (e BaseEdge) Source() Vertex {
	return e.U
}

func (e BaseEdge) Target() Vertex {
	return e.V
}

func (e BaseEdge) Both() (Vertex, Vertex) {
	return e.U, e.V
}

type BaseWeightedEdge struct {
	BaseEdge
	W int
}

func (e BaseWeightedEdge) Weight() int {
	return e.W
}

/* Graph structures */

type Graph interface {
	EachVertex(f func(vertex Vertex))
	EachEdge(f func(edge Edge))
	EachAdjacent(vertex Vertex, f func(adjacent Vertex))
	HasVertex(vertex Vertex) bool
	Order() int
	Size() int
	InDegree(vertex Vertex) (int, bool)
	OutDegree(vertex Vertex) (int, bool)
}

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

type DirectedGraph interface {
	Graph
	Transpose() DirectedGraph
	IsAcyclic() bool
	GetCycles() [][]Vertex
}

type WeightedGraph interface {
	Graph
	EachWeightedEdge(f func(edge WeightedEdge))
}

type MutableWeightedGraph interface {
	WeightedGraph
	EnsureVertex(vertices ...Vertex)
	RemoveVertex(vertices ...Vertex)
	AddEdges(edges ...WeightedEdge)
	RemoveEdges(edges ...WeightedEdge)
}

/* Initialization for immutable graphs */

type IGInitializer interface {
	EnsureVertex(vertex ...Vertex)
	AddEdges(edges ...Edge) bool
	GetGraph() Graph
}

type ImmutableGraph interface {
	Graph
	EnsureVertex(vertex ...Vertex)
	AddEdges(edges ...Edge)
}

var immutableGraphs map[string]ImmutableGraph

func CreateImmutableGraph(name string) (*ImmutableGraphInitializer, error) {
	template := immutableGraphs[name]
	if template == nil {
		return nil, fmt.Errorf("gogl: Unregistered graph type %s", name)
	}

	// Use reflection to make a copy of the graph template
	v := reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()
	graph, ok := v.(ImmutableGraph)
	if !ok {
		panic(fmt.Sprintf("gogl: Unable to copy graph template: %s (%v)", name, reflect.ValueOf(v).Kind().String()))
	}

	initializer := &ImmutableGraphInitializer{
		Graph: graph,
	}

	return initializer, nil
}

// A ImmutableGraphInitializer provides write-only methods to populate an
// immutable graph.
type ImmutableGraphInitializer struct {
	Graph ImmutableGraph
}

func (gi *ImmutableGraphInitializer) EnsureVertex(vertices ...Vertex) {
	gi.Graph.EnsureVertex(vertices...)
}

func (gi *ImmutableGraphInitializer) AddEdges(edges ...Edge) {
	gi.Graph.AddEdges(edges...)
}

func (gi *ImmutableGraphInitializer) GetGraph() Graph {
	defer func() { gi.Graph = nil }()
	return gi.Graph
}

func init() {
	immutableGraphs = map[string]ImmutableGraph{
		"ual": NewUndirected(),
	}
}
