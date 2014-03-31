package dfs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sdboyer/gogl"
	. "launchpad.net/gocheck"
)

// Hook gocheck into the go test runner
func Test(t *testing.T) { TestingT(t) }

var dfEdgeSet = []gogl.Edge{
	&gogl.BaseEdge{"foo", "bar"},
	&gogl.BaseEdge{"bar", "baz"},
	&gogl.BaseEdge{"baz", "qux"},
}

func not(checker Checker) Checker {
	return &notChecker{checker}
}

type notChecker struct {
	sub Checker
}

func (checker *notChecker) Info() *CheckerInfo {
	info := *checker.sub.Info()
	info.Name = "Not(" + info.Name + ")"
	return &info
}

func (checker *notChecker) Check(params []interface{}, names []string) (result bool, error string) {
	result, error = checker.sub.Check(params, names)
	result = !result
	if result {
		return result, ""
	} else {
		info := *checker.sub.Info()
		return result, info.Name + " succeeded."
	}
}

type containsChecker struct {
	*CheckerInfo
}

var contains Checker = &containsChecker{
	&CheckerInfo{Name: "Contains", Params: []string{"haystack", "needle"}},
}

func (cc *containsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	needle := reflect.ValueOf(params[1])
	haystack := reflect.ValueOf(params[0])

	switch haystack.Kind() {
	default:
		return false, "Haystack must be a slice or array."
	case reflect.Slice, reflect.Array:
	}

	if haystack.Len() == 0 {
		return false, "Haystack is empty."
	}

	typealign := true
	if reflect.SliceOf(needle.Type()) != haystack.Type() {
		if haystack.Index(0).Kind() == reflect.Interface &&
			haystack.Index(0).Elem().Type() == needle.Type() {
			typealign = true
		} else {
			typealign = false
		}
	}

	if !typealign {
		return false, "Haystack must be a slice with the same element type as needle."
	}

	length := haystack.Len()
	for i := 0; i < length; i++ {
		if reflect.DeepEqual(haystack.Index(i).Interface(), needle.Interface()) {
			return true, ""
		}
	}

	return false, "Item not found in collection."
}

type DepthFirstSearchSuite struct{}

var _ = Suite(&DepthFirstSearchSuite{})

// Basic test of outermost search functionality.
func (s *DepthFirstSearchSuite) TestSearch(c *C) {
	// directed
	g := gogl.NewDirected()

	// must demonstrate that non-productive search paths are not included
	edgeset := []gogl.Edge{
		&gogl.BaseEdge{"foo", "bar"},
		&gogl.BaseEdge{"bar", "baz"},
		&gogl.BaseEdge{"bar", "quark"},
		&gogl.BaseEdge{"baz", "qux"},
	}

	g.AddEdges(edgeset...)

	path, err := Search(g, "qux", "bar")
	c.Assert(path, DeepEquals, []gogl.Vertex{"qux", "baz", "bar"})
	c.Assert(err, IsNil)

	// undirected
	ug := gogl.NewUndirected()
	ug.AddEdges(edgeset...)

	path, err = Search(g, "qux", "bar")
	c.Assert(path, DeepEquals, []gogl.Vertex{"qux", "baz", "bar"})
	c.Assert(err, IsNil)
}

func (s *DepthFirstSearchSuite) TestSearchVertexVerification(c *C) {
	g := gogl.NewDirected()
	g.EnsureVertex("foo")

	_, err := Search(g, "foo", "bar")
	c.Assert(err, ErrorMatches, "Start vertex.*")
	_, err = Search(g, "bar", "foo")
	c.Assert(err, ErrorMatches, "Target vertex.*")
}

func (s *DepthFirstSearchSuite) TestFindSources(c *C) {
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)

	sources, err := FindSources(g)
	c.Assert(fmt.Sprint(sources), Equals, fmt.Sprint([]gogl.Vertex{"foo"}))
	c.Assert(err, IsNil)

	// Ensure it finds multiple, as well
	g.AddEdges(&gogl.BaseEdge{"quark", "baz"})
	sources, err = FindSources(g)

	possibles := [][]gogl.Vertex{
		[]gogl.Vertex{"foo", "quark"},
		[]gogl.Vertex{"quark", "foo"},
	}
	c.Assert(possibles, contains, sources)
	c.Assert(err, IsNil)
}

func (s *DepthFirstSearchSuite) TestToposort(c *C) {
	// directed
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)

	tsl, err := Toposort(g, "foo")
	c.Assert(err, IsNil)
	c.Assert(tsl, DeepEquals, []gogl.Vertex{"qux", "baz", "bar", "foo"})

	// add a cycle, ensure error comes back
	g.AddEdges(gogl.BaseEdge{"bar", "foo"})
	tsl, err = Toposort(g, "foo")
	c.Assert(err, ErrorMatches, "Cycle detected in graph")

	// undirected
	ug := gogl.NewUndirected()
	ug.AddEdges(dfEdgeSet...)

	_, err = Toposort(ug)
	c.Assert(err, ErrorMatches, ".*do not have sources.*")

	tsl, err = Toposort(ug, "foo")
	// no such thing as a 'cycle' (of that kind) in undirected graphs
	c.Assert(err, IsNil)
	c.Assert(tsl, DeepEquals, []gogl.Vertex{"qux", "baz", "bar", "foo"})
}

func (s *DepthFirstSearchSuite) TestTraverse(c *C) {
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)
}

// This is a bit wackyhacky, but works well enough
var _ = Suite(&TestVisitor{})

type TestVisitor struct {
	c           *C
	vertices    []string
	colors      map[string]int
	found_edges []gogl.Edge
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
	v.c.Assert(v.found_edges, not(contains), edge)
	v.found_edges = append(v.found_edges, edge)
}

func (v *TestVisitor) OnFinishVertex(vertex gogl.Vertex) {
	vtx := vertex.(string)
	v.c.Assert(v.colors[vtx], Equals, grey)
	v.colors[vtx] = black
}

func (v *TestVisitor) TestTraverse(c *C) {
	v.c = c
	g := gogl.NewDirected()

	edgeset := []gogl.Edge{
		gogl.BaseEdge{"foo", "bar"},
		gogl.BaseEdge{"bar", "baz"},
		gogl.BaseEdge{"bar", "foo"},
		gogl.BaseEdge{"bar", "quark"},
		gogl.BaseEdge{"baz", "qux"},
	}
	g.AddEdges(edgeset...)

	v.vertices = []string{"foo", "bar", "baz", "qux", "quark"}

	v.colors = make(map[string]int)
	for _, vtx := range v.vertices {
		v.colors[vtx] = white
	}

	v.found_edges = make([]gogl.Edge, 0)

	Traverse(g, v, "foo")

	for vertex, color := range v.colors {
		c.Log("Checking that vertex '", vertex, "' has been finished")
		c.Assert(color, Equals, black)
	}

	for _, e := range edgeset {
		c.Assert(v.found_edges, contains, e)
	}
	c.Assert(len(v.found_edges), Equals, len(edgeset))
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
