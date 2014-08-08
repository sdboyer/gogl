package dfs

import (
	"fmt"
	"testing"

	. "github.com/sdboyer/gocheck"
	"github.com/sdboyer/gogl"
	"github.com/sdboyer/gogl/graph/al"
)

// Hook gocheck into the go test runner
func Test(t *testing.T) { TestingT(t) }

var dfEdgeSet = gogl.EdgeList{
	gogl.NewEdge("foo", "bar"),
	gogl.NewEdge("bar", "baz"),
	gogl.NewEdge("baz", "qux"),
}

var dfArcSet = gogl.ArcList{
	gogl.NewArc("foo", "bar"),
	gogl.NewArc("bar", "baz"),
	gogl.NewArc("baz", "qux"),
}

type DepthFirstSearchSuite struct{}

var _ = Suite(&DepthFirstSearchSuite{})

// Basic test of outermost search functionality.
func (s *DepthFirstSearchSuite) TestSearch(c *C) {
	// extra edge demonstrates non-productive search paths are not included
	extraSet := append(dfArcSet, gogl.NewArc("bar", "quark"))
	// directed
	g := gogl.Spec().Directed().Using(extraSet).
		Create(al.G).(gogl.Digraph)

	path, err := Search(g, "qux", "bar")
	c.Assert(path, DeepEquals, []gogl.Vertex{"qux", "baz", "bar"})
	c.Assert(err, IsNil)

	// undirected
	//ug := gogl.BuildGraph().Using(extraSet).Create(al.AdjacencyList)

	// TODO accidentally passed wrong graph - fix!
	path, err = Search(g, "qux", "bar")
	c.Assert(path, DeepEquals, []gogl.Vertex{"qux", "baz", "bar"})
	c.Assert(err, IsNil)
}

func (s *DepthFirstSearchSuite) TestSearchVertexVerification(c *C) {
	g := gogl.Spec().Mutable().Directed().
		Create(al.G).(gogl.MutableDigraph)
	g.EnsureVertex("foo")

	_, err := Search(g.(gogl.Digraph), "foo", "bar")
	c.Assert(err, ErrorMatches, "Start vertex.*")
	_, err = Search(g.(gogl.Digraph), "bar", "foo")
	c.Assert(err, ErrorMatches, "Target vertex.*")
}

func (s *DepthFirstSearchSuite) TestFindSources(c *C) {
	g := gogl.Spec().Directed().
		Mutable().Using(dfArcSet).
		Create(al.G).(gogl.Digraph)

	sources, err := FindSources(g)
	c.Assert(fmt.Sprint(sources), Equals, fmt.Sprint([]gogl.Vertex{"foo"}))
	c.Assert(err, IsNil)

	// Ensure it finds multiple, as well
	g.(gogl.MutableDigraph).AddArcs(gogl.NewArc("quark", "baz"))
	sources, err = FindSources(g)

	possibles := [][]gogl.Vertex{
		[]gogl.Vertex{"foo", "quark"},
		[]gogl.Vertex{"quark", "foo"},
	}
	c.Assert(possibles, Contains, sources)
	c.Assert(err, IsNil)
}

func (s *DepthFirstSearchSuite) TestToposort(c *C) {
	g := gogl.Spec().Directed().
		Mutable().Using(dfArcSet).
		Create(al.G).(gogl.Digraph)

	tsl, err := Toposort(g, "foo")
	c.Assert(err, IsNil)
	c.Assert(tsl, DeepEquals, []gogl.Vertex{"qux", "baz", "bar", "foo"})

	// add a cycle, ensure error comes back
	g.(gogl.MutableDigraph).AddArcs(gogl.NewArc("bar", "foo"))
	tsl, err = Toposort(g, "foo")
	c.Assert(err, ErrorMatches, "Cycle detected in graph")

	// undirected
	ug := gogl.Spec().Using(dfEdgeSet).Create(al.G)

	_, err = Toposort(ug)
	c.Assert(err, ErrorMatches, ".*do not have sources.*")

	tsl, err = Toposort(ug, "foo")
	// no such thing as a 'cycle' (of that kind) in undirected graphs
	c.Assert(err, IsNil)
	c.Assert(tsl, DeepEquals, []gogl.Vertex{"qux", "baz", "bar", "foo"})
}

// This is a bit wackyhacky, but works well enough
var _ = Suite(&TestVisitor{})

type TestVisitor struct {
	c           *C
	vertices    []string
	colors      map[string]int
	found_edges []gogl.Arc
}

func (v *TestVisitor) OnBackEdge(vertex gogl.Vertex) {
	vtx := vertex.(string)
	v.c.Assert(v.colors[vtx], Equals, grey)
}

func (v *TestVisitor) OnStartVertex(vertex gogl.Vertex) {
	vtx := vertex.(string)
	v.c.Assert(v.colors[vtx], Equals, white)
	v.colors[vtx] = grey
}

func (v *TestVisitor) OnExamineEdge(edge gogl.Edge) {
	v.c.Assert(v.found_edges, Not(Contains), edge)
	v.found_edges = append(v.found_edges, edge.(gogl.Arc))
}

func (v *TestVisitor) OnFinishVertex(vertex gogl.Vertex) {
	vtx := vertex.(string)
	v.c.Assert(v.colors[vtx], Equals, grey)
	v.colors[vtx] = black
}

func (v *TestVisitor) TestTraverse(c *C) {
	v.c = c

	el := gogl.ArcList{
		gogl.NewArc("foo", "bar"),
		gogl.NewArc("bar", "baz"),
		gogl.NewArc("bar", "foo"),
		gogl.NewArc("bar", "quark"),
		gogl.NewArc("baz", "qux"),
	}
	g := gogl.Spec().Directed().
		Mutable().Using(el).
		Create(al.G).(gogl.Digraph)

	v.vertices = []string{"foo", "bar", "baz", "qux", "quark"}

	v.colors = make(map[string]int)
	for _, vtx := range v.vertices {
		v.colors[vtx] = white
	}

	v.found_edges = make([]gogl.Arc, 0)

	Traverse(g, v, "foo")

	for vertex, color := range v.colors {
		c.Log("Checking that vertex '", vertex, "' has been finished")
		c.Assert(color, Equals, black)
	}

	for _, e := range el {
		c.Assert(v.found_edges, Contains, e)
	}
	c.Assert(len(v.found_edges), Equals, len(el))
}

type LinkedListSuite struct{}

var _ = Suite(&LinkedListSuite{})

func (s *LinkedListSuite) TestStack(c *C) {
	stack := vstack{}

	c.Assert(stack.length(), Equals, 0)

	stack.push("foo")
	c.Assert(stack.length(), Equals, 1)

	stack.push("bar")
	c.Assert(stack.length(), Equals, 2)
	c.Assert(stack.pop(), Equals, "bar")
	c.Assert(stack.pop(), Equals, "foo")
	c.Assert(stack.pop(), IsNil)
	c.Assert(stack.length(), Equals, 0)
}

func (s *LinkedListSuite) TestQueue(c *C) {
	queue := vqueue{}

	c.Assert(queue.length(), Equals, 0)

	queue.push("foo")
	c.Assert(queue.length(), Equals, 1)

	queue.push("bar")
	c.Assert(queue.length(), Equals, 2)
	c.Assert(queue.pop(), Equals, "foo")
	c.Assert(queue.pop(), Equals, "bar")
	c.Assert(queue.pop(), IsNil)
	c.Assert(queue.length(), Equals, 0)
}
