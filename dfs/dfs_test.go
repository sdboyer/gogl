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
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)

	path, err := Search(g, "qux", "bar")
	c.Assert(fmt.Sprint(path), Equals, fmt.Sprint([]gogl.Vertex{"qux", "baz", "bar"}))
	c.Assert(err, IsNil)
}

func sliceEquals(a, b []gogl.Vertex) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestTslGeneration(t *testing.T) {
	g := gogl.NewDirected()
	g.AddEdges(dfEdgeSet...)

	vis := &DFTslVisitor{}
	DepthFirstFromVertices(g, vis, "foo")

	if !sliceEquals(vis.GetTsl(), []gogl.Vertex{"qux", "baz", "bar", "foo"}) {
		t.Error("TSL is not correct.")
	}
}
