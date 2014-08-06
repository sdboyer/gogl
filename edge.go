package gogl

/* Edge interfaces */

// A graph's behaviors are primarily a product of the constraints and
// capabilities it places on its edges. These constraints and capabilities
// determine whether certain types of operations are possible on the graph, as
// well as the efficiencies for various operations.

// gogl aims to provide a range of graph implementations that can meet
// the varying constraints and implementation needs, but still achieve optimal
// performance given those constraints.

// Edge describes an undirected connection between two vertices.
type Edge interface {
	Both() (u Vertex, v Vertex) // No order consistency is implied.
}

// Arc describes a directed connection between two vertices.
type Arc interface {
	Both() (u Vertex, v Vertex) // u is tail/source, v is head/target.
	Source() Vertex
	Target() Vertex
}

// WeightedEdge describes an Edge that carries a numerical weight.
type WeightedEdge interface {
	Edge
	Weight() float64
}

// WeightedArc describes an Arc that carries a numerical weight.
type WeightedArc interface {
	Arc
	Weight() float64
}

// LabeledEdge describes an Edge that also has a string label.
type LabeledEdge interface {
	Edge
	Label() string
}

// LabeledArc describes an Arc that also has a string label.
type LabeledArc interface {
	Arc
	Label() string
}

// DataEdge describes an Edge that also holds arbitrary data.
type DataEdge interface {
	Edge
	Data() interface{}
}

// DataArc describes an Arc that also holds arbitrary data.
type DataArc interface {
	Arc
	Data() interface{}
}

/* Base implementations of Edge interfaces */

// BaseEdge is a struct used to represent edges and meet the Edge interface
// requirements. It uses the standard graph notation, (U,V), for its
// contained vertex pair.
type baseEdge struct {
	u Vertex
	v Vertex
}

func (e baseEdge) Both() (Vertex, Vertex) {
	return e.u, e.v
}

// Create a new basic edge.
func NewEdge(u, v Vertex) Edge {
	return baseEdge{u: u, v: v}
}

type baseArc struct {
	baseEdge
}

func (e baseArc) Source() Vertex {
	return e.u
}

func (e baseArc) Target() Vertex {
	return e.v
}

// Create a new basic arc.
func NewArc(u, v Vertex) Arc {
	return baseArc{baseEdge{u: u, v: v}}
}

// BaseWeightedEdge extends BaseEdge with weight data.
type baseWeightedEdge struct {
	baseEdge
	w float64
}

func (e baseWeightedEdge) Weight() float64 {
	return e.w
}

// Create a new weighted edge.
func NewWeightedEdge(u, v Vertex, weight float64) WeightedEdge {
	return baseWeightedEdge{baseEdge{u: u, v: v}, weight}
}

type baseWeightedArc struct {
	baseArc
	w float64
}

func (e baseWeightedArc) Weight() float64 {
	return e.w
}

// Create a new weighted arc.
func NewWeightedArc(u, v Vertex, weight float64) WeightedArc {
	return baseWeightedArc{baseArc{baseEdge{u: u, v: v}}, weight}
}

// BaseLabeledEdge extends BaseEdge with label data.
type baseLabeledEdge struct {
	baseEdge
	l string
}

func (e baseLabeledEdge) Label() string {
	return e.l
}

// Create a new labeled edge.
func NewLabeledEdge(u, v Vertex, label string) LabeledEdge {
	return baseLabeledEdge{baseEdge{u: u, v: v}, label}
}

// BaseLabeledArc extends BaseArc with label data.
type baseLabeledArc struct {
	baseArc
	l string
}

func (e baseLabeledArc) Label() string {
	return e.l
}

// Create a new labeled arc.
func NewLabeledArc(u, v Vertex, label string) LabeledArc {
	return baseLabeledArc{baseArc{baseEdge{u: u, v: v}}, label}
}

// BaseDataEdge extends BaseEdge with arbitrary data.
type baseDataEdge struct {
	baseEdge
	d interface{}
}

func (e baseDataEdge) Data() interface{} {
	return e.d
}

// Create a new "data" edge - an edge with arbitrary embedded data.
func NewDataEdge(u, v Vertex, data interface{}) DataEdge {
	return baseDataEdge{baseEdge{u: u, v: v}, data}
}

// BaseDataArc extends BaseArc with arbitrary data.
type baseDataArc struct {
	baseArc
	d interface{}
}

func (e baseDataArc) Data() interface{} {
	return e.d
}

// Create a new "data" edge - an edge with arbitrary embedded data.
func NewDataArc(u, v Vertex, data interface{}) DataArc {
	return baseDataArc{baseArc{baseEdge{u: u, v: v}}, data}
}
