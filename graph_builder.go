package gogl

// Defines a builder for use in creating graph objects.
import (
	"github.com/lann/builder"
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

type builderImmutableDirected builder.Builder

// Builder/Immutable/Basic/Directed
var BIBD = builder.Register(builderImmutableDirected{}, immutableDirected{}).(builderImmutableDirected)

func (b builderImmutableDirected) From(g Graph) builderImmutableDirected {
	return builder.Set(b, "from", g).(builderImmutableDirected)
}

func (b builderImmutableDirected) Create() *immutableDirected {
	gv := builder.GetStruct(b).(immutableDirected)
	g := &gv
	g.list = make(map[Vertex]map[Vertex]struct{})

	if from, exists := builder.Get(b, "from"); exists {
		from := from.(Graph)
		createDeferredEdgeLambda(from, g)()

		if g.Order() != from.Order() {
			from.EachVertex(func(vertex Vertex) (terminate bool) {
				g.ensureVertex(vertex)
				return
			})
		}
	}

	return g
}
