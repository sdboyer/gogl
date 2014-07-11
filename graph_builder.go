package gogl

// Defines a builder for use in creating graph objects.
import (
	"sync"

	"github.com/lann/builder"
)

type GraphBuilder interface {
	Using(g GraphSource) GraphBuilder
	Graph() Graph
}

//type graphStructSpawner struct {
//source       Graph
//directed     bool
//edgeType     int
//mutable      bool
//multiplicity int
//}

// Builder/Immutable/Basic/Directed
var BIBD = builder.Register(builderImmutableDirected{}, immutableDirected{al_basic_immut{al_basic{list: make(map[Vertex]map[Vertex]struct{})}}}).(builderImmutableDirected)

type builderImmutableDirected builder.Builder

func (b builderImmutableDirected) From(g GraphSource) builderImmutableDirected {
	return builder.Set(b, "from", g).(builderImmutableDirected)
}

func (b builderImmutableDirected) Create() *immutableDirected {
	gv := builder.GetStruct(b).(immutableDirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]struct{})

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderImmutableDirected) Graph() Graph {
	return b.Create()
}

func (b builderImmutableDirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Basic/Directed
var BMBD = builder.Register(builderMutableDirected{}, mutableDirected{al_basic_mut{al_basic{list: make(map[Vertex]map[Vertex]struct{})}, sync.RWMutex{}}}).(builderMutableDirected)

type builderMutableDirected builder.Builder

func (b builderMutableDirected) From(g GraphSource) builderMutableDirected {
	return builder.Set(b, "from", g).(builderMutableDirected)
}

func (b builderMutableDirected) Create() *mutableDirected {
	gv := builder.GetStruct(b).(mutableDirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]struct{})

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableDirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableDirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Basic/Undirected
var BMBU = builder.Register(builderMutableUndirected{}, mutableUndirected{al_basic_mut{al_basic{list: make(map[Vertex]map[Vertex]struct{})}, sync.RWMutex{}}}).(builderMutableUndirected)

type builderMutableUndirected builder.Builder

func (b builderMutableUndirected) From(g GraphSource) builderMutableUndirected {
	return builder.Set(b, "from", g).(builderMutableUndirected)
}

func (b builderMutableUndirected) Create() *mutableUndirected {
	gv := builder.GetStruct(b).(mutableUndirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]struct{})

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableUndirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableUndirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Weighted/Directed
var BMWD = builder.Register(builderMutableWeightedDirected{}, weightedDirected{baseWeighted{list: make(map[Vertex]map[Vertex]float64), size: 0, mu: sync.RWMutex{}}}).(builderMutableWeightedDirected)

type builderMutableWeightedDirected builder.Builder

func (b builderMutableWeightedDirected) From(g GraphSource) builderMutableWeightedDirected {
	return builder.Set(b, "from", g).(builderMutableWeightedDirected)
}

func (b builderMutableWeightedDirected) Create() *weightedDirected {
	gv := builder.GetStruct(b).(weightedDirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]float64)

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableWeightedDirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableWeightedDirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Weighted/Undirected
var BMWU = builder.Register(builderMutableWeightedUndirected{}, weightedUndirected{baseWeighted{list: make(map[Vertex]map[Vertex]float64), size: 0, mu: sync.RWMutex{}}}).(builderMutableWeightedUndirected)

type builderMutableWeightedUndirected builder.Builder

func (b builderMutableWeightedUndirected) From(g GraphSource) builderMutableWeightedUndirected {
	return builder.Set(b, "from", g).(builderMutableWeightedUndirected)
}

func (b builderMutableWeightedUndirected) Create() *weightedUndirected {
	gv := builder.GetStruct(b).(weightedUndirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]float64)

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableWeightedUndirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableWeightedUndirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Labeled/Directed
var BMLD = builder.Register(builderMutableLabeledDirected{}, labeledDirected{baseLabeled{list: make(map[Vertex]map[Vertex]string), size: 0, mu: sync.RWMutex{}}}).(builderMutableLabeledDirected)

type builderMutableLabeledDirected builder.Builder

func (b builderMutableLabeledDirected) From(g GraphSource) builderMutableLabeledDirected {
	return builder.Set(b, "from", g).(builderMutableLabeledDirected)
}

func (b builderMutableLabeledDirected) Create() *labeledDirected {
	gv := builder.GetStruct(b).(labeledDirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]string)

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableLabeledDirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableLabeledDirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Labeled/Undirected
var BMLU = builder.Register(builderMutableLabeledUndirected{}, labeledUndirected{baseLabeled{list: make(map[Vertex]map[Vertex]string), size: 0, mu: sync.RWMutex{}}}).(builderMutableLabeledUndirected)

type builderMutableLabeledUndirected builder.Builder

func (b builderMutableLabeledUndirected) From(g GraphSource) builderMutableLabeledUndirected {
	return builder.Set(b, "from", g).(builderMutableLabeledUndirected)
}

func (b builderMutableLabeledUndirected) Create() *labeledUndirected {
	gv := builder.GetStruct(b).(labeledUndirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]string)

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableLabeledUndirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableLabeledUndirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Data/Directed
var BMDD = builder.Register(builderMutableDataDirected{}, dataDirected{baseData{list: make(map[Vertex]map[Vertex]interface{}), size: 0, mu: sync.RWMutex{}}}).(builderMutableDataDirected)

type builderMutableDataDirected builder.Builder

func (b builderMutableDataDirected) From(g GraphSource) builderMutableDataDirected {
	return builder.Set(b, "from", g).(builderMutableDataDirected)
}

func (b builderMutableDataDirected) Create() *dataDirected {
	gv := builder.GetStruct(b).(dataDirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]interface{})

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableDataDirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableDataDirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}

// Builder/Mutable/Data/Undirected
var BMDU = builder.Register(builderMutableDataUndirected{}, dataUndirected{baseData{list: make(map[Vertex]map[Vertex]interface{}), size: 0, mu: sync.RWMutex{}}}).(builderMutableDataUndirected)

type builderMutableDataUndirected builder.Builder

func (b builderMutableDataUndirected) From(g GraphSource) builderMutableDataUndirected {
	return builder.Set(b, "from", g).(builderMutableDataUndirected)
}

func (b builderMutableDataUndirected) Create() *dataUndirected {
	gv := builder.GetStruct(b).(dataUndirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]interface{})

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(GraphSource)
		functorToAdjacencyList(from, g)
	}

	return g
}

func (b builderMutableDataUndirected) Graph() Graph {
	return b.Create()
}

func (b builderMutableDataUndirected) Using(g GraphSource) GraphBuilder {
	return b.From(g)
}
