package test_bundle

import (
	"fmt"
	"github.com/sdboyer/gogl"
	//"github.com/sdboyer/gogl/adjacency_list"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var _ = fmt.Println

var edgeSet = []gogl.Edge{
	&gogl.BaseEdge{"foo", "bar"},
	&gogl.BaseEdge{"bar", "baz"},
}

/*
func EnsureBasicGraphBehaviors(g gogl.Graph, t *testing.T) {
	fml("Type:", reflect.TypeOf(g))
}

func DoItWithFCF(f func(...gogl.Edge) gogl.MutableGraph, t *testing.T) {
	g := f()
	fml("FCF Type:", reflect.TypeOf(g))
	fml("FCF Value:", reflect.ValueOf(g))

	g.EnsureVertex("foo")
	g.EnsureVertex("bar")
	ff := func(v gogl.Vertex) {
		fml(v.(string))
	}
	g.EachVertex(ff)
	fml(g)
	//fml(g2)

	rg := reflect.New(reflect.TypeOf(g))
	fml("FCF2 Type:", reflect.TypeOf(rg))
	fml(rg)
}
*/

func GraphTestVertexMembership(f func(...gogl.Edge) gogl.MutableGraph, t *testing.T) {
	g := f()

	Convey("Test adding, removal, and membership of string literal vertex.", t, func() {
		So(g.HasVertex("foo"), ShouldEqual, false)
		g.EnsureVertex("foo")
		So(g.HasVertex("foo"), ShouldEqual, true)
		g.RemoveVertex("foo")
		So(g.HasVertex("foo"), ShouldEqual, false)
	})

	Convey("Test adding, removal, and membership of int literal vertex.", t, func() {
		So(g.HasVertex(1), ShouldEqual, false)
		g.EnsureVertex(1)
		So(g.HasVertex(1), ShouldEqual, true)
		g.RemoveVertex(1)
		So(g.HasVertex(1), ShouldEqual, false)
	})

	Convey("Test adding, removal, and membership of composite literal vertex.", t, func() {
		So(g.HasVertex(edgeSet[0]), ShouldEqual, false)
		g.EnsureVertex(edgeSet[0])
		So(g.HasVertex(edgeSet[0]), ShouldEqual, true)

		Convey("No membership match on new struct with same values or new pointer", func() {
			So(g.HasVertex(gogl.BaseEdge{"foo", "bar"}), ShouldEqual, false)
			So(g.HasVertex(&gogl.BaseEdge{"foo", "bar"}), ShouldEqual, false)
		})

		g.RemoveVertex(edgeSet[0])
		So(g.HasVertex(edgeSet[0]), ShouldEqual, false)
	})

}
