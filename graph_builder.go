package gogl

// Defines a builder for use in creating graph objects.
import (
	"github.com/lann/builder"
	"sync"
)

type GraphBuilder interface {
	From(g Graph) GraphBuilder
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

func (b builderImmutableDirected) From(g Graph) builderImmutableDirected {
	return builder.Set(b, "from", g).(builderImmutableDirected)
}

func (b builderImmutableDirected) Create() *immutableDirected {
	gv := builder.GetStruct(b).(immutableDirected)
	g := &gv

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(Graph)
		functorToAdjacencyList(from, g)
	}

	return g
}

// Builder/Mutable/Basic/Directed
var BMBD = builder.Register(builderMutableDirected{}, mutableDirected{al_basic_mut{al_basic{list: make(map[Vertex]map[Vertex]struct{})}, sync.RWMutex{}}}).(builderMutableDirected)

type builderMutableDirected builder.Builder

func (b builderMutableDirected) From(g Graph) builderMutableDirected {
	return builder.Set(b, "from", g).(builderMutableDirected)
}

func (b builderMutableDirected) Create() *mutableDirected {
	gv := builder.GetStruct(b).(mutableDirected)
	g := &gv

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(Graph)
		functorToAdjacencyList(from, g)
	}

	return g
}
