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

type containsChecker struct {
	*CheckerInfo
}

var contains Checker = &containsChecker{
	&CheckerInfo{Name: "Contains", Params: []string{"haystack", "needle"}},
}

func (cc *containsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	needle := reflect.ValueOf(params[1])
	haystack := reflect.ValueOf(params[0])

	if reflect.SliceOf(needle.Type()) != haystack.Type() {
		return false, "Haystack must be a slice with the same element type as needle."
	}

	length := haystack.Len()
	for i := 0; i < length; i++ {
		if reflect.DeepEqual(haystack.Index(i).Interface(), needle.Interface()) {
			return true, ""
		}
	}

	return false, ""
}

type DepthFirstSearchSuite struct{}

var _ = Suite(&DepthFirstSearchSuite{})

// Basic test of outermost search functionality.
func (s *DepthFirstSearchSuite) TestSearch(c *C) {
	// directed
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)

	path, err := Search(g, "qux", "bar")
	c.Assert(path, DeepEquals, []gogl.Vertex{"qux", "baz", "bar"})
	c.Assert(err, IsNil)

	// undirected
	ug := gogl.NewUndirected()
	ug.AddEdges(dfEdgeSet...)

	path, err = Search(g, "qux", "bar")
	c.Assert(path, DeepEquals, []gogl.Vertex{"qux", "baz", "bar"})
	c.Assert(err, IsNil)
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
	c.Assert(tsl, DeepEquals, []gogl.Vertex{"foo", "bar", "baz", "qux"})

	//possibles := [][]gogl.Vertex{
	//[]gogl.Vertex{},
	//}
	// undirected
	ug := gogl.NewUndirected()
	ug.AddEdges(dfEdgeSet...)

	_, err = Toposort(g)
	c.Assert(err, ErrorMatches, ".*do not have sources.*")

	tsl, err = Toposort(g, "foo")
	c.Assert(err, IsNil)
	c.Assert(tsl, DeepEquals, []gogl.Vertex{"foo", "bar", "baz", "qux"})

}
