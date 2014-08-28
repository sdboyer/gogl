package spec

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	. "github.com/sdboyer/gocheck"
	. "github.com/sdboyer/gogl"
	"gopkg.in/fatih/set.v0"
)

/////////////////////////////////////////////////////////////////////
//
// GRAPH FIXTURES
//
/////////////////////////////////////////////////////////////////////

type loopArc struct {
	v Vertex
}

func (e loopArc) Both() (Vertex, Vertex) {
	return e.v, e.v
}

func (e loopArc) Source() Vertex {
	return e.v
}

func (e loopArc) Target() Vertex {
	return e.v
}

type loopArcList []Arc

func (el loopArcList) Vertices(fn VertexStep) {
	set := set.NewNonTS()

	for _, e := range el {
		set.Add(e.Both())
	}

	for _, v := range set.List() {
		if fn(v) {
			return
		}
	}
}

func (el loopArcList) EachArc(fn ArcStep) {
	for _, e := range el {
		if _, ok := e.(loopArc); !ok {
			if fn(e) {
				return
			}
		}
	}
}

func (el loopArcList) Edges(fn EdgeStep) {
	for _, e := range el {
		if _, ok := e.(loopArc); !ok {
			if fn(e) {
				return
			}
		}
	}
}

var GraphFixtures = map[string]GraphSource{
	// TODO improve naming basis/patterns for these
	"arctest": ArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
		NewArc("foo", "qux"),
		NewArc("qux", "bar"),
	},
	"pair": ArcList{
		NewArc(1, 2),
	},
	"2e3v": ArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
	},
	"3e4v": ArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
		NewArc("foo", "qux"),
	},
	"3e5v1i": loopArcList{
		NewArc("foo", "bar"),
		NewArc("bar", "baz"),
		NewArc("foo", "qux"),
		loopArc{"isolate"},
	},
	"w-2e3v": WeightedArcList{
		NewWeightedArc(1, 2, 5.23),
		NewWeightedArc(2, 3, 5.821),
	},
	"l-2e3v": LabeledArcList{
		NewLabeledArc(1, 2, "foo"),
		NewLabeledArc(2, 3, "bar"),
	},
	"d-2e3v": DataArcList{
		NewDataArc(1, 2, "foo"),
		NewDataArc(2, 3, struct{ a int }{a: 2}),
	},
}

/////////////////////////////////////////////////////////////////////
//
// HELPERS
//
/////////////////////////////////////////////////////////////////////

// Hook gocheck into the go test runner
func TestHookup(t *testing.T) { TestingT(t) }

// Returns an arc with the directionality swapped.
func Swap(a Arc) Arc {
	return NewArc(a.Target(), a.Source())
}

func gdebug(g Graph, args ...interface{}) {
	fmt.Println("DEBUG: graph type", reflect.New(reflect.Indirect(reflect.ValueOf(g)).Type()))
	pretty.Print(args...)
}

/////////////////////////////////////////////////////////////////////
//
// SUITE SETUP
//
/////////////////////////////////////////////////////////////////////

type graphFactory func(GraphSpec) Graph

func SetUpTestsFromSpec(gp GraphProperties, fn graphFactory) bool {
	var directed bool

	g := fn(GraphSpec{Props: gp})

	fact := func(gs GraphSource) Graph {
		return fn(GraphSpec{Props: gp, Source: gs})
	}

	if _, ok := g.(Digraph); ok {
		directed = true
		Suite(&DigraphSuite{Factory: fact})
	}

	// Set up the basic Graph suite unconditionally
	Suite(&GraphSuite{fact, directed})

	if _, ok := g.(SimpleGraph); ok {
		Suite(&SimpleGraphSuite{fact, directed})
	}

	if _, ok := g.(VertexSetMutator); ok {
		Suite(&VertexSetMutatorSuite{fact})
	}

	if _, ok := g.(EdgeSetMutator); ok {
		Suite(&EdgeSetMutatorSuite{fact})
	}

	if _, ok := g.(ArcSetMutator); ok {
		Suite(&ArcSetMutatorSuite{fact})
	}

	if _, ok := g.(WeightedGraph); ok {
		wfact := func(gs GraphSource) WeightedGraph {
			return fact(gs).(WeightedGraph)
		}

		Suite(&WeightedGraphSuite{wfact})

		if _, ok := g.(WeightedDigraph); ok {
			Suite(&WeightedDigraphSuite{wfact})
		}
		if _, ok := g.(WeightedEdgeSetMutator); ok {
			Suite(&WeightedEdgeSetMutatorSuite{wfact})
		}
		if _, ok := g.(WeightedArcSetMutator); ok {
			Suite(&WeightedArcSetMutatorSuite{wfact})
		}
	}

	if _, ok := g.(LabeledGraph); ok {
		wfact := func(gs GraphSource) LabeledGraph {
			return fact(gs).(LabeledGraph)
		}

		Suite(&LabeledGraphSuite{wfact})

		if _, ok := g.(LabeledDigraph); ok {
			Suite(&LabeledDigraphSuite{wfact})
		}
		if _, ok := g.(LabeledEdgeSetMutator); ok {
			Suite(&LabeledEdgeSetMutatorSuite{wfact})
		}
		if _, ok := g.(LabeledArcSetMutator); ok {
			Suite(&LabeledArcSetMutatorSuite{wfact})
		}
	}

	if _, ok := g.(DataGraph); ok {
		wfact := func(gs GraphSource) DataGraph {
			return fact(gs).(DataGraph)
		}

		Suite(&DataGraphSuite{wfact})

		if _, ok := g.(DataDigraph); ok {
			Suite(&DataDigraphSuite{wfact})
		}
		if _, ok := g.(DataEdgeSetMutator); ok {
			Suite(&DataEdgeSetMutatorSuite{wfact})
		}
		if _, ok := g.(DataArcSetMutator); ok {
			Suite(&DataArcSetMutatorSuite{wfact})
		}
	}

	return true
}
